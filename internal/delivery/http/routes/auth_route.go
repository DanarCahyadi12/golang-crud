package routes

import (
	"github.com/gofiber/fiber/v2"
	"go-crud/internal/delivery/http/controllers"
)

type AuthRoute struct {
	App            *fiber.App
	AuthController *controllers.AuthController
}

func NewAuthRoute(app *fiber.App, authController *controllers.AuthController) *AuthRoute {
	return &AuthRoute{
		App:            app,
		AuthController: authController,
	}
}

func (r *AuthRoute) Setup() {
	r.App.Post("/auth/signin", r.AuthController.SignIn)
	r.App.Get("/signout", r.AuthController.SignOut)
}
