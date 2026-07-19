package user

import (
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	userRepository := NewRepository(db)
	userService := NewService(userRepository)
	userHandler := NewHandler(userService)

	auth := e.Group("/api/v1/auth")
	auth.POST("/register", userHandler.CreateUser)
	auth.POST("/login", userHandler.LoginUser)
}
