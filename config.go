package tickets

import (
	"encoding/json"
	"log"
	"net/http"

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
	ExcludeDays []string `json:"excludedays, omitempty"`
	// Booking days
	// BookingDays []string `json:"bookingdays, omitempty"`
}

func ConfigBookingDates(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		var config BookingConfig

		// Decode POST JSON file
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&config)
		if err != nil {
			ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			return
		}

		// Open config collection
		c := session.DB("tickets").C("config")

		// Try to update if configuration is found
		err = c.Update(bson.M{"id": 0}, &config)
		if err != nil {
			switch err {
			default:
				ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
				log.Println("Failed to update config: ", err)
				return
			// Configuration is not present, do an insert
			case mgo.ErrNotFound:
				log.Println("Config not present, creating a new config")
				err = c.Insert(&config)

				if err != nil {
					ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
					log.Println("Failed to insert config: ", err)
					return
				}
			}
		}

		// Write response
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", r.URL.Path + "/0")
		w.WriteHeader(http.StatusCreated)
	}
}
