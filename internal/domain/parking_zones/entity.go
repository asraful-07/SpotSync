package parking_zones

import (
	"SpotSync/internal/domain/parking_zones/dto"

	"gorm.io/gorm"
)

type ParkingZone struct {
	gorm.Model
	Name           string  `json:"name" gorm:"type:varchar(255);not null"`
	Type           string  `json:"type" gorm:"type:varchar(100);not null"`
	TotalCapacity  int     `json:"total_capacity" gorm:"not null"`
	PricePerHour   float64 `json:"price_per_hour" gorm:"not null"`
	AvailableSpots int `json:"available_spots" gorm:"->;-:migration"`
}

	
func (p *ParkingZone) ToResponse() *dto.ParkingZoneResponse {
	return &dto.ParkingZoneResponse{
		ID:             p.ID,
		Name:           p.Name,
		Type:           p.Type,
		TotalCapacity:  p.TotalCapacity,
		AvailableSpots: p.AvailableSpots,
		PricePerHour:   p.PricePerHour,
		CreatedAt:      p.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
		UpdatedAt:      p.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z"),
	}
}