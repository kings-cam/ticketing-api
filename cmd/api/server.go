package main

import (
	// Ticketing package
	"tickets"

	"encoding/json"
	// "fmt"
	"net/http"
	"log"
	"time"

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
	// Recovery and logging
	n := negroni.Classic()

	// Stats middleware
	statsmiddleware := stats.New()
	n.Use(statsmiddleware)

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


	/*
	// Welcome
	mux.HandleFunc("/api/v1/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Welcome to Ticketing API version 1.0\n")
	})
        */
	
	// Stats
	mux.HandleFunc("/api/v1/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		stats := statsmiddleware.Data()
		b, _ := json.Marshal(stats)
		w.Write(b)
	})

	// Booking dates
	mux.HandleFunc("/api/v1/dates/", tickets.BookingDates(session)).Methods("GET")

	// Config Booking dates
	mux.HandleFunc("/api/v1/dates/config/", tickets.ConfigBookingDates(session)).Methods("POST")

	// listen and serve api
	n.UseHandler(mux)

	// Create and launch server
	log.Println("Launching web api in http://localhost"+port)
	server := &http.Server{
		Handler:      n,
		Addr:         "127.0.0.1" + port,
		// Enforce timeouts for servers
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}	
	log.Fatal(server.ListenAndServe())
}

