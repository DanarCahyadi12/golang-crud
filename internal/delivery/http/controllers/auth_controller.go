package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go-crud/internal/models"
	"go-crud/internal/usecase"
	"time"
)

type AuthController struct {
	Log         *logrus.Logger
	AuthUsecase *usecase.AuthUsecase
}

func NewAuthController(log *logrus.Logger, authUsecase *usecase.AuthUsecase) *AuthController {
	return &AuthController{
		Log:         log,
		AuthUsecase: authUsecase,
	}
}

func (c *AuthController) SignIn(ctx *fiber.Ctx) error {
	ctx.Accepts("application/json")
	req := new(models.SignInRequest)
	err := ctx.BodyParser(req)

	if err != nil {
		c.Log.WithError(err).Error(err)
		if e, ok := err.(*fiber.UnmarshalTypeError); ok {
			return fiber.NewError(400, fmt.Sprintf("%s must be string", e.Field))
		}

		return fiber.NewError(500, "Something Error")

	}

	result, err := c.AuthUsecase.SignIn(req)
	if err != nil {
		if e, ok := err.(*models.ErrorResponse); ok {
			return fiber.NewError(e.Code, e.Message)
		}

		c.Log.Errorf("%v", err)
		return fiber.NewError(500, "Something Error")
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "refresh_token"
	cookie.Value = result.RefreshToken
	cookie.Expires = time.Now().Add(3 * (24 * time.Hour))
	cookie.HTTPOnly = true

	ctx.Cookie(cookie)

	return ctx.Status(fiber.StatusOK).JSON(models.Response[*models.SignInResponse]{
		Message: "Signin successfully",
		Data:    result,
	})

}
