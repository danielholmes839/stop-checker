package v2

import (
	"time"

	"stop-checker.com/db/model"
)

type FixedLegTransit struct {
	Origin      string // stop id
	Destination string // stop id
	Route       string // route id
}

type FixedLeg struct {
	Origin      model.Location
	Destination model.Location
	Transit     *FixedLegTransit
}

type LegTransit struct {
	Origin      model.Stop
	Destination model.Stop
	Route       model.Route
	Trip        model.Trip
	Arrival     time.Time
	Departure   time.Time
}

type Leg struct {
	Origin             model.Location
	OriginArrival      time.Time
	Destination        model.Location
	DestinationArrival time.Time
	Transit            *LegTransit
}
