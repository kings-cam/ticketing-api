package tickets

import (
	// Paths
	"os"
	// MongoDB driver
	"gopkg.in/mgo.v2"
)

// Database struct
type DB struct {
	Session *mgo.Session
}

// Connect to MongoDB instance
func (db *DB) Dial() (s *mgo.Session, err error) {
	return mgo.Dial(DBUrl())
}

// Create a DB name tickets
func (db *DB) Name() string {
	return "tickets"
}

// Return MongoDB URL
func DBUrl() string {
	// Try to fetch MongoDB URL
	db_url := os.Getenv("MONGOHQ_URL")

	if db_url == "" {
		db_url = "localhost"
	}

	return db_url
}
