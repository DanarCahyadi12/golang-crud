package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go-crud/internal/models"
	"go-crud/internal/usecase"
)

type ProductController struct {
	Log            *logrus.Logger
	ProductUsecase *usecase.ProductUsecase
}

func NewProductController(log *logrus.Logger, usecase *usecase.ProductUsecase) *ProductController {
	return &ProductController{
		Log:            log,
		ProductUsecase: usecase,
	}
}

func (c *ProductController) CreateProduct(ctx *fiber.Ctx) error {
	var body = new(models.CreateProductRequest)
	userId := ctx.Locals("user_id").(string)
	err := ctx.BodyParser(body)
	if err != nil {
		c.Log.WithError(err).Error("Error while parsing body request")
		if e, ok := err.(*fiber.UnmarshalTypeError); ok {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("%s must be a string", e.Field))
		}

		return fiber.NewError(fiber.StatusInternalServerError, "Something Error")
	}

	result, err := c.ProductUsecase.CreateProduct(body, userId)
	if err != nil {
		c.Log.WithError(err).Error("Error while creating product")
		if e, ok := err.(*models.ErrorResponse); ok {
			return fiber.NewError(e.Code, e.Message)
		}

		return fiber.NewError(fiber.StatusInternalServerError, "Something Error")

	}

	return ctx.Status(fiber.StatusCreated).JSON(&models.Response[*models.ProductResponse]{
		Message: "Product created!",
		Data:    result,
	})
}
