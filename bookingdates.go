package tickets

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	// Mongo DB
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	// Gorilla Mux
	"github.com/gorilla/mux"
)

type DateSession struct {
	Date string `json:"date"`
	// Typical number of tickets in the morning session
	NMorningTickets int `json:"nmorningtickets"`
	// Typical number of tickets in the evening session
	NAfternoonTickets int `json:"nafternoontickets"`
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

// createBookingDates generates allowed booking dates
func createBookingDates(s *mgo.Session, test bool) error {
	// Copy and launch a Mongo session
	session := s.Copy()
	defer session.Close()
	
	// Test Config DB or Production config
	var configtable string
	if (test) {
		configtable = "testconfig"
	} else {
		configtable = "config"
	}
	
	// Open config collections
	dbc := session.DB("tickets").C(configtable)
	
	var config BookingConfig
	// Find the configuration file
	err := dbc.Find(bson.M{"id": 0}).One(&config)
	if err != nil {
		log.Println("Failed to find config: ", err)
		return err
	}

	// Booking dates
	var bookingdates []string
	
	// Start date as tomorrow
	startdate := time.Now().Local().AddDate(0, 0, 1)
	
	// End date as 90 days (3 months) from tomorrow
	enddate := startdate.AddDate(0, 0, config.NDays)
	
	// Exclude dates
	excludedates := config.ExcludeDates
	
	// Exclude days
	excludedays := config.ExcludeDays
	
	// Iterate over dates to print all allowed dates
	for d := startdate; d != enddate; d = d.AddDate(0, 0, 1) {
		// Exclude weekends (0 - Sunday, 6 - Saturday)
		if (!excludedays[d.Weekday()] &&
			!dayinexcludedates(d, excludedates)) {
			bookingdates = append(bookingdates, d.Format("2006-01-02"))
		}
	}
	config.BookingDates = bookingdates
	
	// Insert bookingdates to database
	err = dbc.Update(bson.M{"id": 0}, &config)

        if err != nil {
		switch err {
		default:
			log.Println("Failed to update bookingdates: ", err)
			return err
		case mgo.ErrNotFound:
			log.Println("Error config not found: ", err)
			return err
		}
	}

	// Create session for dates
	var newdate DateSession
	newdate.NMorningTickets = config.NMorningTickets
	newdate.NAfternoonTickets = config.NAfternoonTickets

	for _, bd := range bookingdates {
		newdate.Date = bd
		// Do not update existing sessions
		err = createBookingSessions(session, &newdate, false)
	}
	return err
}


// BookingDates return allowable booking days
func BookingDates(s *mgo.Session, test bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// Copy and launch a Mongo session
		session := s.Copy()
		defer session.Close()

		// Test Config DB or Production config
		var configtable string
		if (test) {
			configtable = "testconfig"
		} else {
			configtable = "config"
		}

		// Open config collections
		dbc := session.DB("tickets").C(configtable)

		var config BookingConfig
		
		// Find the configuration file
		err := dbc.Find(bson.M{"id": 0}).One(&config)
		if err != nil {
			ErrorWithJSON(w, "Database error, failed to find config!", http.StatusInternalServerError)
			log.Println("Failed to find config: ", err)
			return
		}

		// Start date as tomorrow
		startdate := time.Now().Local().AddDate(0, 0, 1)

		// Booking dates
		var bookingdates []string

		datelayout := "2006-01-02"
		// Check if booking date is more than the current date
		for _, bookingdate := range config.BookingDates {
			t, err := time.Parse(datelayout, bookingdate)
			if err == nil {
				if (t.Sub(startdate) >= 0) {
					bookingdates = append(bookingdates, bookingdate)
				}
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


// createBookingSessions generates allowed booking dates
func createBookingSessions(s *mgo.Session, date *DateSession, update bool) error {
	// Copy and launch a Mongo session
	session := s.Copy()
	defer session.Close()
	
	// Open config collections
	dbc := session.DB("tickets").C("dates")

	var err error
	
	if (update) {
		// Try to update if date is found
		err = dbc.Update(bson.M{"date": date.Date}, &date)
	} else {
		// Create a new insert
		var newdate DateSession
		// Find id date exists
		err = dbc.Find(bson.M{"date": date.Date}).One(&newdate)
	}		
	if err != nil {
		switch err {
		default:
			log.Println("Failed to update date: ", err)
			return err
			// Configuration is not present, do an insert
		case mgo.ErrNotFound:
			// Date not present, creating a new date/session
			err = dbc.Insert(&date)
			if err != nil {
				log.Println("Failed to insert session: ", err)
				return err
			}
		}
	}
	return err
}


// BookingSessions return ntickets for each session in a day
func BookingSessions(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// Copy and launch a Mongo session
		session := s.Copy()
		defer session.Close()

		// Open dates collections
		dbc := session.DB("tickets").C("dates")

		// date
		params := mux.Vars(r)
		sessiondate := params["date"]


		var sess DateSession
		// Find the configuration file
		err := dbc.Find(bson.M{"date": sessiondate}).One(&sess)
		if err != nil {
			ErrorWithJSON(w, "Database error, failed to find date!", http.StatusInternalServerError)
			log.Println("Failed to find date: ", err)
			return
		}

		// Marshall booking dates
		respBody, err := json.MarshalIndent(sess, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		ResponseWithJSON(w, respBody, http.StatusOK)
	}
}

