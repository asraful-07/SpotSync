package parking_zones

import (
	"SpotSync/internal/domain/parking_zones/dto"
	"SpotSync/internal/http_response"
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

type handler struct {
	service Service
}

func NewHandler(s Service) *handler {
	return &handler{service: s}
}

func zoneErrorResponse(c *echo.Context, err error) error {
	if errors.Is(err, ErrParkingZoneNotFound) {
		return c.JSON(http.StatusNotFound, http_response.Error{
			Code:    http.StatusNotFound,
			Message: "Parking zone not found",
		})
	}
	return c.JSON(http.StatusInternalServerError, http_response.Error{
		Code:    http.StatusInternalServerError,
		Message: "Something went wrong",
		Details: err.Error(),
	})
}

func (h *handler) CreateParkingZone(c *echo.Context) error {
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

	response, err := h.service.CreateParkingZone(&req)
	if err != nil {
		return zoneErrorResponse(c, err)
	}

	return c.JSON(http.StatusCreated, dto.APIResponse[*dto.ParkingZoneResponse]{
		Success: true,
		Message: "Parking zone created successfully",
		Data:    response,
	})
}

func (h *handler) GetAllParkingZones(c *echo.Context) error {
	zones, err := h.service.GetAllParkingZone()
	if err != nil {
		return zoneErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, dto.APIResponse[[]*dto.ParkingZoneResponse]{
		Success: true,
		Message: "Parking zones retrieved successfully",
		Data:    zones,
	})
}

func (h *handler) GetParkingZoneByID(c *echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, http_response.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid zone ID",
		})
	}

	zone, err := h.service.GetByIDParkingZone(uint(id))
	if err != nil {
		return zoneErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, dto.APIResponse[*dto.ParkingZoneResponse]{
		Success: true,
		Message: "Parking zone retrieved successfully",
		Data:    zone,
	})
}