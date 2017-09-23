package dates

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type BookingDates struct {
	// Start date for booking
	StartDate string `json:"startdate, omitempty" bson:"startdate, omitempty"`
	// Number of days from start date
	EndDays int `json:"enddays, omitempty" bson:"enddays, omitempty"`
}

type bookingconfig BookingDates
