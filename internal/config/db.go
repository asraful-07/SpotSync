package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(cfg *Config) *gorm.DB { 
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{TranslateError: true})
	if err != nil {
		log.Fatal("Failed to connect to database")
	} else {
		fmt.Println("Database connection successful")
	}
	return db
}