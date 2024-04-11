package injector

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go-crud/internal/delivery/http/controllers"
	"go-crud/internal/delivery/http/middleware"
	"go-crud/internal/delivery/http/routes"
	"go-crud/internal/repository"
	"go-crud/internal/usecase"
	"gorm.io/gorm"
)

var authUsecase *usecase.AuthUsecase

func InjectSignupRoute(app *fiber.App, database *gorm.DB, validator *validator.Validate, log *logrus.Logger) *routes.SignupRoute {
	userRepository := repository.NewUserRepository(database)
	signupUsecase := usecase.NewSignUpUsecase(userRepository, validator, log)
	signupController := controllers.NewSignupController(log, signupUsecase)
	signupRoute := routes.NewSignupRoute(app, signupController)

	return signupRoute

}

func InjectAuthRoute(app *fiber.App, database *gorm.DB, validator *validator.Validate, viper *viper.Viper, log *logrus.Logger) *routes.AuthRoute {
	userRepository := repository.NewUserRepository(database)
	authUsecase = usecase.NewAuthUsecase(userRepository, validator, viper, log)
	authController := controllers.NewAuthController(log, authUsecase)
	authRoute := routes.NewAuthRoute(app, authController)

	return authRoute
}

func InjectProductRoute(app *fiber.App, database *gorm.DB, validator *validator.Validate, viper *viper.Viper, log *logrus.Logger) *routes.ProductRoute {
	productRepository := repository.NewProductRepository(database)
	productUsecase := usecase.NewProductUsecase(productRepository, validator, viper, log)
	productController := controllers.NewProductController(log, productUsecase)
	authMiddleware := middleware.NewAuthMiddleware(authUsecase, log)
	productRoute := routes.NewProductRoute(app, productController, authMiddleware)

	return productRoute
}
