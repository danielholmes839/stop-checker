package travel

import "time"

type Route []*FixedLeg

type Schedule []*Leg

func (s Schedule) Arrival() time.Time {
	arrival := s[len(s)-1]
	return arrival.Departure.Add(arrival.Duration)
}

func (s Schedule) Departure() time.Time {
	return s[0].Departure
}

func (s Schedule) Duration() time.Duration {
	return s.Arrival().Sub(s.Departure())
}

type SchedulePlanner interface {
	Depart(at time.Time, route Route) (Schedule, error)
	Arrive(by time.Time, route Route) (Schedule, error)
}

type RoutePlanner interface {
	Depart(at time.Time, origin, destination string) (Route, error)
	Arrive(by time.Time, origin, destination string) (Route, error)
}
