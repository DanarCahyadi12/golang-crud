package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go-crud/internal/models"
	"go-crud/internal/usecase"
)

type SignupController struct {
	Log           *logrus.Logger
	SignupUsecase *usecase.SignUpUsecase
}

func NewSignupController(log *logrus.Logger, usecase *usecase.SignUpUsecase) *SignupController {
	return &SignupController{
		Log:           log,
		SignupUsecase: usecase,
	}
}

func (c *SignupController) Signup(ctx *fiber.Ctx) error {
	ctx.Accepts("application/json")
	request := new(models.SignUpRequest)
	err := ctx.BodyParser(request)

	if err != nil {
		c.Log.Errorf("%v", err)
		if e, ok := err.(*fiber.UnmarshalTypeError); ok {
			return fiber.NewError(400, fmt.Sprintf("%s must be a string", e.Field))

		}
	}

	result, err := c.SignupUsecase.CreateUser(request)
	if err != nil {
		if e, ok := err.(*models.ErrorResponse); ok {
			return fiber.NewError(e.Code, e.Message)
		}
		c.Log.Errorf("Error while creating user: %v", err)
		return fiber.NewError(500, "Something Error")
	}

	return ctx.Status(fiber.StatusCreated).JSON(&models.Response[*models.SignUpResponse]{
		Message: "Signup successfully",
		Data:    result,
	})

}
