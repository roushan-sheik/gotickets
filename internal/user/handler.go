package user

import (
	"gotickets/internal/auth"
	"gotickets/internal/httpresponse"
	"gotickets/internal/user/dto"
	"net/http"

	"github.com/labstack/echo/v5"
)

type handler struct {
	service *service
}

func NewHandler(service *service) *handler {
	return &handler{
		service: service,
	}
}

func (h handler) CreateUser(c *echo.Context) error {
	var req dto.CreateRquest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid request payload", err.Error()))
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Validation failed", err.Error()))
	}

	response, err := h.service.CreateUser(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Error creating user", err.Error()))
	}

	return c.JSON(http.StatusOK, response)

}

func (h handler) LoginUser(c *echo.Context) error {
	var req dto.LoginRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Invalid request payload", err.Error()))
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Validation failed", err.Error()))
	}

	response, err := h.service.LoginUser(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError(http.StatusBadRequest, "Error logging in", err.Error()))
	}

	return c.JSON(http.StatusOK, response)
}

func (h handler) GetMe(c *echo.Context) error {
	userClaims, ok := c.Get("user").(*auth.JwtCustomClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.NewError(http.StatusUnauthorized, "Unauthorized", "User not found in context"))
	}
	return c.JSON(http.StatusOK, userClaims)
}
