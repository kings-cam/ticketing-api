package bookingdates

import (
	// "encoding/json"
	"fmt"
	"time"
	"github.com/labstack/echo"
	"net/http"
)

var bookingdates []string

// Return dates of bookings
func BookingDates(c echo.Context) error {
	// Errors
	// var err error
	// Start date as tomorrow
	startdate := time.Now().Local().AddDate(0, 0, 1)
	fmt.Println("Start date: ", startdate.Format("2006-01-02"))
	// End date as 3 months from tomorrow
	enddate := startdate.AddDate(0, 2, 0)
	fmt.Println("End date: ", enddate.Format("2006-01-02"))

	// Iterate over dates to print all allowed dates
	for d := startdate; d != enddate; d = d.AddDate(0, 0, 1) {
		// Exclude weekends (0 - Sunday, 6 - Saturday)
		if d.Weekday() != 0 {
			bookingdates = append(bookingdates, d.Format("2006-01-02"))
		}
	}
	/*
	dates, err := json.Marshal(bookingdates)
	if err != nil {
		fmt.Println("JSON marshaling failed: %s", err)
	}
	fmt.Printf("%s\n", dates)
        */
	return c.JSON(http.StatusCreated, bookingdates)
}
