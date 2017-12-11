package main

import (
	// Ticketing package
	"tickets"

	"log"
	"net/http"
	"os"
	"time"

	// Mongodb
	"gopkg.in/mgo.v2"
)

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

	// Create API v1 routes
	apiv1configrouter := tickets.V1CONFIGRouter(apirouter)
	tickets.ConfigRoutes(apiv1configrouter, session)

	// Create and launch server
	log.Println("Launching web api in " + os.Getenv("IP") + ":" + os.Getenv("Port"))
	server := &http.Server{
		Handler: apirouter,
		Addr:    os.Getenv("IP") + ":" + os.Getenv("Port"),
		// Enforce timeouts for servers
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(server.ListenAndServeTLS("fullchain.pem", "privkey.pem"))

}
