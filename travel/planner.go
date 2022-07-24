package travel

import (
	"fmt"
	"time"

	"stop-checker.com/db"
	"stop-checker.com/db/model"
)

type FastestTrip struct {
	Arrival         time.Time
	TransitDuration time.Duration
	WaitDuration    time.Duration
	TripId          string
}

func (trip FastestTrip) Duration() time.Duration {
	return trip.WaitDuration + trip.WaitDuration
}

type Set map[string]struct{}

func (s Set) Add(key string) {
	s[key] = struct{}{}
}

func (s Set) Contains(key string) bool {
	_, ok := s[key]
	return ok
}

type Edge struct {
	Node            *Node
	TripId          string
	TransitDuration time.Duration
	WaitDuration    time.Duration
}

type Node struct {
	Previous *Edge
	StopId   string
	Arrival  time.Time // time we arrived at this node (includes a date)
	Blocked  Set       // blocked route ids
}

type Planner struct {
	ScheduleIndex     *db.ScheduleIndex
	StopRouteIndex    *db.StopRouteIndex
	StopTimesFromTrip *db.InvertedIndex[model.StopTime]
}

func (p *Planner) Depart(departure time.Time, origin string, destination string) {
	initial := &Node{
		StopId:  origin,
		Arrival: departure,
		Blocked: Set{},
	}

	next := p.Explore(initial)
	fmt.Println(next)
}

// Explore nodes
func (p *Planner) Explore(n *Node) []*Node {
	blockers := map[string]Set{}			// stopId:set
	fastest := map[string]FastestTrip{}		// stopId:FastestTrip

	for _, route := range p.StopRouteIndex.Get(n.StopId) {
		// skip blocked routes
		if n.Blocked.Contains(route.Id) {
			continue
		}

		// route next stop time
		next, err := p.ScheduleIndex.Get(n.StopId, route.Id).Next(n.Arrival)
		if err != nil {
			continue
		}

		// route wait duration
		waitDuration := stopTimeDiffDuration(n.Arrival, next.Time)

		for _, visit := range p.visits(next) {
			current, seen := fastest[visit.StopId]
			transitDuration := stopTimeDiffDuration(next.Time, visit.Time)
			totalDuration := waitDuration + transitDuration

			// update the fastest trip
			if !seen || totalDuration < current.Duration() {
				fastest[visit.StopId] = FastestTrip{
					Arrival:         n.Arrival.Add(totalDuration),
					TransitDuration: transitDuration,
					WaitDuration:    waitDuration,
					TripId:          next.TripId,
				}
			}

			// update the stop blockers
			if _, ok := blockers[visit.StopId]; !ok {
				blockers[visit.StopId] = Set{}
			}
			stopBlockers := blockers[visit.StopId]
			stopBlockers.Add(route.Id)

		}
	}

	// create the nodes
	connections := []*Node{}

	for stopId, trip := range fastest {
		connections = append(connections, &Node{
			Previous: &Edge{
				Node:            n,
				TripId:          trip.TripId,
				TransitDuration: trip.TransitDuration,
				WaitDuration:    trip.WaitDuration,
			},
			Blocked: blockers[stopId],
			StopId:  stopId,
			Arrival: trip.Arrival,
		})
	}

	return connections
}

func (p *Planner) visits(origin model.StopTime) []model.StopTime {
	all, _ := p.StopTimesFromTrip.Get(origin.TripId)
	connections := []model.StopTime{}

	for _, stopTime := range all {
		if stopTime.StopSeq > origin.StopSeq {
			connections = append(connections, stopTime)
		}
	}

	return connections
}

func stopTimeDiffDuration(from, to time.Time) time.Duration {
	f := from.Hour()*60 + from.Minute()
	t := to.Hour()*60 + to.Minute()

	if t < f {
		t += 60 * 24
	}

	return time.Duration(t-f) * time.Minute
}
