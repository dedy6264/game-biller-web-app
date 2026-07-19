package main

import (
	"gamebiller/configs"
	"gamebiller/connections"
	"gamebiller/routes"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Initialize DB Connection
	connections.Connect()
	defer connections.DBconn().Close()

	e := echo.New()

	// Middleware
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	// Register Routes
	routes.AppRoutes(e)

	// Start Server
	e.Logger.Fatal(e.Start(":" + configs.APP_PORT))
}
