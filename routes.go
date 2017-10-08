package tickets

import (
	"encoding/json"
	"net/http"

	// CORS
	"github.com/rs/cors"
	// JWT
	// "github.com/dgrijalva/jwt-go"
	// JSON Web Tokens middleware Auth0
	// "github.com/auth0/go-jwt-middleware"
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
	apirouter.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ticketing API\n"))
	})
        
    return apirouter
}

func V1Router(apirouter *mux.Router) *mux.Router {
	// Stats middleware
	statsmw := stats.New()

	// CORS for cross-domain access controls
	corsmw := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "HEAD", "PUT", "PATCH", "POST", "DELETE"},

	})

/*
	// Auth0 JWT middleware
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("dsAREIWkSHee604VTbq4kJf0imEeWwdC"), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
*/
	
	// API V1 router
	apiv1router := mux.NewRouter().PathPrefix("/api/v1").Subrouter().StrictSlash(true)
	apirouter.PathPrefix("/api/v1").Handler(negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		// negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
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

// Routes define API version 1.0  router for tickets package
func Routes(apiv1router *mux.Router, session *mgo.Session) {
	// API version 1.0 welcome
	apiv1router.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ticketing API version 1!\n"))
	})

	// Booking dates
	apiv1router.HandleFunc("/dates", BookingDates(session, false)).Methods("GET")

	// Get Session
	apiv1router.HandleFunc("/sessions/{date}", BookingSessions(session)).Methods("GET")

	// Get pricing
	apiv1router.HandleFunc("/prices", GetPrices(session)).Methods("GET")
	
	
	// Config Booking dates
	apiv1router.HandleFunc("/config/dates", ConfigBookingDates(session, false)).Methods("POST")

	// Config Booking dates
	apiv1router.HandleFunc("/config/prices", ConfigPricing(session)).Methods("POST")

	// Test config Booking dates
	apiv1router.HandleFunc("/test/config/dates", ConfigBookingDates(session, true)).Methods("POST")

	// Test booking dates
	apiv1router.HandleFunc("/test/dates", BookingDates(session, true)).Methods("GET")
}
