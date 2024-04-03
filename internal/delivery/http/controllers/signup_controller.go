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

func NewSignupController(usecase *usecase.SignUpUsecase) *SignupController {
	return &SignupController{
		SignupUsecase: usecase,
	}
}

func (c *SignupController) Signup(ctx *fiber.Ctx) error {
	request := new(models.SignUpRequest)
	fmt.Println("REQUEST: ", request)
	err := ctx.BodyParser(request)

	if err != nil {
		c.Log.Errorf("Error while parsing request body. %v", err)
		return fiber.NewError(500, "Something Wrong")
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
