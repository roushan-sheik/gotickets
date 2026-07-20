package server

import (
	"fmt"
	"gotickets/internal/config"
	"gotickets/internal/domain/booking"
	"gotickets/internal/domain/event"
	"gotickets/internal/domain/user"
	"net/http"
	"time"

	_ "gotickets/docs"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"gorm.io/gorm"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.ErrBadRequest.Wrap(err)
	}
	return nil
}

func Start(db *gorm.DB, cfg *config.Config) {
	db.AutoMigrate(&user.User{}, &event.Event{}, &booking.Booking{})

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.RequestLogger())

	e.GET("/", WelcomeHandler)
	e.GET("/health", HealthCheckHandler(db))
	
	// Swagger UI Route
	e.GET("/swagger/*", echo.WrapHandler(httpSwagger.Handler()))

	//routes
	user.RegisterRoutes(e, db, cfg)
	event.RegisterRoutes(e, db)
	booking.RegisterRoutes(e, db, cfg)

	port := fmt.Sprintf(":%s", cfg.Port)
	fmt.Printf("\033[1;32m🚀 Server is running on http://localhost:%s\033[0m\n", cfg.Port)
	if err := e.Start(port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}

// WelcomeHandler godoc
// @Summary      Welcome to GoTickets API
// @Description  Root endpoint to verify the API is running
// @Tags         System
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       / [get]
func WelcomeHandler(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":     "Welcome to GoTickets API",
		"version":     "1.0.0",
		"environment": "development",
		"status":      "active",
	})
}

// HealthCheckHandler godoc
// @Summary      Health Check
// @Description  Check the health status of the API and Database
// @Tags         System
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /health [get]
func HealthCheckHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c *echo.Context) error {
		dbStatus := "up"
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			dbStatus = "down"
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":    "up",
			"database":  dbStatus,
			"timestamp": time.Now().Format(time.RFC3339),
		})
	}
}
