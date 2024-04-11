package test_e2e

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"go-crud/internal/models"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestProduct(t *testing.T) {
	CreateUser(user)
	response, err := SignIn(models.SignInRequest{
		Email:    user.Email,
		Password: "12345678",
	})
	fmt.Println(response)
	require.Nil(t, err)
	defer DeleteUser(user.Id)

	t.Run("Create product", func(t *testing.T) {
		body := models.CreateProductRequest{
			Name:  "Product 1",
			Price: 1500,
			Stock: 20000,
		}

		bodyJson, err := json.Marshal(body)
		require.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(string(bodyJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", response.Data.AccessToken))
		response, err := App.Fiber.Test(req)
		require.Nil(t, err)

		var expectedResponse models.Response[*models.ProductResponse]
		bodyByte, err := io.ReadAll(response.Body)
		require.Nil(t, err)
		err = json.Unmarshal(bodyByte, &expectedResponse)
		require.Nil(t, err)

		require.Equal(t, http.StatusCreated, response.StatusCode)
		require.Equal(t, "Product created!", expectedResponse.Message)
		require.NotNil(t, expectedResponse.Data.Id)
		require.Equal(t, expectedResponse.Data.Name, expectedResponse.Data.Name)
		require.Equal(t, expectedResponse.Data.Stock, expectedResponse.Data.Stock)
		require.Equal(t, expectedResponse.Data.Price, expectedResponse.Data.Price)
	})

	t.Run("Create product with empty product name", func(t *testing.T) {
		body := models.CreateProductRequest{
			Name:  "",
			Price: 1500,
			Stock: 20000,
		}

		bodyJson, err := json.Marshal(body)
		require.Nil(t, err)

		req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(string(bodyJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", response.Data.AccessToken))
		response, err := App.Fiber.Test(req)
		require.Nil(t, err)

		var expectedResponse models.ErrorResponse
		bodyByte, err := io.ReadAll(response.Body)
		err = json.Unmarshal(bodyByte, &expectedResponse)
		require.Nil(t, err)

		require.Equal(t, http.StatusBadRequest, response.StatusCode)
		require.Equal(t, "Name required", expectedResponse.Message)
		require.Equal(t, "Bad Request", expectedResponse.Status)

	})
}
