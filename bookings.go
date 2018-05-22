package tickets

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	// Mongo DB
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	// Gorilla Mux
	"github.com/gorilla/mux"
	// QR Code generator
	"github.com/skip2/go-qrcode"
)

// Type Booking stores the datastructure for a booking
type Booking struct {
	// Unique ID
	UUID string `json:"uuid, omitempty"`
	// Name of booking
	Name string `json:"name, omitempty"`
	// Email address
	Email string `json:"email, omitempty"`
	// Gift aid (true / false)
	Giftaid string `json:"giftaid, omitempty"`
	// Subscribe (true / false)
	Subscribe string `json:"subscribe, omitempty"`
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
	// Address
	Address string `json:"address, omitempty"`
	// City
	City string `json:"city, omitempty"`
	// Country
	Country string `json:"country, omitempty"`
	// City
	Postcode string `json:"postcode, omitempty"`
	// CC Number
	CCNumber string `json:"ccnumber, omitempty"`
	// CVV
	CVV string `json:"cvv, omitempty"`
	// Nadults
	Month int `json:"month, omitempty"`
	// Nchild
	Year int `json:"year, omitempty"`
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

// GetBookingsRange returns the summaray of booking matching a start - end date
// BUG(r) This function returns a booking matching uuid
func GetBookingsRangeSummary(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Copy and launch a Mongo session
		session := s.Copy()
		defer session.Close()

		// Open bookings collections
		dbc := session.DB("tickets").C("bookings")

		vars := mux.Vars(r)
		startdate := strings.Split(vars["date"], ":")[0]
		enddate := strings.Split(vars["date"], ":")[1]

		var bookings []Booking

		// Find all bookings matching the date
		err := dbc.Find(bson.M{"date": bson.M{"$gte": startdate, "$lte": enddate}}).All(&bookings)
		if err != nil {
			ErrorWithJSON(w, "Database error, failed to find booking with date!", http.StatusInternalServerError)
			log.Println("Failed to find booking with the specified date: ", err)
			return
		}

		var summary Booking
		for _, booking := range bookings {
			summary.Total += booking.Total
			summary.Ntickets += booking.Ntickets
			summary.NAdults += booking.NAdults
			summary.Nchild += booking.Nchild
			summary.Nconcession += booking.Nconcession
			summary.Nguides += booking.Nguides
		}
		summary.Name = "Total # of booking: " + strconv.Itoa(len(bookings))

		// Marshall booking
		respBody, err := json.MarshalIndent(summary, "", "  ")
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

		// Payment
		var payment Payment
		payment.Name = booking.Name
		payment.OrderDescription = booking.UUID
		payment.CardNumber = booking.CCNumber
		payment.Cvc = booking.CVV
		payment.Month = strconv.Itoa(booking.Month)
		payment.Year = strconv.Itoa(booking.Year)
		payment.Amount = booking.Total

		// Clear sensitivite data in booking
		booking.CCNumber = "0000-0000-0000-0000"
		booking.CVV = "000"
		booking.Month = 99
		booking.Year = 0000

		// Invoke payment
		resp := makePayment(&payment)
		log.Println("response Status:", resp.Status)
		log.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println("response Body:", string(body))

		if resp.StatusCode == http.StatusOK {
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

			// Create a QR code
			err := qrcode.WriteFile("https://store.kings.cam.ac.uk/bookings/"+booking.UUID, qrcode.Medium, 256, booking.UUID+".png")

			// Send an email
			sendmail(&booking)
			if err != nil {
				log.Println("Failed to create a QR code: ", err)
			}

			// Write response
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Location", r.URL.Path)
			w.WriteHeader(http.StatusCreated)
			return
		}
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
