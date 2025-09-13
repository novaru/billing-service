package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/novaru/billing-service/pkg/logger"
)

type Config struct {
	DatabaseURL string
	Port        string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Error loading .env file", err)
	}

	return &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        os.Getenv("PORT"),
	}
}
