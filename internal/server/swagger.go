package server

import (
	"fmt"
	"net/http"

	_ "gotickets/docs"

	"github.com/labstack/echo/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// WelcomeResponse defines the shape of the root endpoint response.
type WelcomeResponse struct {
	Message     string `json:"message"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
	Status      string `json:"status"`
	DocsURL     string `json:"docs_url"`
}

// WelcomeHandler godoc
// @Summary      Welcome to GoTickets API
// @Description  Root endpoint to verify the API is running.
// @Tags         System
// @Produce      json
// @Success      200  {object}  WelcomeResponse
// @Router       / [get]
func WelcomeHandler(c *echo.Context) error {
	scheme := "http"
	if c.Request().TLS != nil {
		scheme = "https"
	}
	return c.JSON(http.StatusOK, WelcomeResponse{
		Message:     "Welcome to GoTickets API",
		Version:     "1.0.0",
		Environment: "development",
		Status:      "active",
		DocsURL:     fmt.Sprintf("%s://%s/swagger/index.html", scheme, c.Request().Host),
	})
}

// HealthResponse defines the shape of the health check endpoint response.
type HealthResponse struct {
	Status    string `json:"status"`
	Database  string `json:"database"`
	Timestamp string `json:"timestamp"`
}

// RegisterSwagger mounts the Swagger UI documentation route onto the Echo router.
// Access the UI at: /swagger/index.html
func RegisterSwagger(e *echo.Echo) {
	e.GET("/swagger/*", echo.WrapHandler(httpSwagger.Handler()))
}
