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

type BookingConfig struct {
	ID int `json:"id"`
	// Typical number of tickets per session
	NTickets int `json:"ntickets"`
	// Exclude dates
	ExcludeDates []string `json:"excludedates, omitempty"`
	// Exclude days (0 - Sunday, 6 - Saturday)
	ExcludeDays []time.Weekday `json:"excludedays, omitempty"`
	// Booking dates
	BookingDates []string `json:"bookingdates, omitempty"`
}


// ConfigBookingDates assign excludedays and dates
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

		err = createBookingDates(session)
		if err != nil {
			switch err {
			default:
				ErrorWithJSON(w, "Database error, failed to update project", http.StatusInternalServerError)
				return
			case mgo.ErrNotFound:
				ErrorWithJSON(w, "Project not found", http.StatusNotFound)
				return
			}
		}

		// Write response
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", r.URL.Path + "/0")
		w.WriteHeader(http.StatusCreated)
	}
}
