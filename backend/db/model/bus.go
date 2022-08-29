package model

type Bus struct {
	Headsign    string
	Arrival     Time
	LastUpdated Time
	Location    *Location
}
