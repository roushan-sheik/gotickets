package server

import (
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
