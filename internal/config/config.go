package config

import (
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/novaru/billing-service/pkg/logger"
)

type Config struct {
	Env         string
	Port        string
	DatabaseURL string
	JWTSecret   string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Error loading .env file", zap.Error(err))
	}

	return &Config{
		Env:         os.Getenv("ENV"),
		Port:        os.Getenv("PORT"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
	}
}
