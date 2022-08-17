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
}

type node struct {
	walk      bool         // walked to the stop
	routeId   string       // route taken to the stop, empty when walking
	stopId    string       // stop id
	transfers int          // transfers
	arrival   time.Time    // arrival time at this stop
	blockers  dijkstra.Set // routes than cannot be taken from this node
}

func (n *node) ID() string {
	return n.stopId
}

func (n *node) Weight() int {
	return int(n.arrival.Unix()) + n.transfers*60*5 // 5 minute penalty per transfer
}

func (n *node) Arrival() time.Time {
	return n.arrival
}

func (n *node) Blocked(routeId string) bool {
	return n.blockers.Contains(routeId)
}
