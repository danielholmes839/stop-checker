package v3

import (
	"fmt"
	"time"

	"stop-checker.com/db/model"
	"stop-checker.com/db/repository"
)

type scheduleReachResult struct {
	tripId             string
	originDeparture    time.Time
	destinationArrival time.Time
}

type scheduleReach interface {
	// Depart from the origin at a certain time
	Depart(originId, destinationId, routeId string, at time.Time) (*scheduleReachResult, error)
	// Arrive to the destination by a certain time
	Arrive(originId, destinationId, routeId string, by time.Time) (*scheduleReachResult, error)
}

type scheduleReachImpl struct {
	reachIndex      repository.ReachBetween
	stopTimesByTrip repository.InvertedIndex[model.StopTime]
	stopIndex       repository.Stops
}

// Depart from the origin at a certain time
func (s *scheduleReachImpl) Depart(originId, destinationId, routeId string, at time.Time) (*scheduleReachResult, error) {
	originSchedule, _ := s.reachIndex.ReachableBetweenWithSchedule(originId, destinationId, routeId)

	// planned leg by transit
	next, err := originSchedule.Next(at)
	if err != nil {
		return nil, err
	}

	// all stop times next trip
	all, _ := s.stopTimesByTrip.Get(next.TripId)

	// origin stop time
	originArrival, err := s.stoptime(originId, all)
	if err != nil {
		return nil, err
	}

	// destination stop time
	destinationArrival, err := s.stoptime(destinationId, all)
	if err != nil {
		return nil, err
	}

	duration := model.TimeDiff(originArrival.Time, destinationArrival.Time)

	return &scheduleReachResult{
		tripId:             next.TripId,
		destinationArrival: next.Time.Add(duration),
		originDeparture:    next.Time,
	}, nil
}

// Arrive to the destination by a certain time
func (s *scheduleReachImpl) Arrive(originId, destinationId, routeId string, by time.Time) (*scheduleReachResult, error) {
	_, destinationSchedule := s.reachIndex.ReachableBetweenWithSchedule(originId, destinationId, routeId)

	// planned leg by transit
	previous, err := destinationSchedule.Previous(by)
	if err != nil {
		return nil, err
	}

	// all stop times next trip
	all, _ := s.stopTimesByTrip.Get(previous.TripId)

	// origin stop time
	originArrival, err := s.stoptime(originId, all)
	if err != nil {
		return nil, err
	}

	// destination stop time
	destinationArrival, err := s.stoptime(destinationId, all)
	if err != nil {
		return nil, err
	}

	duration := model.TimeDiff(originArrival.Time, destinationArrival.Time)

	return &scheduleReachResult{
		tripId:             previous.TripId,
		destinationArrival: previous.Time,
		originDeparture:    previous.Time.Add(-duration),
	}, nil
}

// helper to get origin and destination stops
func (s *scheduleReachImpl) stops(originId, destinationId string) (model.Stop, model.Stop, error) {
	origin, err := s.stopIndex.Get(originId)
	if err != nil {
		return model.Stop{}, model.Stop{}, err
	}

	destination, err := s.stopIndex.Get(destinationId)
	if err != nil {
		return model.Stop{}, model.Stop{}, err
	}

	return origin, destination, nil
}

// helper to get stop time from a trip
func (s *scheduleReachImpl) stoptime(stopId string, all []model.StopTime) (model.StopTime, error) {
	for _, stopTime := range all {
		if stopTime.StopId == stopId {
			return stopTime, nil
		}
	}
	return model.StopTime{}, fmt.Errorf("stoptime not found stop:%s", stopId)
}
