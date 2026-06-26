package reservations

import (
	"errors"

	"SpotSync/internal/domain/parking_zones"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrReservationNotFound = errors.New("reservation not found")
	ErrZoneFull            = errors.New("parking zone is at full capacity")
	ErrZoneNotFound        = errors.New("parking zone not found")
	ErrForbidden           = errors.New("you do not have permission to modify this reservation")
)

// Repository defines the data-access contract for reservations.
type Repository interface {
	// CreateWithCapacityCheck opens a transaction, locks the zone row, verifies
	// capacity, and creates the reservation atomically.
	CreateWithCapacityCheck(reservation *Reservation) error

	GetAll() ([]*Reservation, error)
	GetByID(id uint) (*Reservation, error)
	GetByUserID(userID uint) ([]*Reservation, error)

	// Cancel sets the reservation status to "cancelled".
	// It returns ErrForbidden if the reservation does not belong to ownerUserID.
	Cancel(id uint, ownerUserID uint) error
}

type repository struct {
	db      *gorm.DB
	zoneRepo parking_zones.Repository // used inside the transaction
}

func NewRepository(db *gorm.DB, zoneRepo parking_zones.Repository) Repository {
	return &repository{db: db, zoneRepo: zoneRepo}
}

// ── Write ─────────────────────────────────────────────────────────────────────

// CreateWithCapacityCheck prevents overbooking using a serialisable transaction
// with a FOR UPDATE row-level lock on the parking zone.
func (r *repository) CreateWithCapacityCheck(reservation *Reservation) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Lock the target zone row so no concurrent transaction can read or
		//    write it until this transaction commits or rolls back.
		var zone parking_zones.ParkingZone
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&zone, reservation.ZoneID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrZoneNotFound
			}
			return err
		}

		// 2. Count currently active reservations for this zone inside the
		//    same transaction so we see the locked, consistent value.
		var activeCount int64
		if err := tx.Model(&Reservation{}).
			Where("zone_id = ? AND status = ?", zone.ID, ReservationActive).
			Count(&activeCount).Error; err != nil {
			return err
		}

		// 3. Reject if the zone is already at capacity.
		if int(activeCount) >= zone.TotalCapacity {
			return ErrZoneFull
		}

		// 4. Safe to create — the lock guarantees no other goroutine can
		//    sneak a reservation in between steps 2 and 4.
		if err := tx.Create(reservation).Error; err != nil {
			return err
		}

		return nil // commits the transaction
	})
}

// Cancel sets status to "cancelled" but only for the owning user.
func (r *repository) Cancel(id uint, ownerUserID uint) error {
	// First fetch to verify ownership.
	var reservation Reservation
	if err := r.db.First(&reservation, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrReservationNotFound
		}
		return err
	}

	if reservation.UserID != ownerUserID {
		return ErrForbidden
	}

	if reservation.Status == ReservationCancelled {
		// Idempotent — already cancelled, nothing to do.
		return nil
	}

	return r.db.Model(&reservation).Update("status", ReservationCancelled).Error
}

// ── Read
func (r *repository) GetAll() ([]*Reservation, error) {
	var reservations []*Reservation
	err := r.db.
		Preload("User").
		Preload("Zone").
		Find(&reservations).Error
	return reservations, err
}

func (r *repository) GetByID(id uint) (*Reservation, error) {
	var reservation Reservation
	err := r.db.
		Preload("User").
		Preload("Zone").
		First(&reservation, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReservationNotFound
		}
		return nil, err
	}
	return &reservation, nil
}

func (r *repository) GetByUserID(userID uint) ([]*Reservation, error) {
	var reservations []*Reservation
	err := r.db.
		Preload("Zone").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&reservations).Error
	return reservations, err
}