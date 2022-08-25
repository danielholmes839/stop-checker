package model

import (
	"fmt"
	"math"
)

type StopRoute struct {
	RouteId     string
	StopId      string
	DirectionId string
	Headsign    string
}

func (r StopRoute) DirectedID() string {
	return fmt.Sprintf("%s:%s", r.RouteId, r.DirectionId)
}

// GraphQL related
type Transit struct {
	RouteId       string
	TripId        string
	OriginId      string
	DestinationId string
}

type Location struct {
	Latitude  float64
	Longitude float64
}

const earthRaidusKm float64 = 6371 // radius of the earth in kilometers.

func degreesToRadians(d float64) float64 {
	return d * math.Pi / 180
}

func (l Location) Distance(other Location) float64 {
	// harversine distance in meters
	lat1 := degreesToRadians(l.Latitude)
	lon1 := degreesToRadians(l.Longitude)
	lat2 := degreesToRadians(other.Latitude)
	lon2 := degreesToRadians(other.Longitude)

	diffLat := lat2 - lat1
	diffLon := lon2 - lon1

	a := math.Pow(math.Sin(diffLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*
		math.Pow(math.Sin(diffLon/2), 2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return c * 6_371_000
}
