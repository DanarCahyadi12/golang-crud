package config

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"go-crud/internal/models"
)

func errorHandlerConfig(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	var status string
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}
	switch code {
	case 400:
		status = "Bad Request"
	case 401:
		status = "Unauthorized"
	case 403:
		status = "Forbidden"
	case 404:
		status = "Not Found"
	case 408:
		status = "Request Timeout"
	case 500:
		status = "Internal Server Error"

	}
	return ctx.Status(code).JSON(models.ErrorResponse{
		Status:  status,
		Message: e.Message,
	})
}
func NewFiber() *fiber.App {
	fiber := fiber.New(fiber.Config{
		ErrorHandler: errorHandlerConfig,
	})
	return fiber
}
