package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(cfg *Config) *gorm.DB {
	if cfg == nil || cfg.DSN == "" {
		log.Fatal("Database DSN is empty")
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  cfg.DSN,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		TranslateError: true,
		PrepareStmt:    false,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	} else {
		fmt.Println("Database connection successful")
	}
	return db
}