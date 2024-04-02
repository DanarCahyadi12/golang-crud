package injector

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go-crud/internal/delivery/http/controllers"
	"go-crud/internal/delivery/http/routes"
	"go-crud/internal/repository"
	"go-crud/internal/usecase"
	"gorm.io/gorm"
)

func InjectSignupRoute(app *fiber.App, database *gorm.DB, validator *validator.Validate, log *logrus.Logger) *routes.SignupRoute {
	userRepository := repository.NewUserRepository(database)
	userUsecase := usecase.NewSignUpUsecase(userRepository, validator, log)
	userController := controllers.NewSignupController(userUsecase)
	userRoute := routes.NewSignupRoute(app, userController)

	return userRoute

}
