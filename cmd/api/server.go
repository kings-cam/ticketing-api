package main

import (
	// Ticketing package
	"tickets"

	"encoding/json"
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
	apirouter := mux.NewRouter()

	// Stats middleware
	statsmw := stats.New()

	// CORS for cross-domain access controls
	corsmw := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "HEAD", "PUT", "PATCH", "POST", "DELETE"},

	})

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
	
	// API Router
	apirouter.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ticketing API\n"))
	})
	
	
	// API V1 router
	apiv1router := mux.NewRouter().PathPrefix("/api/v1").Subrouter().StrictSlash(true)
	apirouter.PathPrefix("/api/v1").Handler(negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		statsmw,
		corsmw,
		negroni.Wrap(apiv1router),
	))

	// API version 1.0 welcome
	apiv1router.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ticketing API version 1!\n"))
	})
	
	// Stats
	apiv1router.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		stats := statsmw.Data()
		b, _ := json.Marshal(stats)
		w.Write(b)
	})
        

	// Booking dates
	apiv1router.HandleFunc("/dates", tickets.BookingDates(session)).Methods("GET")

	// Config Booking dates
	apiv1router.HandleFunc("/dates/config", tickets.ConfigBookingDates(session)).Methods("POST")


	// Create and launch server
	log.Println("Launching web api in http://localhost"+port)
	server := &http.Server{
		Handler:      apirouter,
		Addr:         "127.0.0.1" + port,
		// Enforce timeouts for servers
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}	
	log.Fatal(server.ListenAndServe())
}

