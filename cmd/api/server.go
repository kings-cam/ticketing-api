package main

import (
	// Ticketing package
	"tickets"

	"encoding/json"
	"net/http"
	"log"

	// CORS
	"github.com/rs/cors"
	// Mongodb
	"gopkg.in/mgo.v2"
	// Gorilla Mux
	"github.com/gorilla/mux"
	// Negroni framework
	"github.com/urfave/negroni"
	// Stats
	"github.com/thoas/stats"
)

const port string = ":4000"

// Server API
func main() {
	mux := mux.NewRouter()

	// Includes some default middlewares
	n := negroni.Classic()

	// Recovery
	n.Use(negroni.NewRecovery())

	// Logger
	n.Use(negroni.NewLogger())
	// Stats middleware
	statsmiddleware := stats.New()

	// CORS for cross-domain access controls
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "HEAD", "PUT", "PATCH", "POST", "DELETE"},

	})
	n.Use(c)

	/*
	//For production, keep HTTPSProtection = true
	HTTPSProtection := false
	if HTTPSProtection {
		app.Use(restgate.New("X-Auth-Key", "X-Auth-Secret", restgate.Static, restgate.Config{HTTPSProtectionOff: false, Key: []string{c.API_ENDPOINT_KEY}, Secret: []string{c.API_ENDPOINT_SECRET}}))
        */
	
	// Create database and session
	db := tickets.DB{}
	session, err := db.Dial()

	if err != nil {
		panic(err)
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)


	// Welcome
	// mux.HandleFunc("/api/v1", welcome).Methods("GET")

	// Booking dates
	mux.HandleFunc("/api/v1/dates/", tickets.BookingDates(session)).Methods("GET")

	// Config Booking dates
	mux.HandleFunc("/api/v1/dates/config/", tickets.ConfigBookingDates(session)).Methods("POST")

	mux.HandleFunc("/api/v1/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		stats := statsmiddleware.Data()

		b, _ := json.Marshal(stats)

		w.Write(b)
	})

	n.Use(statsmiddleware)
	// listen and serve api
	n.UseHandler(mux)
	log.Println("Launching web api in http://localhost"+port)
	http.ListenAndServe(port, n)
}

