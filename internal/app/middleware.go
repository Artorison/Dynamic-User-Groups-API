package app

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterMiddleware(e *echo.Echo) {
	e.Use(middleware.Recover())
}
