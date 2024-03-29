package config

import (
	"fmt"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type App struct {
	Fiber     *fiber.App
	Validator *validator.Validate
	Database  *gorm.DB
	Viper     *viper.Viper
	Logger    *logrus.Logger
}

func NewApp(fiber *fiber.App, validator *validator.Validate, database *gorm.DB, viper *viper.Viper, logger *logrus.Logger) *App {
	return &App{
		Fiber:     fiber,
		Validator: validator,
		Database:  database,
		Viper:     viper,
		Logger:    logger,
	}
}

func (app *App) StartServer() {
	port := app.Viper.GetInt("web.port")
	err := app.Fiber.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

}

func (app *App) Setup() {

}
