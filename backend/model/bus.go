package model

import "time"

type Bus struct {
	Headsign    string
	Arrival     time.Time
	LastUpdated time.Time
	Location    *Location
}
