package main

import (
	// Routes
	"net/http"
	// "tickets"
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

	// Initialise routes
	//	tickets.Routes(mux)

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		panic("oh no")
	})
	
	http.ListenAndServe(port, n)
}

