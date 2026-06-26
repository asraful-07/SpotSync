package dto

type ParkingZoneResponse struct {
	ID             uint    `json:"id"`
	Name           string  `json:"name"`
	Type           string  `json:"type"`
	TotalCapacity  int     `json:"total_capacity"`
	AvailableSpots int     `json:"available_spots"`
	PricePerHour   float64 `json:"price_per_hour"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

type APIResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}
