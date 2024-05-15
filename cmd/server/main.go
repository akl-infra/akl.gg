package main

import (
	"github.com/akl-infra/api/internal/auth"
	"github.com/akl-infra/api/internal/handlers"
	"github.com/akl-infra/api/internal/setup"
	"github.com/akl-infra/api/internal/storage"
	"github.com/charmbracelet/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	log.Info("Booting up")

	srv := echo.New()
	setup.Middleware(srv)

	api := srv.Group("/api")

	// Public
	api.GET("/", handlers.Banner)
	api.GET("/layouts", handlers.Layouts)
	api.GET("/layout/:name", handlers.Layout)

	// Protected
	api.PUT("/layout", handlers.AddLayout, middleware.KeyAuth(auth.TokenValidator))

	storage.Init("layouts")
	Server(srv)
}
