package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

const debugFileName = "debug.json"

func SetupLogger() *logrus.Logger {
	var logger = logrus.New()

	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.DebugLevel)

	logFile, err := os.OpenFile(debugFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.SetOutput(os.Stderr)
		logger.Fatalf("Error opening file: %v", err)
	}

	logger.SetOutput(logFile)

	return logger
}
