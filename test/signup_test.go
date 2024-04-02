package test

import (
	"github.com/stretchr/testify/require"
	"go-crud/internal/entity"
	"go-crud/internal/models"
	"go-crud/internal/usecase"
	"testing"
)

func TestSignupUsecase(t *testing.T) {
	signupUsecase := usecase.NewSignUpUsecase(userRepositoryMock, validate, log)
	t.Run("Validate request", func(t *testing.T) {
		t.Run("Empty name", func(t *testing.T) {
			req := &models.SignUpRequest{
				Name:     "",
				Email:    "danar@gmail.com",
				Password: "12339",
			}

			err := signupUsecase.ValidateRequest(req)
			require.NotNil(t, err)

		})

		t.Run("Empty email", func(t *testing.T) {
			req := &models.SignUpRequest{
				Name:     "Danar",
				Email:    "",
				Password: "12339",
			}

			err := signupUsecase.ValidateRequest(req)
			require.NotNil(t, err)

		})

		t.Run("Invalid email", func(t *testing.T) {
			req := &models.SignUpRequest{
				Name:     "Danar",
				Email:    "damar@yworngloam",
				Password: "12339",
			}

			err := signupUsecase.ValidateRequest(req)
			require.NotNil(t, err)

		})

		t.Run("Empty password", func(t *testing.T) {
			req := &models.SignUpRequest{
				Name:     "Danar",
				Email:    "danar@gmail.com",
				Password: "",
			}

			err := signupUsecase.ValidateRequest(req)
			require.NotNil(t, err)

		})

		t.Run("Password less than 8 character", func(t *testing.T) {
			req := &models.SignUpRequest{
				Name:     "Danar",
				Email:    "danar@gmail.com",
				Password: "12312",
			}

			err := signupUsecase.ValidateRequest(req)
			require.NotNil(t, err)

		})

		t.Run("Email already exists", func(t *testing.T) {
			req := &models.SignUpRequest{
				Name:     "Danar",
				Email:    "danar@gmail.com",
				Password: "12345678",
			}
			user := &entity.User{
				Id:       "njie-2ecw21",
				Name:     "Danar",
				Email:    "danar@gmail.com",
				Password: "12345678",
			}
			userRepositoryMock.Mock.On("FindOneByEmail", req.Email).Return(user, nil)
			err := signupUsecase.ValidateRequest(req)
			userRepositoryMock.Mock.AssertCalled(t, "FindOneByEmail", req.Email)
			require.NotNil(t, err)

		})
	})

	t.Run("Hashing password", func(t *testing.T) {
		hashedPassword, err := signupUsecase.HashPassword("password")
		require.Nil(t, err)
		require.NotNil(t, hashedPassword)
	})

}
