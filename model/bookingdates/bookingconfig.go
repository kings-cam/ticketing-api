package bookingconfig

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type BookingConfig struct {
	// Start date for booking
	StartDate string `json:"startdate, omitempty" bson:"startdate, omitempty"`
	// Number of days from start date
	EndDays int `json:"enddays, omitempty" bson:"enddays, omitempty"`
	// Exclude days
	ExcludeDays []string `json:"excludedays, omitempty" bson:"excludedays, omitempty"`
}

type bookingconfig BookingConfig
