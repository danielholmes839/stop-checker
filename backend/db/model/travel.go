package model

import (
	"fmt"
	"time"
)

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
	Legs               []TravelScheduleLeg
}

type TravelScheduleNode struct {
	Id       string // stop id, "ORIGIN" or "DESTINATION"
	Location Location
	Arrival  time.Time
}

type TravelScheduleLeg struct {
	// origin attributes
	Origin TravelScheduleNode

	// destination attributes
	Destination TravelScheduleNode

	// walking, transit information
	Transit *Transit
	Walk    *Path
}

func (l *TravelScheduleLeg) String() string {
	if l.Transit == nil {
		return fmt.Sprintf("%s -> %s. %s, %s. %.2f", l.Origin.Id, l.Destination.Id, l.Origin.Arrival, l.Destination.Arrival, l.Walk.Distance)
	}
	return fmt.Sprintf("%s -> %s. %s, %s. %s %s", l.Origin.Id, l.Destination.Id, l.Origin.Arrival, l.Destination.Arrival, l.Transit.RouteId, l.Transit.OriginDeparture)
}

type Transit struct {
	TripId          string
	TripDuration    time.Duration // time spent on transit
	RouteId         string
	OriginDeparture time.Time
}
