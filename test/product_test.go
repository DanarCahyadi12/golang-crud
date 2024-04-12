package test

import (
	"github.com/stretchr/testify/require"
	"go-crud/internal/entity"
	"go-crud/internal/models"
	"go-crud/internal/usecase"
	"gorm.io/gorm"
	"testing"
)

func TestProduct(t *testing.T) {
	productUsecase := usecase.NewProductUsecase(productRepositoryMock, validate, viperConfig, log)
	t.Run("Validate request", func(t *testing.T) {
		t.Run("Empty name", func(t *testing.T) {
			req := &models.ProductRequest{
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
			req := &models.ProductRequest{
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

	})

	t.Run("Update product", func(t *testing.T) {
		t.Run("Should successfully updating product", func(t *testing.T) {
			const productID string = "product-id"
			request := new(models.ProductRequest)
			request.Name = "Product update"
			request.Stock = 10
			request.Price = 1500

			product := new(entity.Product)
			product.Name = request.Name
			product.Stock = request.Stock
			product.Price = request.Price

			productRepositoryMock.Mock.On("UpdateById", entity.Product{Name: request.Name, Stock: request.Stock, Price: request.Price}, productID).Return(product, nil)
			result, err := productUsecase.UpdateProduct(request, productID)
			require.Nil(t, err)
			require.Equal(t, "Product update", result.Name)
			require.Equal(t, 10, result.Stock)
			require.Equal(t, 1500, product.Price)
		})

		t.Run("Should return a error when updating with empty name", func(t *testing.T) {
			const productID string = "product-id"
			request := new(models.ProductRequest)
			request.Name = ""
			request.Stock = 10
			request.Price = 1500

			product := new(entity.Product)
			product.Name = request.Name
			product.Stock = request.Stock
			product.Price = request.Price

			productRepositoryMock.Mock.On("UpdateById", product, productID).Return(product, nil)
			result, err := productUsecase.UpdateProduct(request, productID)
			require.NotNil(t, err)
			require.Nil(t, result)
		})

		t.Run("Should return a error when updating with wrong id", func(t *testing.T) {
			const productID string = "wrong-id"
			request := new(models.ProductRequest)
			request.Name = "Product updated"
			request.Stock = 10
			request.Price = 1500

			product := entity.Product{
				Name:  request.Name,
				Stock: request.Stock,
				Price: request.Price,
			}

			productRepositoryMock.Mock.On("UpdateById", product, productID).Return(nil, gorm.ErrRecordNotFound)
			result, err := productUsecase.UpdateProduct(request, productID)
			require.NotNil(t, err)
			require.Nil(t, result)
		})
	})
}
