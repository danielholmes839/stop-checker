package travel

import (
	"errors"
	"fmt"
	"time"

	"stop-checker.com/db"
	"stop-checker.com/db/model"
)

type FixedEdge interface {
	Origin() string
	OriginArrival() time.Time
	OriginWait() time.Duration // time spent waiting at the origin
	Destination() string
	DestinationArrival() time.Time
	Transit() string // "WALK" or "TRANSIT"
	TransitDuration() time.Duration
	TotalDuration() time.Duration // transit + origin wait
}

type FixedLeg struct {
	Origin      string // stop id
	Destination string // stop id
	RouteId     string // stop id
	Walk        bool   // if true then: route id is empty
}

type TransitLeg struct {
	TripId                string
	OriginStopTimeId      string
	DestinationStopTimeId string
	WaitDuration          time.Duration
	TransitDuration       time.Duration
}

type WalkingLeg struct {
	Distance float64
	Duration time.Duration
}

func (l *Leg) String() string {
	if l.Walk {
		return fmt.Sprintf("walk {from:%s, to:%s, arrival:%s}",
			l.Origin, l.Destination, l.Arrival.Format("2006-01-02@15:04"),
		)
	}

	return fmt.Sprintf("transit {from:%s, to:%s, arrival:%s, wait:%s, transit:%s}",
		l.Origin, l.Destination, l.Arrival.Format("2006-01-02@15:04"), l.WaitDuration, l.TransitDuration,
	)
}

type Leg struct {
	Arrival time.Time
	*FixedLeg
	*TransitLeg
	*WalkingLeg
}

func (leg *Leg) Duration() time.Duration {
	if leg.Walk {
		return leg.WalkingLeg.Duration
	}
	return leg.WaitDuration + leg.TransitDuration
}

type FixedPlanner struct {
	StopIndex         *db.Index[model.Stop]
	StopTimesFromTrip *db.InvertedIndex[model.StopTime]
	*db.ScheduleIndex
}

func (p *FixedPlanner) DepartAt(at time.Time, fixedLegs []*FixedLeg) ([]*Leg, error) {
	legs := []*Leg{}

	for _, fixed := range fixedLegs {
		leg, err := p.NewLeg(at, fixed)
		if err != nil {
			return nil, err
		}
		at = at.Add(leg.Duration())
		legs = append(legs, leg)
	}

	return legs, nil
}

func (p *FixedPlanner) ArriveBy(by time.Time, fixedLegs []*FixedLeg) ([]*Leg, error) {
	legs := []*Leg{}

	for _, fixed := range fixedLegs {
		leg, err := p.NewLeg(by, fixed)
		if err != nil {
			return nil, err
		}
		by = by.Add(leg.Duration())
		legs = append(legs, leg)
	}

	return legs, nil
}

func (p *FixedPlanner) NewLeg(arrival time.Time, fixed *FixedLeg) (*Leg, error) {
	leg := &Leg{Arrival: arrival, FixedLeg: fixed}
	if leg.Walk {
		// create a walking leg
		walkingLeg, err := p.NewWalkingLeg(arrival, fixed)
		leg.WalkingLeg = walkingLeg
		if err != nil {
			return nil, err
		}
		return leg, nil
	}

	// create a transit leg
	transitLeg, err := p.NewTransitLeg(arrival, fixed)
	leg.TransitLeg = transitLeg
	if err != nil {
		return nil, err
	}
	return leg, nil
}

