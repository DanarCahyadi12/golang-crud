package routes

import (
	"github.com/gofiber/fiber/v2"
	"go-crud/internal/delivery/http/controllers"
)

type SignupRoute struct {
	App              *fiber.App
	SignupController *controllers.SignupController
}

func NewSignupRoute(app *fiber.App, controller *controllers.SignupController) *SignupRoute {
	return &SignupRoute{
		App:              app,
		SignupController: controller,
	}
}

func (r *SignupRoute) Setup() {
	r.App.Post("/signup", r.SignupController.Signup)
}
