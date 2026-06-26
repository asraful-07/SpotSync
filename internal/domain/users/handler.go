package users

import (
	"SpotSync/internal/domain/users/dto"
	"SpotSync/internal/http_response"
	"net/http"

	"github.com/labstack/echo/v5"
)

type handler struct {
	service *service
}

func NewHandler(service *service) *handler {
	return &handler{service: service}
}

func (h *handler) CreateUser(c *echo.Context) error {
	var req dto.RegisterRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, http_response.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Details: err.Error(),
		})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, http_response.Error{
			Code:    http.StatusBadRequest,
			Message: "Validation failed",
			Details: err.Error(),
		})
	}

	response, err := h.service.CreateUser(req)
	if err != nil {
		// Handle duplicate email with 409 Conflict, not 500
		if err.Error() == "email already exists" {
			return c.JSON(http.StatusConflict, http_response.Error{
				Code:    http.StatusConflict,
				Message: "Email already registered",
				Details: err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, http_response.Error{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create user",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, dto.APIResponse[*dto.RegisterResponse]{
		Success: true,
		Message: "User registered successfully",
		Data:    response,
	})
}

func (h *handler) LoginUser(c *echo.Context) error {
	var req dto.LoginRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, http_response.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Details: err.Error(),
		})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, http_response.Error{
			Code:    http.StatusBadRequest,
			Message: "Validation failed",
			Details: err.Error(),
		})
	}

	response, err := h.service.LoginUser(req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, http_response.Error{
			Code:    http.StatusUnauthorized,
			Message: "Invalid credentials",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.APIResponse[*dto.LoginResponse]{
		Success: true,
		Message: "Login successful",
		Data:    response,
	})
}

func (h *handler) GetMe(c *echo.Context) error {
	// Safely obtain user_id which may be stored as different numeric types
	var userID uint
	switch v := c.Get("user_id").(type) {
	case uint:
		userID = v
	case int:
		userID = uint(v)
	case int64:
		userID = uint(v)
	case float64:
		userID = uint(v)
	default:
		return c.JSON(http.StatusUnauthorized, http_response.Error{
			Code:    http.StatusUnauthorized,
			Message: "Cannot get user information",
			Details: "missing or invalid user id in context",
		})
	}

	// Read optional string claims safely to avoid panics when nil
	name, _ := c.Get("user_name").(string)
	email, _ := c.Get("user_email").(string)
	role, _ := c.Get("user_role").(string)
	phone, _ := c.Get("user_phone").(string)
	createdAt, _ := c.Get("created_at").(string)
    updatedAt, _ := c.Get("updated_at").(string)

	return c.JSON(http.StatusOK, dto.UserResponse{
		ID:        userID,
		Name:      name,
		Email:     email,
		Role:      role,
		Phone:     phone,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	})
}