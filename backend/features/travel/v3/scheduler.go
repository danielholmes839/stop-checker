package v3

import (
	"time"

	"stop-checker.com/db"
	"stop-checker.com/db/model"
	"stop-checker.com/db/repository"
)

type Scheduler struct {
	directions      walkingDirections
	directionsCache walkingDirectionsCache
	reach           db.ReachIndex
	stops           repository.Stops
}

func (s *Scheduler) Arrive(plan *model.TravelPlan, by time.Time) (*model.TravelSchedule, error) {
	return nil, nil
}

func (s *Scheduler) arrive(plan *model.TravelPlan, by time.Time) (*model.TravelSchedule, error) {
	return nil, nil
}

func (s *Scheduler) Depart(plan *model.TravelPlan, at time.Time) (*model.TravelSchedule, error) {
	original, err := s.depart(plan, at)
	if err != nil {
		return nil, err
	}

	backward, err := s.arrive(plan, original.DestinationArrival)
	if err != nil {
		return nil, err
	}

	return s.depart(plan, backward.OriginDeparture)
}

func (s *Scheduler) depart(plan *model.TravelPlan, at time.Time) (*model.TravelSchedule, error) {
	// if len(plan.Legs) == 0 {
	// 	return s.departZeroLegs(plan, at)
	// }

	// previousId := "INITIAL"
	// previousLocation := plan.Origin

	// legs := []model.Leg{}

	// for i, leg := range plan.Legs {
	// 	if leg.OriginId != previousId {
	// 		legs = append(legs, model.Leg{
	// 			OriginId:       previousId,
	// 			OriginLocation: previousLocation,
	// 			OriginArrival:  at,
	// 			DestinationId:  leg.OriginId,
	// 		})
	// 	}
	// }

	// schedule := &model.TravelSchedule{}

	return nil, nil
}

func (s *Scheduler) departZeroLegs(plan *model.TravelPlan, at time.Time) (*model.TravelSchedule, error) {
	path, err := s.directions.GetDirections(plan.Origin, plan.Destination)
	if err != nil {
		return nil, err
	}

	duration := walkingDuration(path.Distance)
	originArrival := at
	destinationArrival := at.Add(duration)

	return createZeroLegSchedule(plan, originArrival, destinationArrival, duration, path), nil
}

func (s *Scheduler) arriveZeroLegs(plan *model.TravelPlan, by time.Time) (*model.TravelSchedule, error) {
	path, err := s.directions.GetDirections(plan.Origin, plan.Destination)
	if err != nil {
		return nil, err
	}

	duration := walkingDuration(path.Distance)
	originArrival := by.Add(-duration)
	destinationArrival := by

	return createZeroLegSchedule(plan, originArrival, destinationArrival, duration, path), nil
}

func createZeroLegSchedule(plan *model.TravelPlan, originArrival, destinationArrival time.Time, duration time.Duration, path model.Path) *model.TravelSchedule {
	return &model.TravelSchedule{
		OriginDeparture:    originArrival,
		DestinationArrival: destinationArrival,
		Legs: []model.Leg{{
			OriginId:            "",
			OriginLocation:      plan.Origin,
			OriginArrival:       originArrival,
			DestinationId:       "",
			DestinationLocation: plan.Destination,
			DestinationArrival:  destinationArrival,
			Transit:             nil,
			Walk:                &path,
			Duration:            duration,
		}},
	}
}
