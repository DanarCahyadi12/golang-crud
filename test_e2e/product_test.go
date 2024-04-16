package test_e2e

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go-crud/internal/entity"
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
		body := models.ProductRequest{
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
		body := models.ProductRequest{
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

	t.Run("Update product", func(t *testing.T) {
		product := new(entity.Product)
		product.Id = uuid.New().String()
		product.Name = "Product 2"
		product.Price = 15000
		product.Stock = 120
		product.UserId = user.Id

		err := ProductRepository.Save(product)
		require.Nil(t, err)

		body := models.ProductRequest{
			Name:  "Updated product",
			Price: 1500,
			Stock: 20000,
		}

		bodyJson, err := json.Marshal(body)
		require.Nil(t, err)

		target := fmt.Sprintf("/products/%s", product.Id)
		req := httptest.NewRequest(http.MethodPut, target, strings.NewReader(string(bodyJson)))
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

		require.Equal(t, http.StatusOK, response.StatusCode)
		require.Equal(t, "Product updated", expectedResponse.Message)
		fmt.Println(expectedResponse.Data)
		require.NotNil(t, expectedResponse.Data.Id)
		require.Equal(t, expectedResponse.Data.Name, expectedResponse.Data.Name)
		require.Equal(t, expectedResponse.Data.Stock, expectedResponse.Data.Stock)
		require.Equal(t, expectedResponse.Data.Price, expectedResponse.Data.Price)
	})

	t.Run("Update product with empty name", func(t *testing.T) {
		product := new(entity.Product)
		product.Id = uuid.New().String()
		product.Name = "Product 2"
		product.Price = 15000
		product.Stock = 120
		product.UserId = user.Id

		err := ProductRepository.Save(product)
		require.Nil(t, err)

		body := models.ProductRequest{
			Name:  "",
			Price: 1500,
			Stock: 20000,
		}

		bodyJson, err := json.Marshal(body)
		require.Nil(t, err)

		target := fmt.Sprintf("/products/%s", product.Id)
		req := httptest.NewRequest(http.MethodPut, target, strings.NewReader(string(bodyJson)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", response.Data.AccessToken))
		response, err := App.Fiber.Test(req)
		require.Nil(t, err)

		var expectedResponse models.ErrorResponse
		bodyByte, err := io.ReadAll(response.Body)
		require.Nil(t, err)
		err = json.Unmarshal(bodyByte, &expectedResponse)
		require.Nil(t, err)

		require.Equal(t, http.StatusBadRequest, response.StatusCode)

	})

	t.Run("Should delete product", func(t *testing.T) {
		product := new(entity.Product)
		product.Id = uuid.New().String()
		product.Name = "Product 2"
		product.Price = 15000
		product.Stock = 120
		product.UserId = user.Id

		err := ProductRepository.Save(product)
		require.Nil(t, err)

		target := fmt.Sprintf("/products/%s", product.Id)
		req := httptest.NewRequest(http.MethodDelete, target, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", response.Data.AccessToken))
		response, err := App.Fiber.Test(req)
		require.Nil(t, err)

		var actualResponse models.Response[any]
		bodyByte, err := io.ReadAll(response.Body)
		require.Nil(t, err)
		err = json.Unmarshal(bodyByte, &actualResponse)
		require.Nil(t, err)

		require.Equal(t, http.StatusOK, response.StatusCode)
		require.Equal(t, "Product deleted", actualResponse.Message)
	})

	t.Run("Should return error if product is not found", func(t *testing.T) {
		product := new(entity.Product)
		product.Id = uuid.New().String()
		product.Name = "Product 2"
		product.Price = 15000
		product.Stock = 120
		product.UserId = user.Id

		err := ProductRepository.Save(product)
		require.Nil(t, err)

		target := fmt.Sprintf("/products/%s", "notfound-id")
		req := httptest.NewRequest(http.MethodDelete, target, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", response.Data.AccessToken))
		response, err := App.Fiber.Test(req)
		require.Nil(t, err)

		var actualResponse models.ErrorResponse
		bodyByte, err := io.ReadAll(response.Body)
		require.Nil(t, err)
		err = json.Unmarshal(bodyByte, &actualResponse)
		require.Nil(t, err)

		require.Equal(t, http.StatusNotFound, response.StatusCode)
		require.Equal(t, "Product not found", actualResponse.Message)
	})

	t.Run("Should return products metadata (page 1 and limit 10)", func(t *testing.T) {
		const PAGE = 1
		const LIMIT = 10
		target := fmt.Sprintf("/products?page=%d&limit=%d", PAGE, LIMIT)
		err := CreateManyProduct(user.Id)
		require.Nil(t, err)
		req := httptest.NewRequest(http.MethodGet, target, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", response.Data.AccessToken))

		response, err := App.Fiber.Test(req)
		require.Nil(t, err)

		var actualResponse models.Response[*[]models.ProductResponse]
		bodyByte, err := io.ReadAll(response.Body)
		require.Nil(t, err)
		err = json.Unmarshal(bodyByte, &actualResponse)
		require.Nil(t, err)

		require.NotNil(t, actualResponse.Message)
		require.Equal(t, int64(5), actualResponse.Metadata.PageSize)
		require.Equal(t, int64(50), actualResponse.Metadata.TotalItemCount)
		require.Equal(t, 1, actualResponse.Metadata.PageNumber)
		require.Equal(t, "", actualResponse.Metadata.Prev)
		require.Equal(t, fmt.Sprintf("http://localhost:8080/products?page=%d&limit=%d", PAGE+1, LIMIT), actualResponse.Metadata.Next)
		require.Equal(t, "", actualResponse.Metadata.Prev)

	})

	t.Run("Should return products metadata (page 2 and limit 10)", func(t *testing.T) {
		const PAGE = 2
		const LIMIT = 10
		target := fmt.Sprintf("/products?page=%d&limit=%d", PAGE, LIMIT)
		err := CreateManyProduct(user.Id)
		require.Nil(t, err)
		req := httptest.NewRequest(http.MethodGet, target, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", response.Data.AccessToken))

		response, err := App.Fiber.Test(req)
		require.Nil(t, err)

		var actualResponse models.Response[*[]models.ProductResponse]
		bodyByte, err := io.ReadAll(response.Body)
		require.Nil(t, err)
		err = json.Unmarshal(bodyByte, &actualResponse)
		require.Nil(t, err)

		require.NotNil(t, actualResponse.Message)
		require.Equal(t, int64(5), actualResponse.Metadata.PageSize)
		require.Equal(t, int64(50), actualResponse.Metadata.TotalItemCount)
		require.Equal(t, PAGE, actualResponse.Metadata.PageNumber)
		require.Equal(t, fmt.Sprintf("http://localhost:8080/products?page=%d&limit=%d", PAGE+1, LIMIT), actualResponse.Metadata.Next)
		require.Equal(t, fmt.Sprintf("http://localhost:8080/products?page=%d&limit=%d", PAGE-1, LIMIT), actualResponse.Metadata.Prev)

	})

	t.Run("Should return products metadata (page 5 and limit 10)", func(t *testing.T) {
		const PAGE = 5
		const LIMIT = 10
		target := fmt.Sprintf("/products?page=%d&limit=%d", PAGE, LIMIT)
		err := CreateManyProduct(user.Id)
		require.Nil(t, err)
		req := httptest.NewRequest(http.MethodGet, target, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", response.Data.AccessToken))

		response, err := App.Fiber.Test(req)
		require.Nil(t, err)

		var actualResponse models.Response[*[]models.ProductResponse]
		bodyByte, err := io.ReadAll(response.Body)
		require.Nil(t, err)
		err = json.Unmarshal(bodyByte, &actualResponse)
		require.Nil(t, err)

		require.NotNil(t, actualResponse.Message)
		require.Equal(t, int64(5), actualResponse.Metadata.PageSize)
		require.Equal(t, int64(50), actualResponse.Metadata.TotalItemCount)
		require.Equal(t, PAGE, actualResponse.Metadata.PageNumber)
		require.Equal(t, "", actualResponse.Metadata.Next)
		require.Equal(t, fmt.Sprintf("http://localhost:8080/products?page=%d&limit=%d", PAGE-1, LIMIT), actualResponse.Metadata.Prev)

	})

	t.Run("Get detail product", func(t *testing.T) {
		product, err := CreateOneProduct(user.Id)
		require.Nil(t, err)

		target := fmt.Sprintf("/products/%s", product.Id)
		req := httptest.NewRequest(http.MethodGet, target, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", response.Data.AccessToken))
		response, err := App.Fiber.Test(req)
		require.Nil(t, err)

		var actualResponse models.Response[*models.ProductResponse]
		bodyByte, err := io.ReadAll(response.Body)
		require.Nil(t, err)

		err = json.Unmarshal(bodyByte, &actualResponse)
		require.Nil(t, err)

		require.NotNil(t, actualResponse.Data.Id)
		require.Equal(t, product.Name, actualResponse.Data.Name)
		require.Equal(t, product.Stock, actualResponse.Data.Stock)
		require.Equal(t, product.Price, actualResponse.Data.Price)
	})

}
