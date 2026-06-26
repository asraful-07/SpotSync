package users

import (
	"SpotSync/internal/auth"
	"SpotSync/internal/config"
	"SpotSync/internal/middlewares"
	"time"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	userRepo := NewRepository(db)
	jwtService := auth.NewJWTService(cfg.JWTSecretKey, 24*time.Hour)
	userService := NewService(userRepo, jwtService)
	userHandler := NewHandler(userService)

	api := e.Group("/api/v1/auth")

	api.POST("/register", userHandler.CreateUser)
	api.POST("/login", userHandler.LoginUser)
	api.GET("/me", userHandler.GetMe,  middlewares.AuthMiddleware(jwtService))
}