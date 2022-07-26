package travel

import "time"

type Plan []*FixedLeg

type Schedule []*Leg

type RouteScheduler interface {
	Depart(at time.Time, plan Plan) (Schedule, error)
	Arrive(by time.Time, plan Plan) (Schedule, error)
}

type RoutePlanner interface {
	Depart(at time.Time, origin, destination string) (Plan error)
}
