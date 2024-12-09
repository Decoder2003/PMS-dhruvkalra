package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

// InitLogger initializes the logger
func InitLogger() {
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)
	log.Info("Logger initialized")
}

// Info logs an info message
func Info(message string, fields logrus.Fields) {
	log.WithFields(fields).Info(message)
}

// Error logs an error message
func Error(message string, fields logrus.Fields) {
	log.WithFields(fields).Error(message)
}
