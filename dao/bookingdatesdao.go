package bookingdatesdao

import (
	// "ticketing-api/model/bookingdatesmodel"
	"time"
	// "error"
)

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

	if (len(bookingdates) == 0) {
		
	}
	return bookingdates
}
