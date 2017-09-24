package tickets

import (
	"encoding/json"
	"log"
	"net/http"
	// "time"

	// Mongo DB
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type BookingConfig struct {
	ID int `json:"id"`
	// Start date for booking
	// StartDate string `json:"startdate, omitempty"`
	// Number of days from start date
	// EndDays int `json:"enddays, omitempty"`
	// Exclude days
	// ExcludeDays []string `json:"excludedays, omitempty"`
	// Booking days
	BookingDays []string `json:"bookingdays, omitempty"`
}

func ConfigBookingDates(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {  
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		var config BookingConfig
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&config)
		if err != nil {
			ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			return
		}

		c := session.DB("tickets").C("config")

		err = c.Insert(config)
		if err != nil {
			if mgo.IsDup(err) {
				ErrorWithJSON(w, "Config with this ID already exists", http.StatusBadRequest)
				return
			}

			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed insert config: ", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", r.URL.Path+"/0")
		w.WriteHeader(http.StatusCreated)
	}
}

func BookingDates(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// Copy and launch a Mongo session
		session := s.Copy()
		defer session.Close()

		// Open config collections
		dbc := session.DB("tickets").C("config")

		var config BookingConfig
		// Find the configuration file
		err := dbc.Find(bson.M{}).One(&config)
		if err != nil {
			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed to find config: ", err)
			return
		}

		if len(config.BookingDays) == 0 {
			ErrorWithJSON(w, "Config not found", http.StatusNotFound)
			log.Println("Configuration with id not found: ", err)
			return
		}

		respBody, err := json.MarshalIndent(config.BookingDays, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		
		/*
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
		*/
		ResponseWithJSON(w, respBody, http.StatusOK)
	}
}
