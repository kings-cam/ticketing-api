package tickets

import (
	// HTTP requests
	"net/http"
	"time"
	// "gopkg.in/mgo.v2/bson"
	// Echo webframework
	"github.com/labstack/echo"
)

/*
type BookingConfig struct {
	// Start date for booking
	StartDate string `json:"startdate, omitempty" bson:"startdate, omitempty"`
	// Number of days from start date
	EndDays int `json:"enddays, omitempty" bson:"enddays, omitempty"`
	// Exclude days
	ExcludeDays []string `json:"excludedays, omitempty" bson:"excludedays, omitempty"`
}

type bookingconfig BookingConfig
*/

func BookingDates() []string {
	var bookingdates []string

	// Start date as tomorrow
	startdate := time.Now().Local().AddDate(0, 0, 1)

	// End date as 90 days (3 months) from tomorrow
	enddate := startdate.AddDate(0, 0, 90)

	// Iterate over dates to print all allowed dates
	for d := startdate; d != enddate; d = d.AddDate(0, 0, 1) {
		// Exclude weekends (0 - Sunday, 6 - Saturday)
		if d.Weekday() != 0 {
			bookingdates = append(bookingdates, d.Format("2006-01-02"))
		}
	}

	return bookingdates
}

// Return dates of bookings
func GetBookingDates(c echo.Context) error {
	// bd, _ to get errors
	bd := BookingDates()
	return c.JSON(http.StatusCreated, bd)
}
