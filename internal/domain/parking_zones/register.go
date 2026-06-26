package parking_zones

import (
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc)

	api := e.Group("/api/v1/zones")


	api.POST("", handler.CreateParkingZone)
	api.GET("", handler.GetAllParkingZones)
	api.GET("/:id", handler.GetParkingZoneByID)
}