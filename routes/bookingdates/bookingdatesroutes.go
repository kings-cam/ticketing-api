package bookingdatesroutes

import (
	"ticketing-api/controller/bookingdates"

	// Echo webframework
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func Init(e *echo.Echo) {
	e.Pre(middleware.RemoveTrailingSlash())

	// Booking dates
	e.GET("/api/v1/dates", bookingdatescontroller.BookingDates)
}

