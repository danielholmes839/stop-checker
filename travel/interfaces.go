package travel

import "time"

type Route []*FixedLeg

type Schedule []*Leg

type SchedulePlanner interface {
	Depart(at time.Time, route Route) (Schedule, error)
	Arrive(by time.Time, route Route) (Schedule, error)
}

type RoutePlanner interface {
	Depart(at time.Time, origin, destination string) (Route, error)
}
