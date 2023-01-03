package travel

import (
	"time"

	"stop-checker.com/db/model"
)

type scheduleNode struct {
	Id string
	model.Location
}

type edge struct {
	origin, destination *scheduleNode
}

func (e *edge) Edge() *edge {
	return e
}

/* scheduleEdge core interface */
type scheduleEdge interface {
	Edge() *edge                                          // get the origin and destination information of the end
	Depart(at time.Time) (model.TravelScheduleLeg, error) // create a model leg by departing from the origin at a specific time
	Arrive(by time.Time) (model.TravelScheduleLeg, error) // create a model leg arriving to the destination by a specific time
}

type scheduleWalkEdge struct {
	*edge
	path     *model.Path
	duration time.Duration
}

func (s *scheduleWalkEdge) Depart(at time.Time) (model.TravelScheduleLeg, error) {
	return model.TravelScheduleLeg{
		Origin: model.TravelScheduleNode{
			Id:       s.origin.Id,
			Location: s.origin.Location,
			Arrival:  at,
		},
		Destination: model.TravelScheduleNode{
			Id:       s.destination.Id,
			Location: s.destination.Location,
			Arrival:  at.Add(s.duration),
		},
		Transit: nil,
		Walk:    s.path,
	}, nil
}

// Arrive at the destination of the edge by a certain time
func (s *scheduleWalkEdge) Arrive(by time.Time) (model.TravelScheduleLeg, error) {
	return model.TravelScheduleLeg{
		Origin: model.TravelScheduleNode{
			Id:       s.origin.Id,
			Location: s.origin.Location,
			Arrival:  by.Add(-s.duration),
		},
		Destination: model.TravelScheduleNode{
			Id:       s.destination.Id,
			Location: s.destination.Location,
			Arrival:  by,
		},
		Transit: nil,
		Walk:    s.path,
	}, nil
}

type scheduleTransitEdge struct {
	*edge
	routeId string
	reach   scheduleReach
}

func (s *scheduleTransitEdge) Depart(at time.Time) (model.TravelScheduleLeg, error) {
	res, err := s.reach.Depart(s.origin.Id, s.destination.Id, s.routeId, at)
	if err != nil {
		return model.TravelScheduleLeg{}, err
	}

	return model.TravelScheduleLeg{
		Origin: model.TravelScheduleNode{
			Id:       s.origin.Id,
			Location: s.origin.Location,
			Arrival:  at,
		},
		Destination: model.TravelScheduleNode{
			Id:       s.destination.Id,
			Location: s.destination.Location,
			Arrival:  res.destinationArrival,
		},
		Transit: &model.Transit{
			TripId:          res.tripId,
			TripDuration:    res.destinationArrival.Sub(res.originDeparture),
			RouteId:         s.routeId,
			OriginDeparture: res.originDeparture,
		},
		Walk: nil,
	}, nil
}

func (s *scheduleTransitEdge) Arrive(by time.Time) (model.TravelScheduleLeg, error) {
	res, err := s.reach.Arrive(s.origin.Id, s.destination.Id, s.routeId, by)
	if err != nil {
		return model.TravelScheduleLeg{}, err
	}

	return model.TravelScheduleLeg{
		Origin: model.TravelScheduleNode{
			Id:       s.origin.Id,
			Location: s.origin.Location,
			Arrival:  res.originDeparture,
		},
		Destination: model.TravelScheduleNode{
			Id:       s.destination.Id,
			Location: s.destination.Location,
			Arrival:  res.destinationArrival,
		},
		Transit: &model.Transit{
			TripId:          res.tripId,
			TripDuration:    res.destinationArrival.Sub(res.originDeparture),
			RouteId:         s.routeId,
			OriginDeparture: res.originDeparture,
		},
		Walk: nil,
	}, nil
}
