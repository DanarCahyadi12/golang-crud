package config

import "github.com/sirupsen/logrus"

func NewLogrus() *logrus.Logger {
	logger := logrus.New()
	return logger
}
