package dto

// CreateRequest — UserID is NOT here; it is extracted from the JWT in the handler.
type CreateRequest struct {
	ZoneID       uint   `json:"zone_id" validate:"required"`
	LicensePlate string `json:"license_plate" validate:"required,max=15"`
}

// UpdateStatusRequest is kept for potential admin use (e.g. marking completed).
type UpdateStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=active completed cancelled"`
}