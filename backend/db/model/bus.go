package model

type Bus struct {
	Destination Stop
	Headsign    string
	Arrival     Time
	LastUpdated Time
	Location    *Location
}
