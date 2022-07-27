package travel

import (
	"time"

	"stop-checker.com/travel/dijkstra"
)

type closestWalk struct {
	distance float64
	stopId   string
	duration time.Duration
}

type fastestTransit struct {
	routeId string
	arrival time.Time
	wait    time.Duration
	transit time.Duration
}

type node struct {
	walk     bool
	routeId  string
	stopId   string
	arrival  time.Time
	duration time.Duration
	blockers dijkstra.Set
}

func (n *node) ID() string {
	return n.stopId
}

func (n *node) Weight() int {
	return int(n.arrival.Unix())
}

func (n *node) Arrival() time.Time {
	return n.arrival
}

func (n *node) Blocked(routeId string) bool {
	return n.blockers.Contains(routeId)
}
