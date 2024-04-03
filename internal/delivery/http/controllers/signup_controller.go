package controllers

import (
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
		return fiber.NewError(400, "Fields must be a string")
	}

	result, err := c.SignupUsecase.CreateUser(request)
	if err != nil {
		if e := err.(*models.ErrorResponse); e != nil {
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
