package main

import (
	// Ticketing package
	"tickets"

	// Routes
	"net/http"

	// Mongodb
	"gopkg.in/mgo.v2"
	// Gorilla Mux
	"github.com/gorilla/mux"
	// Negroni framework
	"github.com/urfave/negroni"
)

const port string = ":4000"

// Server API
func main() {
	mux := mux.NewRouter()

	// Includes some default middlewares
	n := negroni.Classic()
	n.UseHandler(mux)

	// Recovery
	n.Use(negroni.NewRecovery())

	// Logger
	n.Use(negroni.NewLogger())

	/*
	//For production, keep HTTPSProtection = true
	HTTPSProtection := false
	if HTTPSProtection {
		app.Use(restgate.New("X-Auth-Key", "X-Auth-Secret", restgate.Static, restgate.Config{HTTPSProtectionOff: false, Key: []string{c.API_ENDPOINT_KEY}, Secret: []string{c.API_ENDPOINT_SECRET}}))
        */
	
	/*
	// CORS for cross-domain access controls
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
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
	mux.HandleFunc("/api/v1/dates", tickets.BookingDates(session)).Methods("GET")

	// Config Booking dates
	mux.HandleFunc("/api/v1/dates/config", tickets.ConfigBookingDates(session)).Methods("POST")
	
	http.ListenAndServe(port, n)
}

