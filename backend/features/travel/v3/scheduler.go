package v3

import (
	"time"

	"stop-checker.com/db/model"
	"stop-checker.com/db/repository"
)

type Scheduler struct {
	*edgeFactory
}

func NewScheduler(
	directions walkingDirections,
	directionsCache walkingDirectionsCache,
	stopIndex repository.Stops,
	reachIndex repository.ReachBetween,
	stopTimesByTrip repository.InvertedIndex[model.StopTime],
) *Scheduler {
	return &Scheduler{
		edgeFactory: &edgeFactory{
			directions:      directions,
			directionsCache: directionsCache,
			stops:           stopIndex,
			reach: &scheduleReachImpl{
				reachIndex:      reachIndex,
				stopTimesByTrip: stopTimesByTrip,
				stopIndex:       stopIndex,
			},
		},
	}
}

func (s *Scheduler) Depart(plan *model.TravelPlan, at time.Time) (*model.TravelSchedule, error) {
	edges, err := s.Edges(plan)
	if err != nil {
		return nil, err
	}

	// create an initial schedule - departure time may be earlier than necessary
	initial, err := s.depart(edges, at)

	if err != nil {
		return nil, err
	}

	return s.arrive(edges, initial.DestinationArrival)
}

func (s *Scheduler) Arrive(plan *model.TravelPlan, by time.Time) (*model.TravelSchedule, error) {
	edges, err := s.Edges(plan)
	if err != nil {
		return nil, err
	}

	// create an initial schedule - arrival time may be later than necessary
	initial, err := s.arrive(edges, by)
	if err != nil {
		return nil, err
	}

	return s.depart(edges, initial.OriginDeparture)
}

func (s *Scheduler) depart(edges []scheduleEdge, at time.Time) (*model.TravelSchedule, error) {
	acc := at
	legs := []model.Leg{}

	for _, edge := range edges {
		leg, err := edge.Depart(acc)
		if err != nil {
			return nil, err
		}
		legs = append(legs, leg)
		acc = leg.DestinationArrival
	}

	return &model.TravelSchedule{
		OriginDeparture:    at,
		DestinationArrival: acc,
		Legs:               legs,
	}, nil
}

func (s *Scheduler) arrive(edges []scheduleEdge, by time.Time) (*model.TravelSchedule, error) {
	acc := by
	legs := []model.Leg{}

	for i := len(edges) - 1; i >= 0; i-- {
		leg, err := edges[i].Arrive(acc)
		if err != nil {
			return nil, err
		}
		legs = append(legs, leg)
		acc = leg.OriginArrival
	}

	for i, j := 0, len(legs)-1; i < j; i, j = i+1, j-1 {
		legs[i], legs[j] = legs[j], legs[i]
	}

	return &model.TravelSchedule{
		OriginDeparture:    acc,
		DestinationArrival: by,
		Legs:               legs,
	}, nil
}
