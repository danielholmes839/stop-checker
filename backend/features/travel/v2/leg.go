package v2

import (
	"time"

	"stop-checker.com/db/model"
)

type Schedule []*Leg

type Route []*FixedLeg

type FixedLeg struct {
	OriginId      string
	Origin        model.Location
	DestinationId string
	Destination   model.Location
	RouteId       string
}

type Leg struct {
	OriginId      string
	Origin        model.Location
	OriginArrival time.Time

	DestinationId      string
	Destination        model.Location
	DestinationArrival time.Time

	Transit *Transit
}

type Transit struct {
	RouteId   string
	TripId    string
	Departure time.Time
}
