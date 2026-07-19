package main

import (
	"fmt"
	"gotickets/internal/config"
	"gotickets/internal/user"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

func main() {
	e := echo.New()

	config := config.LoadEnv()

	dsn := config.Dsn
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})

	db.AutoMigrate(&user.User{})

	e.Validator = &CustomValidator{validator: validator.New()}

	if err != nil {
		panic("failed to connect database")
	}
	e.Logger.Info("Database connected successfully")

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	user.RegisterRoutes(e, db)

	port := fmt.Sprintf(":%s", config.Port)

	fmt.Printf("🚀🚀Server running on http://localhost%s\n", port)

	err = e.Start(port)

	if err != nil {
		e.Logger.Error("failed to start server:", "error", err)
	}
}
