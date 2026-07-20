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

	e.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Welcome to GoTickets API",
			"version": "1.0.0",
			"environment": "development",
			"status":  "active",
		})
	})

	e.GET("/health", func(c *echo.Context) error {
		dbStatus := "up"
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			dbStatus = "down"
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "up",
			"database": dbStatus,
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	//routes
	user.RegisterRoutes(e, db, cfg)
	event.RegisterRoutes(e, db)
	booking.RegisterRoutes(e, db, cfg)

	port := fmt.Sprintf(":%s", cfg.Port)
	if err := e.Start(port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
