package test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go-crud/internal/entity"
	"go-crud/internal/helper"
	"go-crud/internal/models"
	"go-crud/internal/usecase"
	"gorm.io/gorm"
	"math"
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

	t.Run("Delete product", func(t *testing.T) {
		const productID = "product-id"
		t.Run("Should return error if the product is doesn't matched", func(t *testing.T) {
			productRepositoryMock.Mock.On("DeleteById", "").Return(gorm.ErrRecordNotFound)
			err := productUsecase.DeleteProduct("")
			require.Equal(t, "Product not found", err.Error())
		})

		t.Run("Shouldn't return a error", func(t *testing.T) {
			productRepositoryMock.Mock.On("DeleteById", productID).Return(nil)
			err := productUsecase.DeleteProduct(productID)
			require.Nil(t, err)
		})
	})

	t.Run("Get products", func(t *testing.T) {
		expectedResult := &[]models.ProductResponse{
			{
				Id:    "1",
				Name:  "Product 1",
				Price: 15000,
				Stock: 120,
				User: models.UserResponse{
					Id:   "user-id-1",
					Name: "Danar Cahyadi",
				},
			},
			{
				Id:    "2",
				Name:  "Product 2",
				Price: 20000,
				Stock: 150,
				User: models.UserResponse{
					Id:   "user-id-2",
					Name: "Ketut Danar",
				},
			},
		}

		productMock := []entity.Product{
			{
				Id:    "1",
				Name:  "Product 1",
				Price: 15000,
				Stock: 120,
				User: entity.User{
					Id:   "user-id-1",
					Name: "Danar Cahyadi",
				},
			},
			{
				Id:    "2",
				Name:  "Product 2",
				Price: 20000,
				Stock: 150,
				User: entity.User{
					Id:   "user-id-2",
					Name: "Ketut Danar",
				},
			},
		}

		t.Run("Should return products with user entity", func(t *testing.T) {
			productRepositoryMock.Mock.On("FindMany", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
				productsPtr := args.Get(0).(*[]entity.Product)
				*productsPtr = productMock
			})
			result, err := productUsecase.GetProducts(0, 2)
			require.Nil(t, err)
			require.Equal(t, expectedResult, result)

		})

		t.Run("Should return next and previous pagination URL", func(t *testing.T) {
			nextURL := helper.FormatNextURLPagination("products", 1, 10, 20)
			assert.Equal(t, "http://localhost:8080/products?page=2&limit=10", nextURL)

			nextURL = helper.FormatNextURLPagination("products", 20, 10, 20)
			assert.Equal(t, "", nextURL)

			prevURL := helper.FormatPrevURLPagination("products", 2, 10)
			assert.Equal(t, "http://localhost:8080/products?page=1&limit=10", prevURL)

			prevURL = helper.FormatPrevURLPagination("products", 1, 10)
			assert.Equal(t, "", prevURL)
		})

		t.Run("Should return metadata", func(t *testing.T) {
			var returnArgs int64 = 500
			var pageSize int64 = int64(math.Ceil(float64(returnArgs / 50)))
			productRepositoryMock.Mock.On("Count").Return(returnArgs, nil)
			metadata, err := productUsecase.GetMetadataPagination(1, 50)
			require.Nil(t, err)
			require.Equal(t, 1, metadata.PageNumber)
			require.Equal(t, returnArgs, metadata.TotalItemCount)
			require.Equal(t, "http://localhost:8080/products?page=2&limit=50", metadata.Next)
			require.Equal(t, "", metadata.Prev)
			require.Equal(t, pageSize, metadata.PageSize)
		})

		t.Run("Next link & previous link on metadata must not be empty", func(t *testing.T) {
			var count int64 = 500
			size := float64(count) / float64(50)
			pageSize := int64(math.Ceil(size))
			productRepositoryMock.Mock.On("Count").Return(count, nil)
			metadata, err := productUsecase.GetMetadataPagination(5, 50)
			require.Nil(t, err)
			require.Equal(t, 5, metadata.PageNumber)
			require.Equal(t, count, metadata.TotalItemCount)
			require.Equal(t, "http://localhost:8080/products?page=6&limit=50", metadata.Next)
			require.Equal(t, "http://localhost:8080/products?page=4&limit=50", metadata.Prev)
			require.Equal(t, pageSize, metadata.PageSize)
		})

		t.Run("Next link on metadata must be empty", func(t *testing.T) {
			var returnArgs int64 = 500
			size := float64(returnArgs) / float64(50)
			pageSize := int64(math.Ceil(size))
			productRepositoryMock.Mock.On("Count").Return(returnArgs, nil)
			metadata, err := productUsecase.GetMetadataPagination(10, 50)
			require.Nil(t, err)
			require.Equal(t, 10, metadata.PageNumber)
			require.Equal(t, returnArgs, metadata.TotalItemCount)
			require.Equal(t, "", metadata.Next)
			require.Equal(t, "http://localhost:8080/products?page=9&limit=50", metadata.Prev)
			require.Equal(t, pageSize, metadata.PageSize)
		})

	})

	t.Run("Get detail products", func(t *testing.T) {
		t.Run("Should return error 404 not found", func(t *testing.T) {
			productRepositoryMock.Mock.On("FindOneById", mock.Anything, "invalid-id").Return(gorm.ErrRecordNotFound)
			result, err := productUsecase.GetDetailProduct("invalid-id")
			require.Nil(t, result)
			require.Equal(t, "Product not found", err.Error())
		})

		t.Run("Should return detail product", func(t *testing.T) {
			var productMock entity.Product
			productMock.Id = "id"
			productMock.Name = "Product 1"
			productMock.Price = 1500
			productMock.Stock = 120
			productMock.User = entity.User{
				Id:   "user-id",
				Name: "Danar",
			}
			productRepositoryMock.Mock.On("FindOneById", mock.Anything, "id").Return(nil).Run(func(args mock.Arguments) {
				productPtr := args.Get(0).(*entity.Product)
				*productPtr = productMock
			})
			result, err := productUsecase.GetDetailProduct("id")
			require.Nil(t, err)
			require.Equal(t, productMock.Id, result.Id)
			require.Equal(t, productMock.Name, result.Name)
			require.Equal(t, productMock.Price, result.Price)
			require.Equal(t, productMock.Stock, result.Stock)
			require.Equal(t, productMock.User.Id, result.User.Id)
			require.Equal(t, productMock.User.Name, result.User.Name)
		})

	})
}
