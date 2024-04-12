package middleware

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go-crud/internal/entity"
	"go-crud/internal/repository"
	"gorm.io/gorm"
)

type ProductMiddleware struct {
	ProductRepository repository.ProductRepositoryInterface
	Log               *logrus.Logger
}

func NewProductMiddleware(productRepository repository.ProductRepositoryInterface, Log *logrus.Logger) *ProductMiddleware {
	return &ProductMiddleware{
		ProductRepository: productRepository,
		Log:               Log,
	}
}

func (m *ProductMiddleware) ProductAuth(ctx *fiber.Ctx) error {
	productID := ctx.Params("id", "")
	userID := ctx.Locals("user_id").(string)
	product := new(entity.Product)
	err := m.ProductRepository.FindOneById(product, productID)
	if err != nil {
		m.Log.WithError(err).Error("Error while finding product by id")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Product not found")
		}

		return fiber.NewError(fiber.StatusInternalServerError, "Something Wrong")
	}

	if userID != product.UserId {
		return fiber.NewError(fiber.StatusForbidden, "You're not allowed to update/delete this resource")
	}

	return ctx.Next()

}
