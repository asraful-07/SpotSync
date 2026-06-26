package dto

// ── Shared ────────────────────────────────────────────────────────────────────

// APIResponse is the standard envelope for every endpoint.
type APIResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

// ── Create / Admin list ───────────────────────────────────────────────────────

// ReservationResponse is returned by Create and the admin GET /reservations list.
// It includes flat user_id / zone_id fields plus optional nested User and Zone.
type ReservationResponse struct {
	ID           uint              `json:"id"`
	UserID       uint              `json:"user_id"`
	ZoneID       uint              `json:"zone_id"`
	LicensePlate string            `json:"license_plate"`
	Status       string            `json:"status"`
	User         *UserSummary      `json:"user,omitempty"`
	Zone         *ZoneSummary      `json:"zone,omitempty"`
	CreatedAt    string            `json:"created_at"`
	UpdatedAt    string            `json:"updated_at"`
}

// ── GET /my-reservations ─────────────────────────────────────────────────────

// MyReservationResponse is the shape returned for the authenticated user's own list.
// Spec requires a nested zone object and omits user_id.
type MyReservationResponse struct {
	ID           uint         `json:"id"`
	LicensePlate string       `json:"license_plate"`
	Status       string       `json:"status"`
	Zone         *ZoneSummary `json:"zone"`
	CreatedAt    string       `json:"created_at"`
}

// ── Nested summaries ──────────────────────────────────────────────────────────

type ZoneSummary struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type UserSummary struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}