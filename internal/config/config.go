package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DSN          string
	PORT         string
	JWTSecretKey string
}

func LoadConfig() (*Config, error) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env file")
	}

	return &Config{
		DSN:          os.Getenv("DSN"),
		PORT:         os.Getenv("PORT"),
		JWTSecretKey: os.Getenv("JWT_SECRET_KEY"),
	}, nil
}