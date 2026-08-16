package main

import (
	"encoding/json"
	"fmt"
	"gamebiller/configs"
	"gamebiller/connections"
	"gamebiller/helpers"
	"gamebiller/repositories"
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
	var db = connections.DBconn()
	productSegment, err := repositories.GetProductSegmentJoinProvider(db, 6, "ML_5", "")
	if err != nil {
		helpers.ProcessLogger(nil, "svc", err.Error(), "Failed to get product segment")
	}
	aa, _ := json.Marshal(productSegment)
	fmt.Println(string(aa))
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
