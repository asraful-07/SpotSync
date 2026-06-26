package dto

// Register Request
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Phone    string `json:"phone,omitempty" validate:"omitempty"`
	Role  string `json:"role" validate:"required,oneof=driver admin"`
	// Role     string `json:"role,omitempty" validate:"omitempty,oneof=driver admin"`
}

// Login Request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}