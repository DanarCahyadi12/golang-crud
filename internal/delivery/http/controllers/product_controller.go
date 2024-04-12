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
	var body = new(models.ProductRequest)
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

func (c *ProductController) UpdateProduct(ctx *fiber.Ctx) error {
	request := new(models.ProductRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.WithError(err).Error("Error while unmarshal request body")
		if e, ok := err.(*fiber.UnmarshalTypeError); ok {
			return fiber.NewError(400, fmt.Sprintf("Invalid type %s type for %s field", e.Type, e.Field))
		}

		return fiber.NewError(fiber.StatusInternalServerError, "Something Wrong")

	}

	productID := ctx.Params("id")
	result, err := c.ProductUsecase.UpdateProduct(request, productID)
	if err != nil {
		if e, ok := err.(*models.ErrorResponse); ok {
			return fiber.NewError(e.Code, e.Message)
		}
		c.Log.WithError(err).Error("Unknown error while updating product")
		return fiber.NewError(fiber.StatusInternalServerError, "Something Wrong")
	}

	return ctx.Status(fiber.StatusOK).JSON(&models.Response[*models.ProductResponse]{
		Message: "Product updated",
		Data:    result,
	})

}

func (c *ProductController) DeleteProduct(ctx *fiber.Ctx) error {
	productID := ctx.Params("id")
	err := c.ProductUsecase.DeleteProduct(productID)
	if err != nil {
		if e, ok := err.(*models.ErrorResponse); ok {
			return fiber.NewError(e.Code, e.Message)
		}
		c.Log.WithError(err).Error("Unknown error while deleting product")
		return fiber.NewError(fiber.StatusInternalServerError, "Something Wrong")

	}
	return ctx.Status(fiber.StatusOK).JSON(&models.Response[any]{
		Message: "Product deleted",
	})
}
