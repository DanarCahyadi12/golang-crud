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
	ProductMiddleware *middleware.ProductMiddleware
}

func NewProductRoute(app *fiber.App, productController *controllers.ProductController, authMiddleware *middleware.AuthMiddleware, productMiddleware *middleware.ProductMiddleware) *ProductRoute {
	return &ProductRoute{
		App:               app,
		ProductController: productController,
		AuthMiddleware:    authMiddleware,
		ProductMiddleware: productMiddleware,
	}
}

func (r *ProductRoute) Setup() {
	r.App.Post("/products", r.AuthMiddleware.Auth, r.ProductController.CreateProduct)
	r.App.Put("/products/:id", r.AuthMiddleware.Auth, r.ProductMiddleware.ProductAuth, r.ProductController.UpdateProduct)
}
