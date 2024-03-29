package config

import "github.com/sirupsen/logrus"

func NewLogrus() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	return logger
}
