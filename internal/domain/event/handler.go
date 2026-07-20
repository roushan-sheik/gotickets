package event

import (
	"errors"
	"gotickets/internal/domain/event/dto"
	"gotickets/internal/httpresponse"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

type handler struct {
	service *service
}

func NewHandler(s *service) *handler {
	return &handler{service: s}
}

func eventErrorResponse(c *echo.Context, err error) error {
	if errors.Is(err, ErrEventNotFound) {
		return c.JSON(http.StatusNotFound, httpresponse.Error{
			Code:    http.StatusNotFound,
			Message: "Event not found",
		})
	}

	return c.JSON(http.StatusInternalServerError, httpresponse.Error{
		Code:    http.StatusInternalServerError,
		Message: "Something went wrong",
		Details: err.Error(),
	})
}

// CreateEvent godoc
// @Summary      Create a new event
// @Description  Creates a new event with the provided details.
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreateRequest  true  "Event Creation Details"
// @Success      201      {object}  dto.Response
// @Failure      400      {object}  httpresponse.Error
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/events [post]
func (h *handler) CreateEvent(c *echo.Context) error {
	var req dto.CreateRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Details: err.Error(),
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Validation failed",
			Details: err.Error(),
		})
	}

	response, err := h.service.CreateEvent(req)
	if err != nil {
		return eventErrorResponse(c, err)
	}

	return c.JSON(http.StatusCreated, response)
}

// GetEvents godoc
// @Summary      List all events
// @Description  Retrieves a list of all available events.
// @Tags         Events
// @Produce      json
// @Success      200      {array}   dto.Response
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/events [get]
func (h *handler) GetEvents(c *echo.Context) error {
	events, err := h.service.GetEvents()
	if err != nil {
		return eventErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, events)
}

// GetEventsByID godoc
// @Summary      Get event by ID
// @Description  Retrieves the details of a specific event by its ID.
// @Tags         Events
// @Produce      json
// @Param        id   path      int  true  "Event ID"
// @Success      200  {object}  dto.Response
// @Failure      400  {object}  httpresponse.Error
// @Failure      404  {object}  httpresponse.Error
// @Failure      500  {object}  httpresponse.Error
// @Router       /api/v1/events/{id} [get]
func (h *handler) GetEventsByID(c *echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid event id",
			Details: err.Error(),
		})
	}

	response, err := h.service.GetEventByID(uint(id)) // err => re-assign

	if err != nil {
		return eventErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, response)
}

// UpdateEvent godoc
// @Summary      Update an event
// @Description  Updates the details of a specific event by its ID.
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        id       path      int                true  "Event ID"
// @Param        request  body      dto.UpdateRequest  true  "Event Update Details"
// @Success      200      {object}  dto.Response
// @Failure      400      {object}  httpresponse.Error
// @Failure      404      {object}  httpresponse.Error
// @Failure      500      {object}  httpresponse.Error
// @Router       /api/v1/events/{id} [patch]
func (h *handler) UpdateEvent(c *echo.Context) error {
	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid event id",
			Details: err.Error(),
		})
	}

	var req dto.UpdateRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Details: err.Error(),
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Validation failed",
			Details: err.Error(),
		})
	}

	response, err := h.service.UpdateEvent(uint(eventId), req)
	if err != nil {
		return eventErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, response)
}
