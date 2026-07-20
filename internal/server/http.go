package server

import (
	"fmt"
	"gotickets/internal/config"
	"gotickets/internal/domain/booking"
	"gotickets/internal/domain/event"
	"gotickets/internal/domain/user"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"gorm.io/gorm"
)

type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.ErrBadRequest.Wrap(err)
	}
	return nil
}

// Start initializes and runs the HTTP server.
func Start(db *gorm.DB, cfg *config.Config) {
	db.AutoMigrate(&user.User{}, &event.Event{}, &booking.Booking{})

	e := echo.New()
	e.Validator = &customValidator{validator: validator.New()}
	e.Use(middleware.RequestLogger())

	// System routes
	e.GET("/", WelcomeHandler)
	e.GET("/health", HealthCheckHandler(db))

	// Documentation
	RegisterSwagger(e)

	// Domain routes
	user.RegisterRoutes(e, db, cfg)
	event.RegisterRoutes(e, db)
	booking.RegisterRoutes(e, db, cfg)

	port := fmt.Sprintf(":%s", cfg.Port)
	fmt.Printf("\033[1;32m🚀 Server is running on http://localhost:%s\033[0m\n", cfg.Port)
	if err := e.Start(port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}




// HealthCheckHandler godoc
// @Summary      Health Check
// @Description  Check the health status of the API and the database connection.
// @Tags         System
// @Produce      json
// @Success      200  {object}  HealthResponse
// @Router       /health [get]
func HealthCheckHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c *echo.Context) error {
		dbStatus := "up"
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			dbStatus = "down"
		}

		return c.JSON(http.StatusOK, HealthResponse{
			Status:    "up",
			Database:  dbStatus,
			Timestamp: time.Now().Format(time.RFC3339),
		})
	}
}