func (p *FixedPlanner) NewTransitLeg(arrival time.Time, fixed *FixedLeg) (*TransitLeg, error) {
	// orgin stop time
	originStopTime, err := p.ScheduleIndex.Get(fixed.Origin, fixed.RouteId).Next(arrival)
	if err != nil {
		return nil, err
	}

	// stop times for the trip
	stopTimes, ok := p.StopTimesFromTrip.Get(originStopTime.TripId)
	if !ok {
		return nil, errors.New("no stop times")
	}

	// destination stop time
	destinationStopTime := model.StopTime{}
	destinationStopTimeFound := false

	for _, stopTime := range stopTimes {
		if stopTime.StopId == fixed.Destination {
			destinationStopTime = stopTime
			destinationStopTimeFound = true
			break
		}
	}

	// destination stop time not found
	if !destinationStopTimeFound {
		return nil, errors.New("route does not reach destination")
	}

	waitDuration := stopTimeDiffDuration(arrival, originStopTime.Time)
	transitDuration := stopTimeDiffDuration(originStopTime.Time, destinationStopTime.Time)

	return &TransitLeg{
		WaitDuration:          waitDuration,
		TransitDuration:       transitDuration,
		OriginStopTimeId:      originStopTime.ID(),
		DestinationStopTimeId: destinationStopTime.ID(),
		TripId:                originStopTime.TripId,
	}, nil
}

func (p *FixedPlanner) NewWalkingLeg(arrival time.Time, fixed *FixedLeg) (*WalkingLeg, error) {
	origin, ok := p.StopIndex.Get(fixed.Origin)
	if !ok {
		return nil, fmt.Errorf("stop:%s not found", fixed.Origin)
	}

	destination, ok := p.StopIndex.Get(fixed.Destination)
	if !ok {
		return nil, fmt.Errorf("stop:%s not found", fixed.Destination)
	}

	// calculate distance and duration
	distance := origin.Location.Distance(destination.Location)
	duration := time.Duration(distance*1.4) * time.Second
	return &WalkingLeg{Distance: distance, Duration: duration}, nil
}

func (p *FixedPlanner) Explore() {

}

func (p *FixedPlanner) ExploreReversed(t time.Time, fixed *FixedLeg) (*Leg, error) {
	origin, destination, err := p.findStops(fixed)
	if err != nil {
		return nil, err
	}

	if fixed.Walk {
		distance := origin.Location.Distance(destination.Location)
		return &Leg{
			Arrival: t,
			WalkingLeg: &WalkingLeg{
				Distance: distance,
			},
		}, nil
	}

	destinationStopTime, err := p.ScheduleIndex.Get(destination.Id, fixed.RouteId).Previous(t)
	if err != nil {
		return nil, err
	}

	all, _ := p.StopTimesFromTrip.Get(destinationStopTime.TripId)
	originStopTime, err := p.findStopTime(fixed.Origin, all)
	if err != nil {
		return nil, err
	}

	waitDuration := stopTimeDiffDuration(destinationStopTime.Time, t)
	transitDuration := stopTimeDiffDuration(originStopTime.Time, destinationStopTime.Time)

	return &Leg{
		Arrival:  t,
		FixedLeg: fixed,
		TransitLeg: &TransitLeg{
			TripId:                destinationStopTime.TripId,
			OriginStopTimeId:      originStopTime.ID(),
			DestinationStopTimeId: destinationStopTime.ID(),
			WaitDuration:          waitDuration,
			TransitDuration:       transitDuration,
		},
	}, nil

}

func (p *FixedPlanner) findStops(fixed *FixedLeg) (origin, destination model.Stop, err error) {
	empty := model.Stop{}
	origin, ok := p.StopIndex.Get(fixed.Origin)
	if !ok {
		return empty, empty, fmt.Errorf("origin stop %s not found", fixed.Origin)
	}

	destination, ok = p.StopIndex.Get(fixed.Destination)
	if !ok {
		return empty, empty, fmt.Errorf("destination stop %s not found", fixed.Destination)
	}
	return origin, destination, nil
}

func (p *FixedPlanner) findStopTime(stopId string, all []model.StopTime) (model.StopTime, error) {
	for _, stopTime := range all {
		if stopTime.StopId == stopId {
			return stopTime, nil
		}
	}
	return model.StopTime{}, errors.New("not found")
}
