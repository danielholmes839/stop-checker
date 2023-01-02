package model

import "time"

type TravelPlan struct {
	Origin      Location
	Destination Location
	Legs        []TravelPlanLeg
}

type TravelPlanLeg struct {
	OriginId      string
	DestinationId string
	RouteId       string
}

type TravelSchedule struct {
	OriginDeparture    time.Time
	DestinationArrival time.Time
	Legs               []Leg
}

type Leg struct {
	// origin attributes
	OriginId       string
	OriginLocation Location
	OriginArrival  time.Time // when we reach the origin

	// destination attributes
	DestinationId       string
	DestinationLocation Location
	DestinationArrival  time.Time // when we reach the destination

	// walking, transit information
	Transit *LegTransit
	Walk    *Path

	Duration time.Duration // the walking duration or bus duration
}

type LegTransit struct {
	TripId          string
	TripDuration    time.Duration // time spent on transit
	RouteId         string
	OriginDeparture time.Time
}
