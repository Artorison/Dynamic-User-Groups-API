package app

import (
	"log"

	"github.com/labstack/echo/v4"
)

type App struct {
	Router      *echo.Echo
	DIContainer *DIContainer
}

func NewApp(router *echo.Echo, container *DIContainer) *App {
	return &App{
		Router:      router,
		DIContainer: container,
	}
}

func (a *App) Start(address string) {
	log.Printf("Starting server at http://localhost%s/", address)
	if err := a.Router.Start(address); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
