package tickets

import (
	"encoding/json"
	"net/http"
	"os"

	// CORS
	"github.com/rs/cors"
	// JWT
	"github.com/dgrijalva/jwt-go"
	// JSON Web Tokens middleware Auth0
	"github.com/auth0/go-jwt-middleware"
	// Mongodb
	"gopkg.in/mgo.v2"
	// Gorilla Mux
	"github.com/gorilla/mux"
	// Negroni framework
	"github.com/urfave/negroni"
	// Stats
	"github.com/thoas/stats"
)

// Router Returns the mux api router with CORS, stats and logging
func Router() *mux.Router {
	apirouter := mux.NewRouter()

	/*
		//For production, keep HTTPSProtection = true
		HTTPSProtection := false
		if HTTPSProtection {
			app.Use(restgate.New("X-Auth-Key", "X-Auth-Secret", restgate.Static, restgate.Config{HTTPSProtectionOff: false, Key: []string{c.API_ENDPOINT_KEY}, Secret: []string{c.API_ENDPOINT_SECRET}}))
	*/

	// API Router
	apirouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ticketing API\n"))
	})

	return apirouter
}

func V1Router(apirouter *mux.Router) *mux.Router {
	// Stats middleware
	statsmw := stats.New()

	// CORS for cross-domain access controls
	corsmw := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "HEAD", "PUT", "PATCH", "POST", "DELETE"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
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

	// Stats
	apiv1router.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		stats := statsmw.Data()
		b, _ := json.Marshal(stats)
		w.Write(b)
	})

	return apiv1router
}

func V1CONFIGRouter(apirouter *mux.Router) *mux.Router {
	// Stats middleware
	statsmw := stats.New()

	// CORS for cross-domain access controls
	corsmw := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "HEAD", "PUT", "PATCH", "POST", "DELETE"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
	})

	// Auth0 JWT middleware
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("Auth0")), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	// API V1CONFIG router
	apiconfigrouter := mux.NewRouter().PathPrefix("/config").Subrouter().StrictSlash(true)
	apirouter.PathPrefix("/config").Handler(negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		statsmw,
		corsmw,
		negroni.Wrap(apiconfigrouter),
	))

	// Stats
	apiconfigrouter.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		stats := statsmw.Data()
		b, _ := json.Marshal(stats)
		w.Write(b)
	})

	return apiconfigrouter
}

// Routes define API version 1.0  router for tickets package
func Routes(apiv1router *mux.Router, session *mgo.Session) {
	// API version 1.0 welcome
	apiv1router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ticketing API version 1!\n"))
	})

	// Booking dates
	apiv1router.HandleFunc("/dates", BookingDates(session, false)).Methods("GET")

	// Get Session
	apiv1router.HandleFunc("/sessions/{date}", BookingSessions(session)).Methods("GET")

	// Get pricing
	apiv1router.HandleFunc("/prices", GetPrices(session)).Methods("GET")

	// Create a new booking
	apiv1router.HandleFunc("/bookings/{uuid}", CreateBooking(session)).Methods("POST")
}

// Routes define API version 1.0  router for tickets package
func ConfigRoutes(apiv1router *mux.Router, session *mgo.Session) {
	// API version 1.0 welcome
	apiv1router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ticketing API config version 1!\n"))
	})

	// Config Booking dates
	apiv1router.HandleFunc("/dates", ConfigBookingDates(session, false)).Methods("POST")

	// Get Config Booking dates
	apiv1router.HandleFunc("/dates", GetConfigDates(session, false)).Methods("GET")

	// Config Booking dates
	apiv1router.HandleFunc("/prices", ConfigPricing(session)).Methods("POST")

	// Test configuration
	// Test config Booking dates
	apiv1router.HandleFunc("/test/dates", ConfigBookingDates(session, true)).Methods("POST")

	// Test booking dates
	apiv1router.HandleFunc("/test/dates", BookingDates(session, true)).Methods("GET")

	// Get existing booking
	apiv1router.HandleFunc("/bookings/{uuid}", GetBooking(session)).Methods("GET")
	// Return all bookings
	apiv1router.HandleFunc("/bookings", GetBookings(session)).Methods("GET")

	// Return all bookings matching a date
	apiv1router.HandleFunc("/bookings/date/{date}", GetBookingsDate(session)).Methods("GET")

	// Update an existing booking
	apiv1router.HandleFunc("/bookings/{uuid}", UpdateBooking(session)).Methods("PUT")

	// Delete an existing booking
	apiv1router.HandleFunc("/bookings/{uuid}", DeleteBooking(session)).Methods("DELETE")

}
