package usecase

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go-crud/internal/entity"
	"go-crud/internal/helper"
	"go-crud/internal/models"
	"go-crud/internal/repository"
	"gorm.io/gorm"
)

type ProductUsecase struct {
	Repository repository.ProductRepositoryInterface
	Validate   *validator.Validate
	Viper      *viper.Viper
	Log        *logrus.Logger
}

func NewProductUsecase(repository repository.ProductRepositoryInterface, validate *validator.Validate, viper *viper.Viper, log *logrus.Logger) *ProductUsecase {
	return &ProductUsecase{
		Repository: repository,
		Validate:   validate,
		Viper:      viper,
		Log:        log,
	}
}

func (c *ProductUsecase) ValidateRequest(req *models.ProductRequest) error {
	err := c.Validate.Struct(req)
	if err != nil {
		c.Log.WithError(err).Error("Error validating request")
		message := helper.GetFirstValidationErrorAndConvert(err)
		return &models.ErrorResponse{
			Code:    400,
			Status:  "Bad Request",
			Message: message,
		}
	}
	return nil
}

func (c *ProductUsecase) CreateProduct(request *models.ProductRequest, userId string) (*models.ProductResponse, error) {
	err := c.ValidateRequest(request)
	if err != nil {
		c.Log.WithError(err).Error("Error validating request")
		if e, ok := err.(*models.ErrorResponse); ok {
			return nil, e
		}

		return nil, &models.ErrorResponse{
			Code:    500,
			Message: "Something Wrong",
			Status:  "Internal Server Error",
		}
	}

	var product entity.Product
	product.Id = uuid.New().String()
	product.Name = request.Name
	product.Stock = request.Stock
	product.Price = request.Price
	product.UserId = userId
	err = c.Repository.Save(&product)
	if err != nil {
		c.Log.WithError(err).Error("Error while creating product")
		return nil, &models.ErrorResponse{
			Code:    500,
			Message: "Something Wrong",
			Status:  "Internal Server Error",
		}
	}

	return &models.ProductResponse{Id: product.Id, Name: product.Name, Price: product.Price, Stock: product.Stock}, nil
}

func (c *ProductUsecase) UpdateProduct(request *models.ProductRequest, productId string) (*models.ProductResponse, error) {
	err := c.ValidateRequest(request)
	if err != nil {
		if e, ok := err.(*models.ErrorResponse); ok {
			return nil, e
		}
		c.Log.WithError(err).Error("Error while validating request")
		return nil, &models.ErrorResponse{
			Code:    500,
			Message: "Something Error",
			Status:  "Internal Server Error",
		}

	}

	product := entity.Product{
		Name:  request.Name,
		Price: request.Price,
		Stock: request.Stock,
	}

	result, err := c.Repository.UpdateById(product, productId)
	if err != nil {
		c.Log.WithError(err).Error("Error while updating product")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &models.ErrorResponse{
				Code:    404,
				Message: "Product not found",
				Status:  "Not Found",
			}
		}

		return nil, &models.ErrorResponse{
			Code:    500,
			Message: "Something Wrong",
			Status:  "Internal Server Error",
		}
	}

	return &models.ProductResponse{
		Id:    result.Id,
		Name:  result.Name,
		Stock: result.Stock,
		Price: result.Price,
	}, nil
}
