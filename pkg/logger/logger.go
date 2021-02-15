package logger

import (
	"os"

	"go.uber.org/zap"
)

func Create() {
	var logger *zap.Logger

	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "error"
	}

	switch level {
	case "debug":
		logger, _ = zap.NewDevelopment()
	case "error":
		logger, _ = zap.NewProduction()
	default:
		logger, _ = zap.NewProduction()
	}
	logger.Sync()
	zap.ReplaceGlobals(logger)
}
