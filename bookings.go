package tickets

import (
	"encoding/json"
	"log"
	"net/http"

	// Mongo DB
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	// Gorilla Mux
	"github.com/gorilla/mux"
)

// Type Booking stores the datastructure for a booking
type Booking struct {
	// Unique ID
	UUID string `json:"uuid, omitempty"`
	// Name of booking
	Name string `json:"name, omitempty"`
	// date
	Date string `json:"date, omitempty"`
	// Session
	Session string `json:"session, omitempty"`
	// Total
	Total float64 `json:"total, omitempty"`
	// Ntickets
	Ntickets int `json:"ntickets, omitempty"`
	// Nadults
	NAdults int `json:"nadults, omitempty"`
	// Nchild
	Nchild int `json:"nchild, omitempty"`
	// Nconcession
	Nconcession int `json:"nconcession, omitempty"`
	// Nguides
	Nguides int `json:"nguides, omitempty"`
	// Guidebooks
	Guidebooks []string `json:"guidebooks, omitempty"`
}

// GetBookings returns all bookings
func GetBookings(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Copy and launch a Mongo session
		session := s.Copy()
		defer session.Close()

		// Open bookings collections
		dbc := session.DB("tickets").C("bookings")

		var bookings []Booking

		// Find all Bookings
		err := dbc.Find(bson.M{}).All(&bookings)
		if err != nil {
			ErrorWithJSON(w, "Database error, failed to find bookings!", http.StatusInternalServerError)
			log.Println("Failed to find bookings: ", err)
			return
		}

		// Marshall booking
		respBody, err := json.MarshalIndent(bookings, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		ResponseWithJSON(w, respBody, http.StatusOK)
	}
}

// GetBooking returns booking with a matching uuid
// BUG(r) This function returns a booking matching uuid
func GetBooking(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Copy and launch a Mongo session
		session := s.Copy()
		defer session.Close()

		// Open bookings collections
		dbc := session.DB("tickets").C("bookings")

		vars := mux.Vars(r)
		bookingid := vars["uuid"]

		var booking Booking

		// Find the Booking matching the id
		err := dbc.Find(bson.M{"uuid": bookingid}).One(&booking)
		if err != nil {
			ErrorWithJSON(w, "Database error, failed to find booking with uuid!", http.StatusInternalServerError)
			log.Println("Failed to find booking with the specified uuid: ", err)
			return
		}

		// Marshall booking
		respBody, err := json.MarshalIndent(booking, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		ResponseWithJSON(w, respBody, http.StatusOK)
	}
}

// GetBookingsDate returns booking matching a date
// BUG(r) This function returns a booking matching uuid
func GetBookingsDate(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Copy and launch a Mongo session
		session := s.Copy()
		defer session.Close()

		// Open bookings collections
		dbc := session.DB("tickets").C("bookings")

		vars := mux.Vars(r)
		date := vars["date"]

		var bookings []Booking

		// Find all bookings matching the date
		err := dbc.Find(bson.M{"date": date}).All(&bookings)
		if err != nil {
			ErrorWithJSON(w, "Database error, failed to find booking with date!", http.StatusInternalServerError)
			log.Println("Failed to find booking with the specified date: ", err)
			return
		}

		// Marshall booking
		respBody, err := json.MarshalIndent(bookings, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		ResponseWithJSON(w, respBody, http.StatusOK)
	}
}

// CreateBooking creates a new booking object and writes to MongoDB booking collection.
// BUG(r) This function doesn't check if the incoming body matches the booking structure
// BUG(r) This function doesn't check if UUID is a valid UUID
func CreateBooking(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Open bookings collection
		session := s.Copy()
		defer session.Close()
		dbc := session.DB("tickets").C("bookings")

		// UUID
		params := mux.Vars(r)
		uuid := params["uuid"]

		// Create a booking to store incoming JSON booking request
		var booking Booking

		// Decode POST JSON file
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&booking)
		if err != nil {
			ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			return
		}

		// Check if the URL uuid and the booking JSON UUID match
		if uuid != booking.UUID || booking.UUID == "" {
			ErrorWithJSON(w, "Database error, uuid of request doesn't match booking uuid", http.StatusInternalServerError)
			log.Println("Error booking uuids don't match: ", uuid, booking.UUID)
			return
		}

		// Insert booking to database
		err = dbc.Insert(&booking)
		if err != nil {
			// Check if a booking with a same UUID exists
			if mgo.IsDup(err) {
				ErrorWithJSON(w, "Booking with uuid already exists", http.StatusBadRequest)
				return
			}
			ErrorWithJSON(w, "Database error, failed to insert booking", http.StatusInternalServerError)
			log.Println("Failed to insert booking: ", err)
			return
		}

		// Write response
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", r.URL.Path)
		w.WriteHeader(http.StatusCreated)
	}
}

// UpdateBooking updates an existing booking and updates MongoDB collection.
// BUG(r) This function doesn't check if the incoming body matches the booking structure
// BUG(r) This function doesn't check if UUID is a valid UUID
func UpdateBooking(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Open bookings collection
		session := s.Copy()
		defer session.Close()
		dbc := session.DB("tickets").C("bookings")

		// UUID
		params := mux.Vars(r)
		uuid := params["uuid"]

		// Create a booking to store incoming JSON booking request
		var booking Booking

		// Decode POST JSON file
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&booking)
		if err != nil {
			ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			return
		}

		// Check if the URL uuid and the booking JSON UUID match
		if uuid != booking.UUID || booking.UUID == "" {
			ErrorWithJSON(w, "Database error, uuid of request doesn't match booking uuid", http.StatusInternalServerError)
			log.Println("Error booking uuids don't match: ", uuid, booking.UUID)
			return
		}

		// Insert booking to database
		err = dbc.Update(bson.M{"uuid": uuid}, &booking)
		if err != nil {
			switch err {
			default:
				ErrorWithJSON(w, "Database error, failed to update booking", http.StatusInternalServerError)
				log.Println("Failed to update booking: ", err)
				return
			case mgo.ErrNotFound:
				ErrorWithJSON(w, "Booking not found", http.StatusNotFound)
				log.Println("Error booking with uuid not found: ", booking.UUID)

				return
			}
		}

		// Write response
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", r.URL.Path)
		w.WriteHeader(http.StatusCreated)
	}
}

// DeleteBooking deletes an existing booking
// BUG(r) This function doesn't check if UUID is a valid UUID
func DeleteBooking(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Open bookings collection
		session := s.Copy()
		defer session.Close()
		dbc := session.DB("tickets").C("bookings")

		// UUID
		params := mux.Vars(r)
		uuid := params["uuid"]

		err := dbc.Remove(bson.M{"uuid": uuid})
		if err != nil {
			switch err {
			default:
				ErrorWithJSON(w, "Database error, failed to delete booking", http.StatusInternalServerError)
				log.Println("Failed to delete booking: ", err)
				return
			case mgo.ErrNotFound:
				ErrorWithJSON(w, "Booking not found", http.StatusNotFound)
				log.Println("Booking not found: ", err)
				return
			}
		}

		// Write response
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", r.URL.Path)
		w.WriteHeader(http.StatusCreated)
	}
}
