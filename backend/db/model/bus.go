package model

import "time"

type Bus struct {
	Destination Stop
	Headsign    string
	Arrival     time.Time
	LastUpdated time.Time
	Location    *Location
}
