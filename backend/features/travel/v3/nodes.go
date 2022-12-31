package v3

import (
	"time"

	"stop-checker.com/db/model"
	"stop-checker.com/features/travel/algorithms"
)

type transit struct {
	tripId  string
	routeId string
}

type node struct {
	prev     *node
	id       string
	kind     Kind
	time     time.Time
	blockers algorithms.Set
	transit  *transit

	// heuristics
	location  model.Location
	transfers int     // the cumulative number of buses taken
	walking   float64 // the cumulative walking distance

	// weight
	computeWeight bool
	weight        time.Duration
}

func (n *node) ID() string {
	return n.id
}

func (n *node) Transfers() int {
	/*
		the first transfer is "free".
		taking 1 bus doesn't count as a transfer
		taking 2 buses counts as 1 transfer
	*/
	if n.transfers > 0 {
		return n.transfers - 1
	}
	return 0
}

func (n *node) Blocked(directedRouteId string) bool {
	return n.blockers.Contains(directedRouteId)
}

func (n *node) Weight(target model.Location, initial time.Time, mode Mode) time.Duration {
	if n.computeWeight {
		var duration time.Duration
		if mode == DEPART_AT {
			duration = n.time.Sub(initial)
		} else {
			duration = initial.Sub(n.time)
		}

		// calculating heuristics
		distancePenalty := time.Duration(n.location.Distance(target)) * DISTANCE_PENALTY
		walkPenalty := walkingDuration(n.walking * WALK_PENALTY)
		transferPenalty := time.Duration(n.Transfers()) * TRANSFER_PENALTY
		penalty := distancePenalty + walkPenalty + transferPenalty

		// update the node
		n.weight = duration + penalty
		n.computeWeight = false
	}

	return n.weight
}

type transitNodeParams struct {
	id       string
	location model.Location
	arrival  time.Time // the time we arrive at this node
	transit  *transit
	blockers algorithms.Set
}

type walkingNodeParams struct {
	id       string
	location model.Location
	arrival  time.Time // the time we arrive at this node
	distance float64
}

type initialNodeParams struct {
	location model.Location
	start    time.Time // the arrive by or depart at time
}

// walkingNodeParams without an id
type targetNodeParams struct {
	location model.Location
	arrival  time.Time // the time we arrive at this node
	distance float64
}

func createInitialNode(t time.Time, initial model.Location) *node {
	return &node{
		prev:          nil,
		id:            "INITIAL",
		kind:          INITIAL,
		time:          t,
		location:      initial,
		blockers:      algorithms.Set{},
		transit:       nil, // not necessary since we will not explore transit from this node
		transfers:     0,
		walking:       0,
		weight:        0,
		computeWeight: true,
	}
}

func createTargetNode(prev *node, params *targetNodeParams) *node {
	return &node{
		prev:          prev,
		id:            "TARGET",
		kind:          TARGET,
		time:          params.arrival,
		blockers:      algorithms.Set{}, // not necessary since we will not explore transit from this node
		transit:       nil,
		location:      params.location,
		transfers:     prev.transfers,
		walking:       prev.walking + params.distance, // increase cumulative walking distance
		weight:        0,                              // updated separately
		computeWeight: true,
	}
}

func createWalkingNode(prev *node, params *walkingNodeParams) *node {
	return &node{
		prev:          prev,
		id:            params.id,
		kind:          STOP,
		time:          params.arrival,
		blockers:      prev.blockers, // ignore the same routes as the previous node.
		transit:       nil,
		location:      params.location,
		transfers:     prev.transfers,
		walking:       prev.walking + params.distance, // increase cumulative walking distance
		weight:        0,                              // updated separately
		computeWeight: true,
	}
}

func createTransitNode(prev *node, params *transitNodeParams) *node {
	return &node{
		prev:          prev,
		id:            params.id,
		kind:          STOP,
		time:          params.arrival,
		blockers:      params.blockers,
		transit:       params.transit,
		location:      params.location,
		transfers:     prev.transfers + 1, // increase the number of transfers
		walking:       prev.walking,
		weight:        0, // updated separately
		computeWeight: true,
	}
}
