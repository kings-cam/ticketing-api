package tickets

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	// Mongo DB
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// dayinexcludedays returns true if the day is found in excluded days
func dayinexcludedays(d time.Time, excludedays []time.Weekday) bool {
	// Get day 
	day := d.Weekday()
	
	excludeday := false
	// Check if the given day is in the exclude list
	for i := range excludedays {
		if excludedays[i] == day {
			excludeday = true
		}
	}
	return excludeday
}

// dayinexcludedates returns true if the day is found in excluded days
func dayinexcludedates(date time.Time, excludedates []string) bool {
	// Get day 
	excludedate := false
	// Check if the given day is in the exclude list
	for i := range excludedates {
		if excludedates[i] == date.Format("2006-01-02") {
			excludedate = true
		}
	}
	return excludedate
}

// BookingDates return allowable booking days
func BookingDates(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// Copy and launch a Mongo session
		session := s.Copy()
		defer session.Close()

		// Open config collections
		dbc := session.DB("tickets").C("config")

		var config BookingConfig
		// Find the configuration file
		err := dbc.Find(bson.M{"id": 0}).One(&config)
		if err != nil {
			ErrorWithJSON(w, "Database error, failed to find config!", http.StatusInternalServerError)
			log.Println("Failed to find config: ", err)
			return
		}

		// Booking dates
		var bookingdates []string

		// Start date as tomorrow
		startdate := time.Now().Local().AddDate(0, 0, 1)
		
		// End date as 90 days (3 months) from tomorrow
		enddate := startdate.AddDate(0, 0, 90)

		// Exclude dates
		excludedates := config.ExcludeDays
		log.Println(excludedates)

		// Exclude days
		excludedays := []time.Weekday{0,6}
		
		// Iterate over dates to print all allowed dates
		for d := startdate; d != enddate; d = d.AddDate(0, 0, 1) {
			// Exclude weekends (0 - Sunday, 6 - Saturday)
			if (!dayinexcludedays(d, excludedays) &&
				!dayinexcludedates(d, excludedates)) {
				bookingdates = append(bookingdates, d.Format("2006-01-02"))
			}
		}
		
		// Marshall booking dates
		respBody, err := json.MarshalIndent(bookingdates, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		ResponseWithJSON(w, respBody, http.StatusOK)
	}
}
