package routes

import (
	"github.com/gofiber/fiber/v2"
	"go-crud/internal/delivery/http/controllers"
	"go-crud/internal/delivery/http/middleware"
)

type ProductRoute struct {
	App               *fiber.App
	ProductController *controllers.ProductController
	AuthMiddleware    *middleware.AuthMiddleware
}

func NewProductRoute(app *fiber.App, productController *controllers.ProductController, authMiddleware *middleware.AuthMiddleware) *ProductRoute {
	return &ProductRoute{
		App:               app,
		ProductController: productController,
		AuthMiddleware:    authMiddleware,
	}
}

func (r *ProductRoute) Setup() {
	r.App.Post("/products", r.AuthMiddleware.Auth, r.ProductController.CreateProduct)
}
