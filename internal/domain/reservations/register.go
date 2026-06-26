package reservations

import (
	"time"

	"SpotSync/internal/auth"
	"SpotSync/internal/config"
	"SpotSync/internal/domain/parking_zones"
	"SpotSync/internal/middlewares"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	// Repositories
	zoneRepo := parking_zones.NewRepository(db)
	reservationRepo := NewRepository(db, zoneRepo) // zoneRepo used inside tx

	// JWT service (same secret & TTL as auth module)
	jwtService := auth.NewJWTService(cfg.JWTSecretKey, 24*time.Hour)

	// Service and handler
	reservationService := NewService(reservationRepo)
	h := NewHandler(reservationService)

	api := e.Group("/api/v1/reservations")

	// ── Authenticated routes (driver + admin)
	protected := api.Group("")
	protected.Use(middlewares.AuthMiddleware(jwtService))

	protected.POST("", h.CreateReservation)                     
	protected.GET("/my-reservations", h.GetMyReservations)       
	protected.DELETE("/:id", h.CancelReservation)              

	// ── Admin-only routes
	admin := api.Group("")
	admin.Use(middlewares.AuthMiddleware(jwtService))
	admin.Use(middlewares.AdminOnly()) 

	admin.GET("", h.GetAllReservations) 
}