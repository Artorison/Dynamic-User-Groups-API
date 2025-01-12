package main

import (
	"API/internal/app"
	"API/internal/config"
	"log"
	"log/slog"

	_ "API/docs"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Dynamic User Groups API
// @version 1.0
// @description API documentation.

// @contact.name API Support
// @contact.email artorison@gmail.com

// @host localhost:8080
// @BasePath /
func main() {

	cfg, err := config.LoadDBConfig("config/config.yml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	container := app.InitDI(*cfg)

	router := echo.New()
	application := app.NewApp(router, container)

	app.StartTTLConsumer(container.KafkaConsumer, container.UserSegmentService)

	application.Router.GET("/swagger/*", echoSwagger.WrapHandler)
	slog.Info("Swagger page: http://localhost:8080/swagger/index.html")
	slog.Info("pgadmin: http://localhost:5050")

	app.RegisterMiddleware(application.Router)
	app.RegisterRoutes(application.Router, application.DIContainer)

	application.Start(cfg.Server.Address)
	slog.Error("Server stopped")
}
