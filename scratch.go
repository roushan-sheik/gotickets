package main

import (
	"github.com/labstack/echo/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func main() {
	e := echo.New()
	e.GET("/swagger/*", echo.WrapHandler(httpSwagger.Handler()))
}
