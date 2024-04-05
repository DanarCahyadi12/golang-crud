package test_e2e

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/require"
	"go-crud/internal/entity"
	"go-crud/internal/models"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var user = &entity.User{
	Name:     "Danar",
	Email:    "danar@gmail.com",
	Password: "$2a$10$aOySpFRuA2uE8gGNNCuAleiBvNRyMJpZuyhZ21kf/Tpy5c8KHNRTe",
}

func createUser() {
	userFound, err := UserRepository.FindOneByEmail(user.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	}

	if userFound != nil {
		err := UserRepository.DeleteOneById(userFound.Id)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			panic(err)
		}
	}

	err = UserRepository.Save(user)
	if err != nil {
		panic(err)
	}
}

func deleteUser() {
	err := UserRepository.DeleteOneById(user.Id)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	}
}

func TestAuth(t *testing.T) {
	createUser()
	defer deleteUser()

	t.Run("Sign in with empty email", func(t *testing.T) {
		body := models.SignInRequest{
			Email:    "danar@gmail.com",
			Password: "12345678",
		}

		bodyJson, err := json.Marshal(body)
		require.Nil(t, err)
		req := httptest.NewRequest(http.MethodPost, "/auth/signin", strings.NewReader(string(bodyJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		response, err := App.Fiber.Test(req)

		require.Nil(t, err)
		var expectedResponse models.ErrorResponse
		bodyByte, err := io.ReadAll(response.Body)
		require.Nil(t, err)

		err = json.Unmarshal(bodyByte, &expectedResponse)
		require.Nil(t, err)

		require.Equal(t, http.StatusUnauthorized, response.StatusCode)
		require.Equal(t, "Email or password is incorrect", expectedResponse.Message)
		require.Equal(t, "Unauthorized", expectedResponse.Status)

	})

	t.Run("Sign in with invalid email format", func(t *testing.T) {
		body := models.SignInRequest{
			Email:    "danar@invalid",
			Password: "12345678",
		}

		bodyJson, err := json.Marshal(body)
		require.Nil(t, err)
		req := httptest.NewRequest(http.MethodPost, "/auth/signin", strings.NewReader(string(bodyJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		response, err := App.Fiber.Test(req)

		require.Nil(t, err)
		var expectedResponse models.ErrorResponse
		bodyByte, err := io.ReadAll(response.Body)
		require.Nil(t, err)

		err = json.Unmarshal(bodyByte, &expectedResponse)
		require.Nil(t, err)

		require.Equal(t, http.StatusBadRequest, response.StatusCode)
		require.Equal(t, "Email format is invalid", expectedResponse.Message)
		require.Equal(t, "Bad Request", expectedResponse.Status)

	})

	t.Run("Sign in with empty password", func(t *testing.T) {
		body := models.SignInRequest{
			Email:    "danar@gmail.com",
			Password: "",
		}

		bodyJson, err := json.Marshal(body)
		require.Nil(t, err)
		req := httptest.NewRequest(http.MethodPost, "/auth/signin", strings.NewReader(string(bodyJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		response, err := App.Fiber.Test(req)

		require.Nil(t, err)
		var expectedResponse models.ErrorResponse
		bodyByte, err := io.ReadAll(response.Body)
		require.Nil(t, err)

		err = json.Unmarshal(bodyByte, &expectedResponse)
		require.Nil(t, err)

		require.Equal(t, http.StatusBadRequest, response.StatusCode)
		require.Equal(t, "Password required", expectedResponse.Message)
		require.Equal(t, "Bad Request", expectedResponse.Status)

	})

	t.Run("Signin with incorrect email", func(t *testing.T) {
		body := models.SignInRequest{
			Email:    "invalid@gmail.com",
			Password: "12345678",
		}

		bodyJson, err := json.Marshal(body)
		require.Nil(t, err)
		req := httptest.NewRequest(http.MethodPost, "/auth/signin", strings.NewReader(string(bodyJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		response, err := App.Fiber.Test(req)

		require.Nil(t, err)
		var expectedResponse models.ErrorResponse
		bodyByte, err := io.ReadAll(response.Body)
		require.Nil(t, err)

		err = json.Unmarshal(bodyByte, &expectedResponse)
		require.Nil(t, err)

		require.Equal(t, http.StatusUnauthorized, response.StatusCode)
		require.Equal(t, "Email or password is incorrect", expectedResponse.Message)
		require.Equal(t, "Unauthorized", expectedResponse.Status)

	})

	t.Run("Signin with incorrect password", func(t *testing.T) {
		body := models.SignInRequest{
			Email:    "danar@gmail.com",
			Password: "invalid",
		}

		bodyJson, err := json.Marshal(body)
		require.Nil(t, err)
		req := httptest.NewRequest(http.MethodPost, "/auth/signin", strings.NewReader(string(bodyJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		response, err := App.Fiber.Test(req)

		require.Nil(t, err)
		var expectedResponse models.ErrorResponse
		bodyByte, err := io.ReadAll(response.Body)
		require.Nil(t, err)

		err = json.Unmarshal(bodyByte, &expectedResponse)
		require.Nil(t, err)

		require.Equal(t, http.StatusUnauthorized, response.StatusCode)
		require.Equal(t, "Email or password is incorrect", expectedResponse.Message)
		require.Equal(t, "Unauthorized", expectedResponse.Status)

	})

	t.Run("Signin with correct credentials", func(t *testing.T) {
		body := models.SignInRequest{
			Email:    "danar@gmail.com",
			Password: "12345678",
		}

		bodyJson, err := json.Marshal(body)
		require.Nil(t, err)
		req := httptest.NewRequest(http.MethodPost, "/auth/signin", strings.NewReader(string(bodyJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		response, err := App.Fiber.Test(req)

		require.Nil(t, err)
		var expectedResponse models.Response[*models.SignInResponse]
		bodyByte, err := io.ReadAll(response.Body)
		require.Nil(t, err)

		err = json.Unmarshal(bodyByte, &expectedResponse)
		require.Nil(t, err)

		require.Equal(t, http.StatusOK, response.StatusCode)
		require.Equal(t, "Signin successfully", expectedResponse.Message)
		require.NotNil(t, expectedResponse.Data.AccessToken)
		require.NotNil(t, expectedResponse.Data.RefreshToken)

		cookieHeader := response.Header.Get("Set-Cookie")
		cookieValue := strings.Split(strings.Split(cookieHeader, ";")[0], "=")[1]
		require.NotNil(t, cookieValue)

	})

}
