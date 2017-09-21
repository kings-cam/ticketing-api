package main

import (
	// Booking dates
	"ticketing-api/controller"
	// HTTP requests
	"net/http"
	// Echo framework
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Handle root
func welcome(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to King's Chapel Ticketing API!")
}

// Server API
func main() {
	e := echo.New()
	
	// Middleware
	// Server log
	e.Use(middleware.Logger())
	// Recovers from panics anywhere in the chain
	e.Use(middleware.Recover())

	// CORS for cross-domain access controls
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	// Route welcome
	e.GET("/api/v1/", welcome)

	// Route dates
	e.GET("/api/v1/dates/", bookingdates.BookingDates)
	
	// Server
	e.Logger.Fatal(e.Start(":4000"))
}
