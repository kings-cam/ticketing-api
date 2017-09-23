package routes

import (
	"ticketing-api/controller"
	// HTTP requests
	"net/http"
	// Echo webframework
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Welcome message
func welcome(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to King's Chapel Ticketing API version 1!")
}

func Init(e *echo.Echo) {
	e.Pre(middleware.RemoveTrailingSlash())

	// Route welcome
	e.GET("/api/v1", welcome)

	// Booking dates
	e.GET("/api/v1/dates", bookingdates.BookingDates)
}
