package test

import (
	"github.com/stretchr/testify/require"
	"go-crud/internal/entity"
	"go-crud/internal/models"
	"go-crud/internal/usecase"
	"testing"
)

func TestProduct(t *testing.T) {
	productUsecase := usecase.NewProductUsecase(productRepositoryMock, validate, viperConfig, log)
	t.Run("Validate request", func(t *testing.T) {
		t.Run("Empty name", func(t *testing.T) {
			req := &models.CreateProductRequest{
				Name:  "",
				Price: 12000,
				Stock: 302,
			}
			err := productUsecase.ValidateRequest(req)
			require.NotNil(t, err)
			require.Equal(t, "Name required", err.Error())
		})
	})

	t.Run("Create product", func(t *testing.T) {
		t.Run("Should return error when creating with empty name", func(t *testing.T) {
			req := &models.CreateProductRequest{
				Name:  "",
				Price: 12000,
				Stock: 302,
			}

			result, err := productUsecase.CreateProduct(req, "user-id")

			expectedResponse := new(models.ErrorResponse)
			expectedResponse.Code = 400
			expectedResponse.Message = "Name required"
			expectedResponse.Status = "Bad Request"
			require.Equal(t, expectedResponse, err)
			require.Nil(t, result)
		})

		t.Run("Should successfully creating product", func(t *testing.T) {
			req := &models.CreateProductRequest{
				Name:  "Product 1",
				Price: 12000,
				Stock: 302,
			}
			product := new(entity.Product)
			product.Name = "Product 1"
			product.Stock = 302
			product.Price = 12000
			product.UserId = "user-id"

			expectedResponse := new(models.ProductResponse)
			expectedResponse.Name = product.Name
			expectedResponse.Stock = product.Stock
			expectedResponse.Price = product.Price

			productRepositoryMock.Mock.On("Save", product).Return(nil)
			result, err := productUsecase.CreateProduct(req, "user-id")
			require.Nil(t, err)
			require.Equal(t, expectedResponse, result)

		})
	})
}
