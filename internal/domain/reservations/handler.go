package reservations

import (
	"errors"
	"net/http"
	"strconv"

	"SpotSync/internal/domain/reservations/dto"
	"SpotSync/internal/http_response"

	"github.com/labstack/echo/v5"
)

type handler struct {
	service Service
}

func NewHandler(s Service) *handler {
	return &handler{service: s}
}

//  Error helper

func reservationErrorResponse(c *echo.Context, err error) error {
	switch {
	case errors.Is(err, ErrReservationNotFound):
		return c.JSON(http.StatusNotFound, http_response.Error{
			Code:    http.StatusNotFound,
			Message: "Reservation not found",
		})
	case errors.Is(err, ErrZoneNotFound):
		return c.JSON(http.StatusNotFound, http_response.Error{
			Code:    http.StatusNotFound,
			Message: "Parking zone not found",
		})
	case errors.Is(err, ErrZoneFull):
		return c.JSON(http.StatusConflict, http_response.Error{
			Code:    http.StatusConflict,
			Message: "Parking zone is at full capacity",
		})
	case errors.Is(err, ErrForbidden):
		return c.JSON(http.StatusForbidden, http_response.Error{
			Code:    http.StatusForbidden,
			Message: "You do not have permission to modify this reservation",
		})
	default:
		return c.JSON(http.StatusInternalServerError, http_response.Error{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong",
			Details: err.Error(),
		})
	}
}

// ── Handlers ──────────────────────────────────────────────────────────────────

// CreateReservation POST /api/v1/reservations
// Access: authenticated users (driver, admin)
func (h *handler) CreateReservation(c *echo.Context) error {
	// Extract user ID injected by AuthMiddleware into the Echo context.
	userID, ok := c.Get("user_id").(uint)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, http_response.Error{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	var req dto.CreateRequest
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

	reservation, err := h.service.Create(userID, &req)
	if err != nil {
		return reservationErrorResponse(c, err)
	}

	return c.JSON(http.StatusCreated, dto.APIResponse[*dto.ReservationResponse]{
		Success: true,
		Message: "Reservation confirmed successfully",
		Data:    reservation,
	})
}

// GetMyReservations GET /api/v1/reservations/my-reservations
// Access: authenticated users
func (h *handler) GetMyReservations(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, http_response.Error{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	reservations, err := h.service.GetMyReservations(userID)
	if err != nil {
		return reservationErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, dto.APIResponse[[]*dto.MyReservationResponse]{
		Success: true,
		Message: "My reservations retrieved successfully",
		Data:    reservations,
	})
}

// CancelReservation DELETE /api/v1/reservations/:id
// Access: authenticated users — drivers can only cancel their own reservations.
func (h *handler) CancelReservation(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, http_response.Error{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, http_response.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid reservation ID",
		})
	}

	if err := h.service.Cancel(uint(id), userID); err != nil {
		return reservationErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, dto.APIResponse[any]{
		Success: true,
		Message: "Reservation cancelled successfully",
	})
}

// GetAllReservations GET /api/v1/reservations
// Access: admin only
func (h *handler) GetAllReservations(c *echo.Context) error {
	reservations, err := h.service.GetAll()
	if err != nil {
		return reservationErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, dto.APIResponse[[]*dto.ReservationResponse]{
		Success: true,
		Message: "All reservations retrieved successfully",
		Data:    reservations,
	})
}