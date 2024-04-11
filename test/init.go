package test

import (
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go-crud/internal/config"
	"go-crud/test/mocks"
)

var viperConfig *viper.Viper
var userRepositoryMock *mocks.UserRepositoryMock
var productRepositoryMock *mocks.ProductRepositoryMock
var validate *validator.Validate
var log *logrus.Logger

func init() {
	viperConfig = config.NewViper("./../")
	userRepositoryMock = mocks.NewRepositoryMock()
	productRepositoryMock = mocks.NewProductRepositoryMock()
	validate = config.NewValidator()
	log = config.NewLogrus()
}
