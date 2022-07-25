package travel

import (
	"errors"
	"fmt"
	"strings"
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

func (s Set) String() string {
	keys := []string{}
	for key := range s {
		keys = append(keys, key)
	}
	return fmt.Sprintf("[%s]", strings.Join(keys, ", "))
}

type Edge struct {
	Node            *Node
	Walking         bool // this edge was walked. wait duration is zero.
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

func (n *Node) String() string {
	if n.Previous == nil {
		return fmt.Sprintf("{origin: %s}", n.StopId)
	}

	return fmt.Sprintf("{from:%s} -> {stop:%s, arrival:%s, wait:%s, transit:%s, walk:%t, blocked:%s}", n.Previous.Node.StopId,
		n.StopId, n.Arrival.Format("2006-01-02@15:04"), n.Previous.WaitDuration, n.Previous.TransitDuration, n.Previous.Walking, n.Blocked.String())
}

type Planner struct {
	ScheduleIndex     *db.ScheduleIndex
	StopLocationIndex *db.StopLocationIndex
	StopRouteIndex    *db.StopRouteIndex
	StopIndex         *db.Index[model.Stop]
	StopTimesFromTrip *db.InvertedIndex[model.StopTime]
}

func (p *Planner) Depart(departure time.Time, origin string, destination string) (*Node, error) {
	initial := &Node{
		StopId:  origin,
		Arrival: departure,
		Blocked: Set{},
	}

	return p.Djikstra(initial, destination)
}

// Djikstra's algorithm
func (p *Planner) Djikstra(initial *Node, destination string) (*Node, error) {
	seen := Set{}
	pq := NewNodePriorityQueue()
	pq.Push(initial)

	for !pq.Empty() {
		node := pq.Pop()

		// track seen nodes
		if seen.Contains(node.StopId) {
			continue
		}
		seen.Add(node.StopId)

		// check destination
		if node.StopId == destination {
			return node, nil
		}

		// explore nodes
		neighbors := p.Explore(node)

		for _, neighbor := range neighbors {
			pq.Push(neighbor)
		}
	}

	return nil, errors.New("node not found")
}

// Explore nodes
func (p *Planner) Explore(n *Node) []*Node {
	t := p.visitTransitNeighbors(n)
	transitNodes := p.nodesTransitNeighbors(n, t)

	if len(transitNodes) == 0 {
		return nil
	}

	return append(transitNodes, p.walkingNeighbors(transitNodes, t.stops)...)
}
