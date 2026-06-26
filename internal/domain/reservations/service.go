package reservations

import (
	"SpotSync/internal/domain/reservations/dto"
)

// Service defines the business-logic contract for the reservations domain.
type Service interface {
	// Create validates, enforces capacity, and persists a new reservation.
	// userID comes from the JWT claim — never from the request body.
	Create(userID uint, req *dto.CreateRequest) (*dto.ReservationResponse, error)

	// GetMyReservations returns all reservations belonging to the given user.
	GetMyReservations(userID uint) ([]*dto.MyReservationResponse, error)

	// Cancel marks a reservation as cancelled.
	// Returns ErrForbidden if the reservation belongs to a different user.
	Cancel(reservationID uint, requesterUserID uint) error

	// GetAll returns every reservation in the system (admin only).
	GetAll() ([]*dto.ReservationResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// ── Implementations

func (s *service) Create(userID uint, req *dto.CreateRequest) (*dto.ReservationResponse, error) {
	reservation := &Reservation{
		UserID:       userID,
		ZoneID:       req.ZoneID,
		LicensePlate: req.LicensePlate,
		Status:       ReservationActive,
	}

	// Delegates concurrency-safe capacity enforcement to the repository.
	if err := s.repo.CreateWithCapacityCheck(reservation); err != nil {
		return nil, err
	}

	return reservation.ToResponse(), nil
}

func (s *service) GetMyReservations(userID uint) ([]*dto.MyReservationResponse, error) {
	reservations, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.MyReservationResponse, 0, len(reservations))
	for _, r := range reservations {
		responses = append(responses, r.ToMyResponse())
	}
	return responses, nil
}

func (s *service) Cancel(reservationID uint, requesterUserID uint) error {
	return s.repo.Cancel(reservationID, requesterUserID)
}

func (s *service) GetAll() ([]*dto.ReservationResponse, error) {
	reservations, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.ReservationResponse, 0, len(reservations))
	for _, r := range reservations {
		responses = append(responses, r.ToResponse())
	}
	return responses, nil
}