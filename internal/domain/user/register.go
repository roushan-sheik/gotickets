package user

import (
	"gotickets/internal/auth"
	"gotickets/internal/config"
	"gotickets/internal/middlewares"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	userRepository := NewRepository(db)
	jwtService := auth.NewJWTService(cfg.JwtSecret)
	userService := NewService(userRepository, jwtService)
	userHandler := NewHandler(userService)

	authMiddleware := middlewares.AuthMiddleware(jwtService)

	authGroup := e.Group("/api/v1/auth")
	authGroup.POST("/register", userHandler.CreateUser)
	authGroup.POST("/login", userHandler.LoginUser)
	authGroup.POST("/refresh", userHandler.RefreshToken)

	userGroup := e.Group("/api/v1/users", authMiddleware)
	userGroup.GET("/me", userHandler.GetMe)
}
