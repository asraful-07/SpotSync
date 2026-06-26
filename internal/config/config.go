package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DSN          string
	PORT         string
	JWTSecretKey string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("load .env file: %w", err)
		}
	}

	cfg := &Config{
		DSN:          os.Getenv("DSN"),
		PORT:         os.Getenv("PORT"),
		JWTSecretKey: os.Getenv("JWT_SECRET_KEY"),
	}

	if cfg.DSN == "" {
		return nil, fmt.Errorf("DSN is required")
	}
	if cfg.PORT == "" {
		cfg.PORT = "8080"
	}
	if cfg.JWTSecretKey == "" {
		return nil, fmt.Errorf("JWT_SECRET_KEY is required")
	}

	return cfg, nil
}