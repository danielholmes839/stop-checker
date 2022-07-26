package travel

import (
	"fmt"
	"math"
	"time"

	"stop-checker.com/db"
	"stop-checker.com/db/model"
)

type SchedulerConfig struct {
	StopIndex         *db.Index[model.Stop]
	StopTimesFromTrip *db.InvertedIndex[model.StopTime]
	ScheduleIndex     *db.ScheduleIndex
}

type Scheduler struct {
	stopIndex         *db.Index[model.Stop]
	stopTimesFromTrip *db.InvertedIndex[model.StopTime]
	scheduleIndex     *db.ScheduleIndex
}

func NewScheduler(config *SchedulerConfig) *Scheduler {
	return &Scheduler{
		stopIndex:         config.StopIndex,
		stopTimesFromTrip: config.StopTimesFromTrip,
		scheduleIndex:     config.ScheduleIndex,
	}
}

func (p *Scheduler) Depart(at time.Time, plan Plan) (Schedule, error) {
	planned := []*Leg{}
	acc := at

	for _, leg := range plan {
		plan, err := p.planDepart(acc, leg)
		if err != nil {
			return nil, err
		}
		acc = plan.Departure.Add(plan.Duration)
		planned = append(planned, plan)
	}

	return planned, nil
}

func (p *Scheduler) Arrive(by time.Time, plan Plan) (Schedule, error) {
	schedule, err := p.arrive(by, plan)
	if err != nil {
		return nil, err
	}

	first := schedule[0]

	optimized, err := p.Depart(first.Departure, plan)
	if err != nil {
		return nil, err
	}

	return optimized, nil
}

func (p *Scheduler) arrive(by time.Time, fixed []*FixedLeg) ([]*Leg, error) {
	planned := []*Leg{}
	acc := by

	// iterate over fixed legs in review
	for i := len(fixed) - 1; i >= 0; i-- {
		leg := fixed[i]
		plan, err := p.planArrive(acc, leg)

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

func (p *Scheduler) planDepart(acc time.Time, fixed *FixedLeg) (*Leg, error) {
	// origin and destination stops
	origin, destination, err := p.stops(fixed)
	if err != nil {
		return nil, err
	}

	if fixed.Walk {
		// planned leg by walking
		distance := origin.Distance(destination.Location)
		duration := time.Duration(math.Round(distance*1.4/60)) * time.Minute

		return &Leg{
			Origin:      fixed.Origin,
			Destination: fixed.Destination,
			Walk:        true,
			Departure:   acc,
			Duration:    duration,
		}, nil
	}

	// planned leg by transit
	next, err := p.scheduleIndex.Get(fixed.Origin, fixed.RouteId).Next(acc)
	if err != nil {
		return nil, err
	}

	// all stop times next trip
	all, _ := p.stopTimesFromTrip.Get(next.TripId)

	// origin stop times
	originArrival, err := p.stopTime(fixed.Origin, all)
	if err != nil {
		return nil, err
	}

	// destination stop times
	destinationArrival, err := p.stopTime(fixed.Destination, all)
	if err != nil {
		return nil, err
	}

	waitDuration := stopTimeDiffDuration(acc, originArrival.Time)
	transitDuration := stopTimeDiffDuration(originArrival.Time, destinationArrival.Time)
	departure := acc.Add(waitDuration)

	// planned leg
	return &Leg{
		Origin:      fixed.Origin,
		Destination: fixed.Destination,
		Walk:        false,
		Departure:   departure,
		Duration:    transitDuration,
		transit: &transit{
			TripId:                next.TripId,
			OriginStopTimeId:      originArrival.ID(),
			DestinationStopTimeId: destinationArrival.ID(),
		},
	}, nil
}

func (p *Scheduler) planArrive(acc time.Time, fixed *FixedLeg) (*Leg, error) {
	// get origin and destination stops
	origin, destination, err := p.stops(fixed)
	if err != nil {
		return nil, err
	}

	if fixed.Walk {
		// planned leg by walking
		distance := origin.Distance(destination.Location)
		duration := time.Duration(math.Round(distance*1.4/60)) * time.Minute

		return &Leg{
			Origin:      fixed.Origin,
			Destination: fixed.Destination,
			Walk:        true,
			Departure:   acc.Add(-duration),
			Duration:    duration,
		}, nil
	}

	// planned leg by transit
	previous, err := p.scheduleIndex.Get(fixed.Destination, fixed.RouteId).Previous(acc)
	if err != nil {
		return nil, err
	}

	// all stop times next trip
	all, _ := p.stopTimesFromTrip.Get(previous.TripId)

	// origin stop times
	originArrival, err := p.stopTime(fixed.Origin, all)
	if err != nil {
		return nil, err
	}

	// destination stop times
	destinationArrival, err := p.stopTime(fixed.Destination, all)
	if err != nil {
		return nil, err
	}

	excess := stopTimeDiffDuration(destinationArrival.Time, acc)
	transitDuration := stopTimeDiffDuration(originArrival.Time, destinationArrival.Time)

	// planned leg
	return &Leg{
		Origin:      fixed.Origin,
		Destination: fixed.Destination,
		Walk:        false,
		Departure:   acc.Add(-(transitDuration + excess)),
		Duration:    transitDuration,
		transit: &transit{
			TripId:                previous.TripId,
			OriginStopTimeId:      originArrival.ID(),
			DestinationStopTimeId: destinationArrival.ID(),
		},
	}, nil
}

// helper to get origin and destination stops
func (p *Scheduler) stops(fixed *FixedLeg) (model.Stop, model.Stop, error) {
	empty := model.Stop{}
	origin, ok := p.stopIndex.Get(fixed.Origin)
	if !ok {
		return empty, empty, fmt.Errorf("origin stop %s not found", fixed.Origin)
	}

	destination, ok := p.stopIndex.Get(fixed.Destination)
	if !ok {
		return empty, empty, fmt.Errorf("destination stop %s not found", fixed.Destination)
	}
	return origin, destination, nil
}

// helper to get stop time from a trip
func (p *Scheduler) stopTime(stopId string, all []model.StopTime) (model.StopTime, error) {
	for _, stopTime := range all {
		if stopTime.StopId == stopId {
			return stopTime, nil
		}
	}
	return model.StopTime{}, fmt.Errorf("stoptime not found stop:%s", stopId)
}
