package usecase

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go-crud/internal/helper"
	"go-crud/internal/models"
	"go-crud/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"sync"
	"time"
)

type AuthUsecase struct {
	Repository repository.UserRepositoryInterface
	Validate   *validator.Validate
	Viper      *viper.Viper
	Log        *logrus.Logger
}

func NewAuthUsecase(repository repository.UserRepositoryInterface, validator *validator.Validate, viper *viper.Viper, log *logrus.Logger) *AuthUsecase {
	return &AuthUsecase{
		Repository: repository,
		Validate:   validator,
		Viper:      viper,
		Log:        log,
	}
}

func (c *AuthUsecase) ValidateRequest(request *models.SignInRequest) error {
	err := c.Validate.Struct(request)

	if err != nil {
		c.Log.Errorf("%v", err)
		if e := err.(validator.ValidationErrors); e != nil {
			message := helper.GetFirstValidationErrorAndConvert(err)
			return &models.ErrorResponse{
				Code:    400,
				Status:  "Bad Request",
				Message: message,
			}
		}

		return &models.ErrorResponse{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Something Wrong",
		}

	}

	return nil
}

func (c *AuthUsecase) GenerateAccessToken(userID string) (string, error) {
	accessTokenKey := c.Viper.GetString("token.key.access")
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(1 * time.Hour).Unix(),
		"sub": userID,
	})
	token, err := jwtToken.SignedString([]byte(accessTokenKey))
	if err != nil {
		c.Log.Errorf("%v", err)
		return "", &models.ErrorResponse{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Something error",
		}
	}

	return token, nil

}

func (c *AuthUsecase) GenerateRefreshToken(userID string) (string, error) {
	refreshTokenKey := c.Viper.GetString("token.key.refresh")
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(1 * time.Hour).Unix(),
		"sub": userID,
	})
	token, err := jwtToken.SignedString([]byte(refreshTokenKey))
	if err != nil {
		c.Log.Errorf("%v", err)
		return "", &models.ErrorResponse{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Something error",
		}
	}

	return token, nil

}
func (c *AuthUsecase) SignIn(request *models.SignInRequest) (*models.SignInResponse, error) {
	err := c.ValidateRequest(request)
	if err != nil {
		return nil, err
	}
	user, err := c.Repository.FindOneByEmail(request.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.Log.Errorf("%v", err)
		return nil, &models.ErrorResponse{
			Code:    500,
			Message: "Something error",
			Status:  "Internal Server Error",
		}
	}

	if user == nil {
		c.Log.WithFields(logrus.Fields{
			"email": request.Email,
		}).Warn("User not found")
		return nil, &models.ErrorResponse{
			Code:    401,
			Message: "Email or password is incorrect",
			Status:  "Unauthorized",
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		c.Log.WithFields(logrus.Fields{
			"password": request.Password,
		}).Warn("Password not match")
		return nil, &models.ErrorResponse{
			Code:    401,
			Message: "Email or password is incorrect",
			Status:  "Unauthorized",
		}
	}

	var wg sync.WaitGroup
	wg.Add(2)
	var accessToken string
	var refreshToken string
	var errorChannel = make(chan error)
	go func() {
		defer wg.Done()
		accessToken, err = c.GenerateAccessToken(user.Id)
		if err != nil {
			errorChannel <- err
			return
		}

		errorChannel <- nil

	}()

	go func() {
		defer wg.Done()
		refreshToken, err = c.GenerateRefreshToken(user.Id)
		if err != nil {
			errorChannel <- err
			return
		}

		errorChannel <- nil

	}()

	for i := 0; i < 2; i++ {
		if e := <-errorChannel; e != nil {
			c.Log.Errorf("%v", e)
			return nil, &models.ErrorResponse{
				Code:    500,
				Message: "Something Error",
				Status:  "Internal Server Error",
			}
		}
	}
	wg.Wait()

	return &models.SignInResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}
