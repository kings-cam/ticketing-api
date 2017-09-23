package tickets

import (
	"fmt"
	// HTTP requests
	"net/http"
	// "time"
	// "encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	// Echo webframework
	"github.com/labstack/echo"
)


type BookingConfig struct {
	ID int `json:"_id"`
	// Start date for booking
	// StartDate string `json:"startdate, omitempty"`
	// Number of days from start date
	// EndDays int `json:"enddays, omitempty"`
	// Exclude days
	// ExcludeDays []string `json:"excludedays, omitempty"`
	// Booking days
	BookingDays []string `json:"bookingdays, omitempty"`
}

func ConfigBookingDates(c echo.Context) (err error) {
	
	config := new(BookingConfig)
	if err = c.Bind(config); err != nil {
		return
	}
	
	db := DB{}
	session, err := db.Dial()
	
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	
	dbc := session.DB("tickets").C("config")

	err = dbc.Insert(config)
	if err != nil {
		if mgo.IsDup(err) {
			panic(err)
			return
		}
		panic(err)
		fmt.Println("Failed insert book: ", err)
		return
	}
	return c.JSON(http.StatusOK, config)
}

func BookingDates() []string {

	db := DB{}
	session, err := db.Dial()
	
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	var bookingdates []string


	c := session.DB("tickets").C("config")
	var config BookingConfig

	err = c.Find(bson.M{"_id": 0}).One(config)
	if err != nil {
		fmt.Println("Failed find book: ", err)
	}

	if len(config.BookingDays) == 0 {
		fmt.Println("Booking dates are empty")
	}

	//respBody, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(config.BookingDays)

	bookingdates = config.BookingDays
	/*
	// Start date as tomorrow
	startdate := time.Now().Local().AddDate(0, 0, 1)

	// End date as 90 days (3 months) from tomorrow
	enddate := startdate.AddDate(0, 0, 90)

	// Iterate over dates to print all allowed dates
	for d := startdate; d != enddate; d = d.AddDate(0, 0, 1) {
		// Exclude weekends (0 - Sunday, 6 - Saturday)
		if d.Weekday() != 0 {
			bookingdates = append(bookingdates, d.Format("2006-01-02"))
		}
	}
        */
	return bookingdates
}

// Return dates of bookings
func GetBookingDates(c echo.Context) error {
	// bd, _ to get errors
	return c.JSON(http.StatusCreated, BookingDates)
}
