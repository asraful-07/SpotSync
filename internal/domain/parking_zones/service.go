package parking_zones

import (
	"SpotSync/internal/domain/parking_zones/dto"
)

// Service interface — enables mocking in tests
type Service interface {
	CreateParkingZone(req *dto.CreateRequest) (*dto.ParkingZoneResponse, error)
	GetAllParkingZone() ([]*dto.ParkingZoneResponse, error)
	GetByIDParkingZone(id uint) (*dto.ParkingZoneResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateParkingZone(req *dto.CreateRequest) (*dto.ParkingZoneResponse, error) {
	zone := &ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}

	if err := s.repo.Create(zone); err != nil {
		return nil, err
	}

	// On creation, no reservations exist yet — available = total
	zone.AvailableSpots = zone.TotalCapacity
	return zone.ToResponse(), nil
}

func (s *service) GetAllParkingZone() ([]*dto.ParkingZoneResponse, error) {
	zones, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.ParkingZoneResponse, 0, len(zones))
	for _, zone := range zones {
		responses = append(responses, zone.ToResponse())
	}
	return responses, nil
}

func (s *service) GetByIDParkingZone(id uint) (*dto.ParkingZoneResponse, error) {
	zone, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return zone.ToResponse(), nil
}