package travel

import (
	"fmt"
	"time"

	"stop-checker.com/db"
	"stop-checker.com/db/model"
)

type SchedulerConfig struct {
	StopIndex       *db.Index[model.Stop]
	ReachIndex      *db.ReachIndex
	StopTimesByTrip *db.InvertedIndex[model.StopTime]
}

type Scheduler struct {
	stopIndex       *db.Index[model.Stop]
	reachIndex      *db.ReachIndex
	stopTimesByTrip *db.InvertedIndex[model.StopTime]
}

func NewScheduler(config *SchedulerConfig) *Scheduler {
	return &Scheduler{
		stopIndex:       config.StopIndex,
		reachIndex:      config.ReachIndex,
		stopTimesByTrip: config.StopTimesByTrip,
	}
}

func (s *Scheduler) Depart(at time.Time, route Route) (Schedule, error) {
	schedule := Schedule{}
	acc := at

	for _, fixedLeg := range route {
		leg, err := s.planDepart(acc, fixedLeg)
		if err != nil {
			return nil, err
		}
		acc = leg.Departure.Add(leg.Duration)
		schedule = append(schedule, leg)
	}

	return schedule, nil
}

func (s *Scheduler) Arrive(by time.Time, route Route) (Schedule, error) {
	schedule, err := s.arrive(by, route)
	if err != nil {
		return nil, err
	}

	first := schedule[0]

	optimized, err := s.Depart(first.Departure, route)
	if err != nil {
		return nil, err
	}

	return optimized, nil
}

func (s *Scheduler) arrive(by time.Time, fixed []*FixedLeg) ([]*Leg, error) {
	planned := []*Leg{}
	acc := by

	// iterate over fixed legs in review
	for i := len(fixed) - 1; i >= 0; i-- {
		leg := fixed[i]
		plan, err := s.planArrive(acc, leg)

		if err != nil {
			return nil, err
		}

		// adjust the accumulator time
		acc = plan.Departure
		planned = append(planned, plan)
	}

	// reverse the list of planned legs
	order := []*Leg{}
	for i := 0; i < len(planned); i++ {
		order = append(order, planned[len(planned)-(i+1)])
	}

	return order, nil
}

func (s *Scheduler) planDepart(acc time.Time, fixed *FixedLeg) (*Leg, error) {
	// origin and destination stops
	origin, destination, err := s.stops(fixed)
	if err != nil {
		return nil, err
	}

	if fixed.Walk {
		// planned leg by walking
		distance := origin.Distance(destination.Location)
		duration := walkingDuration(distance)

		return &Leg{
			Origin:      fixed.Origin,
			Destination: fixed.Destination,
			Walk:        true,
			Departure:   acc,
			Duration:    duration,
		}, nil
	}

	// stop times guaranteed to reach between the origin and destination
	originSchedule, _ := s.reachIndex.ReachableBetweenWithSchedule(fixed.Origin, fixed.Destination, fixed.RouteId)

	// planned leg by transit
	next, err := originSchedule.Next(acc)
	if err != nil {
		return nil, err
	}

	// all stop times next trip
	all, _ := s.stopTimesByTrip.Get(next.TripId)

	// origin stop time
	_, err = s.stopTime(fixed.Origin, all)
	if err != nil {
		return nil, err
	}

	// destination stop times
	destinationArrival, err := s.stopTime(fixed.Destination, all)
	if err != nil {
		return nil, err
	}

	transitDuration := model.TimeDiff(next.StopTime.Time, destinationArrival.Time)
	departure := next.Time

	// planned leg
	return &Leg{
		Origin:      fixed.Origin,
		Destination: fixed.Destination,
		Walk:        false,
		Departure:   departure,
		Duration:    transitDuration,
		Transit: &transit{
			TripId: next.TripId,
		},
	}, nil
}

func (s *Scheduler) planArrive(acc time.Time, fixed *FixedLeg) (*Leg, error) {
	// get origin and destination stops
	origin, destination, err := s.stops(fixed)
	if err != nil {
		return nil, err
	}

	if fixed.Walk {
		// planned leg by walking
		distance := origin.Distance(destination.Location)
		duration := walkingDuration(distance)

		return &Leg{
			Origin:      fixed.Origin,
			Destination: fixed.Destination,
			Walk:        true,
			Departure:   acc.Add(-duration),
			Duration:    duration,
		}, nil
	}

	_, destinationSchedule := s.reachIndex.ReachableBetweenWithSchedule(fixed.Origin, fixed.Destination, fixed.RouteId)

	// planned leg by transit
	previous, err := destinationSchedule.Previous(acc)
	if err != nil {
		return nil, err
	}

	// all stop times next trip
	all, _ := s.stopTimesByTrip.Get(previous.TripId)

	// origin stop time
	originArrival, err := s.stopTime(fixed.Origin, all)
	if err != nil {
		return nil, err
	}

	// destination stop time
	destinationArrival, err := s.stopTime(fixed.Destination, all)
	if err != nil {
		return nil, err
	}

	transitDuration := model.TimeDiff(originArrival.Time, destinationArrival.Time)

	// planned leg
	return &Leg{
		Origin:      fixed.Origin,
		Destination: fixed.Destination,
		Walk:        false,
		Departure:   previous.Add(-transitDuration),
		Duration:    transitDuration,
		Transit: &transit{
			TripId: previous.TripId,
		},
	}, nil
}

// helper to get origin and destination stops
func (s *Scheduler) stops(fixed *FixedLeg) (model.Stop, model.Stop, error) {
	empty := model.Stop{}
	origin, err := s.stopIndex.Get(fixed.Origin)
	if err != nil {
		return empty, empty, fmt.Errorf("origin stop %s not found", fixed.Origin)
	}

	destination, err := s.stopIndex.Get(fixed.Destination)
	if err != nil {
		return empty, empty, fmt.Errorf("destination stop %s not found", fixed.Destination)
	}
	return origin, destination, nil
}

// helper to get stop time from a trip
func (s *Scheduler) stopTime(stopId string, all []model.StopTime) (model.StopTime, error) {
	for _, stopTime := range all {
		if stopTime.StopId == stopId {
			return stopTime, nil
		}
	}
	return model.StopTime{}, fmt.Errorf("stoptime not found stop:%s", stopId)
}
