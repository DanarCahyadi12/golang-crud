package test_e2e

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"go-crud/internal/models"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuth(t *testing.T) {
	CreateUser(user)
	defer DeleteUser(user.Id)

	t.Run("Sign in with empty email", func(t *testing.T) {
		body := models.SignInRequest{
			Email:    "",
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
		require.Equal(t, "Email required", expectedResponse.Message)
		require.Equal(t, "Bad Request", expectedResponse.Status)

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
		var expectedResponse models.Response[*models.AuthResponse]
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
