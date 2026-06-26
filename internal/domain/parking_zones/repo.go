package parking_zones

import (
	"errors"

	"gorm.io/gorm"
)

var ErrParkingZoneNotFound = errors.New("parking zone not found")

type Repository interface {
	Create(zone *ParkingZone) error
	GetAll() ([]*ParkingZone, error)
	GetByID(id uint) (*ParkingZone, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(zone *ParkingZone) error {
	return r.db.Create(zone).Error
}

func (r *repository) GetAll() ([]*ParkingZone, error) {
	var zones []*ParkingZone

	err := r.db.
		Select("parking_zones.*").
		Find(&zones).Error

	if err != nil {
		return nil, err
	}

	return zones, nil
}

func (r *repository) GetByID(id uint) (*ParkingZone, error) {
	var zone ParkingZone

	err := r.db.
		Select("parking_zones.*").
		Where("id = ?", id).
		First(&zone).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrParkingZoneNotFound
		}
		return nil, err
	}

	return &zone, nil
}