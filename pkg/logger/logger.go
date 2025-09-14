package logger

import (
	"go.uber.org/zap"
)

// Global logger instance
var Log *zap.Logger

// Initialize sets up the logger configuration
func Initialize(env string) error {
	var config zap.Config

	if env == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	var err error
	Log, err = config.Build()
	if err != nil {
		return err
	}

	return nil
}

// TODO: Structured logging helpers
func Debug(msg string, fields ...zap.Field) {
	Log.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	Log.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Log.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Log.Fatal(msg, fields...)
}
