package main

import (
	// Ticketing package
	"tickets"

	"net/http"
	"log"
	"time"

	// Mongodb
	"gopkg.in/mgo.v2"
)

const port string = ":4000"

// main Runs the tickets api server
func main() {
	// Get router
	apirouter := tickets.Router()
	
	// Create database and session
	db := tickets.DB{}
	session, err := db.Dial()
	
	if err != nil {
		panic(err)
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	// Ensure Index
	tickets.EnsureIndex(session)

	// Create API v1 routes
	apiv1router := tickets.V1Router(apirouter)
	tickets.Routes(apiv1router, session)

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
