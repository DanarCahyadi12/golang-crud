package test_e2e

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/require"
	"go-crud/internal/models"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSignup(t *testing.T) {
	t.Run("Signup with empty name", func(t *testing.T) {
		body := models.SignUpRequest{
			Name:     "",
			Email:    "danar@gmail.com",
			Password: "password",
		}
		bodyJson, err := json.Marshal(body)
		require.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(string(bodyJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		response, err := App.Fiber.Test(req)
		require.Nil(t, err)

		var expectedResponse models.ErrorResponse
		bodyByte, err := io.ReadAll(response.Body)
		require.Nil(t, err)
		err = json.Unmarshal(bodyByte, &expectedResponse)

		require.Nil(t, err)

		require.Equal(t, "Name required", expectedResponse.Message)
		require.Equal(t, "Bad Request", expectedResponse.Status)
		require.Equal(t, http.StatusBadRequest, response.StatusCode)

	})

	t.Run("Signup with empty email", func(t *testing.T) {
		body := models.SignUpRequest{
			Name:     "Danar",
			Email:    "",
			Password: "password",
		}
		bodyJson, err := json.Marshal(body)
		require.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(string(bodyJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		response, err := App.Fiber.Test(req)
		require.Nil(t, err)

		var expectedResponse models.ErrorResponse
		bodyByte, err := io.ReadAll(response.Body)
		require.Nil(t, err)
		err = json.Unmarshal(bodyByte, &expectedResponse)

		require.Nil(t, err)

		require.Equal(t, "Email required", expectedResponse.Message)
		require.Equal(t, "Bad Request", expectedResponse.Status)
		require.Equal(t, http.StatusBadRequest, response.StatusCode)

	})

	t.Run("Signup with invalid email", func(t *testing.T) {
		body := models.SignUpRequest{
			Name:     "Danar",
			Email:    "danar@invalid",
			Password: "password",
		}
		bodyJson, err := json.Marshal(body)
		require.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(string(bodyJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		response, err := App.Fiber.Test(req)
		require.Nil(t, err)

		var expectedResponse models.ErrorResponse
		bodyByte, err := io.ReadAll(response.Body)
		require.Nil(t, err)
		err = json.Unmarshal(bodyByte, &expectedResponse)

		require.Nil(t, err)

		require.Equal(t, "Email is invalid", expectedResponse.Message)
		require.Equal(t, "Bad Request", expectedResponse.Status)
		require.Equal(t, http.StatusBadRequest, response.StatusCode)

	})

	t.Run("Signup with password less than 8 character", func(t *testing.T) {
		body := models.SignUpRequest{
			Name:     "Danar",
			Email:    "danar@gmail.com",
			Password: "123",
		}
		bodyJson, err := json.Marshal(body)
		require.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(string(bodyJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		response, err := App.Fiber.Test(req)
		require.Nil(t, err)

		var expectedResponse models.ErrorResponse
		bodyByte, err := io.ReadAll(response.Body)
		require.Nil(t, err)
		err = json.Unmarshal(bodyByte, &expectedResponse)

		require.Nil(t, err)

		require.Equal(t, "Password min 8 character", expectedResponse.Message)
		require.Equal(t, "Bad Request", expectedResponse.Status)
		require.Equal(t, http.StatusBadRequest, response.StatusCode)

	})

	t.Run("Signup with empty password", func(t *testing.T) {
		body := models.SignUpRequest{
			Name:     "Danar",
			Email:    "danar@gmail.com",
			Password: "",
		}
		bodyJson, err := json.Marshal(body)
		require.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(string(bodyJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		response, err := App.Fiber.Test(req)
		require.Nil(t, err)

		var expectedResponse models.ErrorResponse
		bodyByte, err := io.ReadAll(response.Body)
		require.Nil(t, err)
		err = json.Unmarshal(bodyByte, &expectedResponse)

		require.Nil(t, err)

		require.Equal(t, "Password required", expectedResponse.Message)
		require.Equal(t, "Bad Request", expectedResponse.Status)
		require.Equal(t, http.StatusBadRequest, response.StatusCode)

	})

	t.Run("Signup without error", func(t *testing.T) {
		body := models.SignUpRequest{
			Name:     "Danar",
			Email:    "danar@gmail.com",
			Password: "12345678",
		}
		bodyJson, err := json.Marshal(body)
		require.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(string(bodyJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		user, err := UserRepository.FindOneByEmail(body.Email)
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			require.NotEqual(t, gorm.ErrRecordNotFound, err)
		}
		if user != nil {
			err := UserRepository.DeleteOneById(user.Id)
			require.Nil(t, err)
		}
		response, err := App.Fiber.Test(req)
		require.Nil(t, err)

		var expectedResponse models.Response[*models.SignUpResponse]
		bodyByte, err := io.ReadAll(response.Body)
		require.Nil(t, err)
		err = json.Unmarshal(bodyByte, &expectedResponse)

		require.Nil(t, err)

		require.NotNil(t, expectedResponse.Data.Id)
		require.Equal(t, body.Name, expectedResponse.Data.Name)
		require.NotNil(t, expectedResponse.Data.CreatedAt)
		require.NotNil(t, expectedResponse.Data.UpdatedAt)
		require.Equal(t, http.StatusCreated, response.StatusCode)
		err = UserRepository.DeleteOneById(expectedResponse.Data.Id)
		require.Nil(t, err)
	})
}
