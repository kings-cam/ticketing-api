package main

import (
	// Routes
	"ticketing-api/routes"
	// Echo framework
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const port string = ":4000"

// Server API
func main() {
	e := echo.New()
	
	// Middleware
	// Server log
	e.Use(middleware.Logger())
	// Recovers from panics anywhere in the chain
	e.Use(middleware.Recover())

	// CORS for cross-domain access controls
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	// Initialise routes
	routes.Init(e)
	
	// Server
	e.Logger.Fatal(e.Start(":4000"))
}
