package bookingdatescontroller

import (
	"ticketing-api/dao/bookingdates"
	// "encoding/json"
	"net/http"

	"github.com/labstack/echo"
)


// Return dates of bookings
func BookingDates(c echo.Context) error {
	// bd, _ to get errors
	bd := bookingdatesdao.BookingDates()
	return c.JSON(http.StatusCreated, bd)
}
