package user

import (
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	userRepository := NewRepository(db)
	userService := NewService(userRepository)
	userHandler := NewHandler(userService)

	user := e.Group("/api/v1")
	user.POST("/users", userHandler.CreateUser)
}
