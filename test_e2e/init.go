package test_e2e

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go-crud/internal/config"
	"go-crud/internal/repository"
	"gorm.io/gorm"
)

var ViperConfig *viper.Viper
var Validate *validator.Validate
var Database *gorm.DB
var Log *logrus.Logger
var FiberApp *fiber.App
var App *config.App
var UserRepository *repository.UserRepository
var ProductRepository *repository.ProductRepository

func init() {
	ViperConfig = config.NewViper("./../")
	Validate = config.NewValidator()
	Database = config.NewGorm(ViperConfig)
	Log = config.NewLogrus()
	FiberApp = config.NewFiber()
	UserRepository = repository.NewUserRepository(Database)
	ProductRepository = repository.NewProductRepository(Database)
	App = config.NewApp(FiberApp, Validate, Database, ViperConfig, Log)
	App.Setup()
}
