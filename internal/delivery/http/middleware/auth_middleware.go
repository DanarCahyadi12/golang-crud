package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go-crud/internal/models"
	"go-crud/internal/usecase"
)

type AuthMiddleware struct {
	AuthUsecase *usecase.AuthUsecase
	Log         *logrus.Logger
}

func NewAuthMiddleware(authUsecase *usecase.AuthUsecase, log *logrus.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		AuthUsecase: authUsecase,
		Log:         log,
	}
}

func (m *AuthMiddleware) Auth(ctx *fiber.Ctx) error {
	accessToken, err := m.AuthUsecase.ParseTokenFromHeader(ctx.Get("Authorization"))

	if err != nil {
		m.Log.WithError(err).Error("Error while parsing token from header")
		if e, ok := err.(*models.ErrorResponse); ok {
			return fiber.NewError(fiber.StatusUnauthorized, e.Message)
		} else {
			return fiber.NewError(fiber.StatusInternalServerError, "Something Error")
		}

	}

	sub, err := m.AuthUsecase.VerifyAccessToken(accessToken)
	if err != nil {
		m.Log.WithError(err).Error("Error while verifying access token")
		if e, ok := err.(*models.ErrorResponse); ok {
			return fiber.NewError(e.Code, e.Message)
		} else {
			return fiber.NewError(fiber.StatusInternalServerError, "Something Error")
		}
	}

	ctx.Locals("user_id", sub)
	return ctx.Next()
}
