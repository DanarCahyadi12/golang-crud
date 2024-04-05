package test

import (
	"github.com/stretchr/testify/require"
	"go-crud/internal/entity"
	"go-crud/internal/models"
	"go-crud/internal/usecase"
	"sync"
	"testing"
	"time"
)

func TestAuth(t *testing.T) {
	authUsecase := usecase.NewAuthUsecase(userRepositoryMock, validate, viperConfig, log)
	t.Run("Validate request", func(t *testing.T) {
		req := &models.SignInRequest{
			Email:    "",
			Password: "12345678",
		}
		err := authUsecase.ValidateRequest(req)
		require.NotNil(t, err)
		require.Equal(t, &models.ErrorResponse{Code: 400, Message: "Email required", Status: "Bad Request"}, err)

	})

	t.Run("Generate access token", func(t *testing.T) {
		token, err := authUsecase.GenerateAccessToken("my-id")
		require.Nil(t, err)
		require.NotNil(t, token)
	})

	t.Run("Generate refresh token", func(t *testing.T) {
		token, err := authUsecase.GenerateRefreshToken("my-id")
		require.Nil(t, err)
		require.NotNil(t, token)
	})

	t.Run("Generate access and refresh token using goroutine", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			accessToken, err := authUsecase.GenerateAccessToken("my-id")
			require.Nil(t, err)
			require.NotNil(t, accessToken)
		}()

		go func() {
			defer wg.Done()
			refreshToken, err := authUsecase.GenerateRefreshToken("my-id")
			require.Nil(t, err)
			require.NotNil(t, refreshToken)
		}()

		wg.Wait()
	})

	t.Run("Sign in with invalid email", func(t *testing.T) {
		request := &models.SignInRequest{
			Email:    "danar@invalid",
			Password: "12345678",
		}

		result, err := authUsecase.SignIn(request)
		require.Nil(t, result)
		require.Equal(t, &models.ErrorResponse{
			Code:    400,
			Message: "Email is invalid",
			Status:  "Bad Request",
		}, err)
	})

	t.Run("Sign in with empty password", func(t *testing.T) {
		request := &models.SignInRequest{
			Email:    "danar@gmail.com",
			Password: "",
		}

		result, err := authUsecase.SignIn(request)
		require.Nil(t, result)
		require.Equal(t, &models.ErrorResponse{
			Code:    400,
			Message: "Password required",
			Status:  "Bad Request",
		}, err)
	})

	t.Run("Signin with valid credentials", func(t *testing.T) {
		user := &entity.User{
			Id:        "my-id",
			Name:      "Danar Cahyadi",
			Email:     "danar@gmail.com",
			Password:  "$2a$10$aOySpFRuA2uE8gGNNCuAleiBvNRyMJpZuyhZ21kf/Tpy5c8KHNRTe",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		request := &models.SignInRequest{
			Email:    "danar@gmail.com",
			Password: "12345678",
		}
		userRepositoryMock.Mock.On("FindOneByEmail", request.Email).Return(user, nil)
		result, err := authUsecase.SignIn(request)
		require.Nil(t, err)
		require.NotNil(t, result)
	})

	t.Run("Signin with invalid email", func(t *testing.T) {
		request := &models.SignInRequest{
			Email:    "invalid@gmail.com",
			Password: "12345678",
		}
		userRepositoryMock.Mock.On("FindOneByEmail", request.Email).Return(nil, nil)
		result, err := authUsecase.SignIn(request)
		require.Nil(t, result)
		require.Equal(t, &models.ErrorResponse{
			Code:    401,
			Message: "Email or password is incorrect",
			Status:  "Unauthorized",
		}, err)

	})

	t.Run("Signin with invalid password", func(t *testing.T) {
		user := &entity.User{
			Id:        "my-id",
			Name:      "Danar Cahyadi",
			Email:     "danar@gmail.com",
			Password:  "$2a$10$aOySpFRuA2uE8gGNNCuAleiBvNRyMJpZuyhZ21kf/Tpy5c8KHNRTe",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		request := &models.SignInRequest{
			Email:    "danar@gmail.com",
			Password: "wrongpassword",
		}
		userRepositoryMock.Mock.On("FindOneByEmail", request.Email).Return(user, nil)
		result, err := authUsecase.SignIn(request)
		require.Nil(t, result)
		require.Equal(t, &models.ErrorResponse{
			Code:    401,
			Message: "Email or password is incorrect",
			Status:  "Unauthorized",
		}, err)

	})
}
