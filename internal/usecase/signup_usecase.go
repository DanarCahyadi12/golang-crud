package usecase

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go-crud/internal/entity"
	"go-crud/internal/helper"
	"go-crud/internal/models"
	"go-crud/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SignUpUsecase struct {
	Repository repository.UserRepositoryInterface
	Validate   *validator.Validate
	Log        *logrus.Logger
}

func NewSignUpUsecase(repository repository.UserRepositoryInterface, validator *validator.Validate, log *logrus.Logger) *SignUpUsecase {
	return &SignUpUsecase{
		Repository: repository,
		Validate:   validator,
		Log:        log,
	}
}

func (u *SignUpUsecase) ValidateRequest(req *models.SignUpRequest) error {
	err := u.Validate.Struct(req)
	if err != nil {
		u.Log.Warnf("Error validating request from client: %v", err)
		message := helper.GetFirstValidationErrorAndConvert(err)
		return &models.ErrorResponse{
			Code:    400,
			Status:  "Bad Request",
			Message: message,
		}
	}
	user, err := u.Repository.FindOneByEmail(req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		u.Log.Warnf("Error while getting user by email: %v", err)
		return &models.ErrorResponse{
			Code:    500,
			Message: "Something error",
			Status:  "Internal Server Error",
		}

	}

	if user != nil {
		return &models.ErrorResponse{
			Code:    400,
			Message: "Email already exists",
			Status:  "Bad Request",
		}
	}

	return nil

}

func (u *SignUpUsecase) CreateUser(req *models.SignUpRequest) (*models.SignUpResponse, error) {
	err := u.ValidateRequest(req)
	if err != nil {
		return nil, err
	}
	hashedPassword, err := u.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	user := &entity.User{
		Id:       uuid.New().String(),
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}

	err = u.Repository.Save(user)
	if err != nil {
		u.Log.Warnf("Error while creating user: %v", err)
		return nil, &models.ErrorResponse{Code: 500, Message: "Something error", Status: "Internal Server Error"}
	}

	return &models.SignUpResponse{
		Id:        user.Id,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil

}

func (u *SignUpUsecase) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", &models.ErrorResponse{
			Code:    500,
			Message: "Something Error",
			Status:  "Internal Server Error",
		}
	}

	return string(hashedPassword), nil
}
