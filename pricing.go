package tickets

import (
	"encoding/json"
	"log"
	"net/http"

	// Mongo DB
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Pricing struct {
	ID int `json:"id"`
	// Adult ticket price
	AdultPrice float32 `json:"adultprice"`
	// Child ticket price
	ChildPrice float32 `json:"childprice"`
	// Concession ticket price
	ConcessionPrice float32 `json:"concessionprice"`
	// Guidebook price
	GuidePrice float32 `json:"guideprice"`
	// Guide book laguages
	GuideBooks []string `json:"guidebooks"`
}


// ConfigPricing assigns ticket prices
func ConfigPricing(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		var price Pricing

		// Decode POST JSON file
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&price)
		if err != nil {
			ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			return
		}

		// Open price collection
		c := session.DB("tickets").C("pricing")

		// Try to update if price is found
		err = c.Update(bson.M{"id": 0}, &price)
		if err != nil {
			switch err {
			default:
				ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
				log.Println("Failed to update price: ", err)
				return
			// Price is not present, do an insert
			case mgo.ErrNotFound:
				log.Println("Price not present, creating a new price")
				err = c.Insert(&price)

				if err != nil {
					ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
					log.Println("Failed to insert price: ", err)
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
