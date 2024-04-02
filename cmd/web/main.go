package main

import (
	"go-crud/internal/config"
)

func main() {
	viper := config.NewViper("./../../")
	fiber := config.NewFiber()
	logrus := config.NewLogrus()
	validator := config.NewValidator()
	database := config.NewGorm(viper)

	app := config.NewApp(fiber, validator, database, viper, logrus)
	app.Setup()
	app.StartServer()

}
