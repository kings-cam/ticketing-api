package bookingdates

import (
	// "encoding/json"
	"time"
	"github.com/labstack/echo"
	"net/http"
)

var bookingdates []string

// Return dates of bookings
func BookingDates(c echo.Context) error {
	// Start date as tomorrow
	startdate := time.Now().Local().AddDate(0, 0, 1)

	// End date as 3 months from tomorrow
	enddate := startdate.AddDate(0, 3, 0)
	
	// Iterate over dates to print all allowed dates
	for d := startdate; d != enddate; d = d.AddDate(0, 0, 1) {
		// Exclude weekends (0 - Sunday, 6 - Saturday)
		if d.Weekday() != 0 {
			bookingdates = append(bookingdates, d.Format("2006-01-02"))
		}
	}
	return c.JSON(http.StatusCreated, bookingdates)
}
