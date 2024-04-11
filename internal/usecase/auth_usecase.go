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
	"strings"
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

func (c *AuthUsecase) VerifyRefreshToken(refreshToken string) (string, error) {
	if refreshToken == "" {
		return "", &models.ErrorResponse{
			Code:    401,
			Status:  "Unauthorized",
			Message: "You're unauthorized",
		}
	}
	refreshTokenKey := c.Viper.GetString("token.key.refresh")

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(refreshTokenKey), nil
	})

	if err != nil {
		c.Log.WithError(err).Error("Error while parsing refresh token")
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", &models.ErrorResponse{
				Code:    401,
				Status:  "Unauthorized",
				Message: "Token is expired",
			}
		}

		if errors.Is(err, jwt.ErrInvalidKey) {
			return "", &models.ErrorResponse{
				Code:    401,
				Status:  "Unauthorized",
				Message: "Refresh token key is invalid",
			}
		}

		return "", &models.ErrorResponse{
			Code:    401,
			Status:  "Unauthorized",
			Message: "Invalid token",
		}
	}

	var sub string
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		sub = claims["sub"].(string)
	}

	return sub, nil

}

func (c *AuthUsecase) RefreshToken(refreshToken string) (*models.AuthResponse, error) {
	sub, err := c.VerifyRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	accessToken, err := c.GenerateAccessToken(sub)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		AccessToken: accessToken,
	}, nil
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
func (c *AuthUsecase) SignIn(request *models.SignInRequest) (*models.AuthResponse, error) {
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

	return &models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

func (c *AuthUsecase) ParseTokenFromHeader(authorization string) (string, error) {
	if authorization == "" {
		return "", &models.ErrorResponse{
			Code:    401,
			Message: "You're unauthorized",
			Status:  "Unauthorized",
		}
	}
	if !strings.HasPrefix(authorization, "Bearer ") {
		return "", &models.ErrorResponse{
			Code:    401,
			Message: "Invalid authorization",
			Status:  "Unauthorized",
		}
	}

	return strings.TrimPrefix(authorization, "Bearer "), nil

}

func (c *AuthUsecase) VerifyAccessToken(accessToken string) (string, error) {
	if accessToken == "" {
		return "", &models.ErrorResponse{
			Code:    401,
			Status:  "Unauthorized",
			Message: "You're unauthorized",
		}
	}
	accessTokenKey := c.Viper.GetString("token.key.access")
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(accessTokenKey), nil
	})

	if err != nil {
		c.Log.WithError(err).Error("Error parsing token")
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", &models.ErrorResponse{
				Code:    401,
				Status:  "Unauthorized",
				Message: "Token expired",
			}

		}

		if errors.Is(err, jwt.ErrInvalidKey) {
			return "", &models.ErrorResponse{
				Code:    401,
				Status:  "Unauthorized",
				Message: "Invalid key",
			}
		}

		if !token.Valid {
			return "", &models.ErrorResponse{
				Code:    401,
				Status:  "Unauthorized",
				Message: "Invalid token",
			}
		}

		return "", &models.ErrorResponse{
			Code:    401,
			Status:  "Unauthorized",
			Message: err.Error(),
		}
	}

	var sub string
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		sub = claims["sub"].(string)
	}

	return sub, nil

}
