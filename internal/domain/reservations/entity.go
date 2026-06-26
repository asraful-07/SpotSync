package reservations

import (
	"SpotSync/internal/domain/parking_zones"
	"SpotSync/internal/domain/reservations/dto"
	"SpotSync/internal/domain/users"

	"gorm.io/gorm"
)

const (
	ReservationActive    = "active"
	ReservationCompleted = "completed"
	ReservationCancelled = "cancelled"
)

type Reservation struct {
	gorm.Model
	UserID       uint                      `json:"user_id" gorm:"not null;index"`
	User         users.User                `json:"user,omitempty" gorm:"foreignKey:UserID"`
	ZoneID       uint                      `json:"zone_id" gorm:"not null;index"`
	Zone         parking_zones.ParkingZone `json:"zone,omitempty" gorm:"foreignKey:ZoneID"`
	LicensePlate string                    `json:"license_plate" gorm:"type:varchar(15);not null"`
	Status       string                    `json:"status" gorm:"type:varchar(20);default:'active';not null;index"`
}

// ToResponse builds the flat response used by Create and the admin list.
func (r *Reservation) ToResponse() *dto.ReservationResponse {
	resp := &dto.ReservationResponse{
		ID:           r.ID,
		UserID:       r.UserID,
		ZoneID:       r.ZoneID,
		LicensePlate: r.LicensePlate,
		Status:       r.Status,
		CreatedAt:    r.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    r.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z"),
	}

	// Populate nested User summary only when the association was preloaded.
	if r.User.ID != 0 {
		resp.User = &dto.UserSummary{
			ID:    r.User.ID,
			Name:  r.User.Name,
			Email: r.User.Email,
		}
	}

	// Populate nested Zone summary only when the association was preloaded.
	if r.Zone.ID != 0 {
		resp.Zone = &dto.ZoneSummary{
			ID:   r.Zone.ID,
			Name: r.Zone.Name,
			Type: r.Zone.Type,
		}
	}

	return resp
}

// ToMyResponse builds the user-facing "my reservations" shape from the spec.
func (r *Reservation) ToMyResponse() *dto.MyReservationResponse {
	resp := &dto.MyReservationResponse{
		ID:           r.ID,
		LicensePlate: r.LicensePlate,
		Status:       r.Status,
		CreatedAt:    r.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
	}

	if r.Zone.ID != 0 {
		resp.Zone = &dto.ZoneSummary{
			ID:   r.Zone.ID,
			Name: r.Zone.Name,
			Type: r.Zone.Type,
		}
	}

	return resp
}