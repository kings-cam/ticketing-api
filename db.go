package tickets

import (
	"crypto/tls"
	"net"
	// Paths
	"os"
	// MongoDB driver
	"gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
)

// Database struct
type DB struct {
	Session *mgo.Session
}

// Connect to MongoDB instance
func (db *DB) Dial() (s *mgo.Session, err error) {
	// TLS
	tlsConfig := &tls.Config{}
	tlsConfig.InsecureSkipVerify = true

	// fmt.Println("MongoHost:", os.Getenv("MongoHost"))
	//connect URL:
	// "mongodb://<username>:<password>@<hostname>:<port>,<hostname>:<port>/<db-name>
	dialInfo, err := mgo.ParseURL("mongodb://" + os.Getenv("MongoUser") + ":" + os.Getenv("MongoPW") + "@kings-shard-00-00-zzes4.mongodb.net:" + os.Getenv("MongoPort") + ",kings-shard-00-01-zzes4.mongodb.net:" + os.Getenv("MongoPort") + ",kings-shard-00-02-zzes4.mongodb.net:" + os.Getenv("MongoPort") + "/")
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}

	session, err := mgo.DialWithInfo(dialInfo)
	return session, err
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

// EnsureIndex checks database for duplicates
func EnsureIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	dbconfig := session.DB("tickets").C("config")

	projectindex := mgo.Index{
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	projecterr := dbconfig.EnsureIndex(projectindex)
	if projecterr != nil {
		panic(projecterr)
	}


	dbdates := session.DB("tickets").C("dates")

	datesindex := mgo.Index{
		Key:        []string{"date"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	dateserr := dbdates.EnsureIndex(datesindex)
	if dateserr != nil {
		panic(dateserr)
	}
}
