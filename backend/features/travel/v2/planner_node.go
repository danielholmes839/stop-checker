package v2

import (
	"fmt"
	"time"

	"stop-checker.com/db/model"
	"stop-checker.com/features/travel/algorithms"
)

type transit struct {
	route string
	trip  string
}

type node struct {
	previous  *node
	id        string
	location  model.Location
	time      time.Time
	kind      kind
	blockers  algorithms.Set
	transfers int
	walkTotal float64

	// store the weight
	computeWeight bool
	weight        time.Duration

	// how this node was reached
	walk    bool
	transit *transit
}

func (n *node) ID() string {
	return n.id
}

func (n *node) Location() model.Location {
	return n.location
}

func (n *node) Time() time.Time {
	return n.time
}

func (n *node) String() string {
	if n.transit != nil {
		return fmt.Sprintf("node{id: %s, transfers: %d, route: %s, time: %s}", n.ID(), n.transfers, n.transit.route, n.Time())
	}
	return fmt.Sprintf("node{id: %s, transfers: %d, time: %s}", n.ID(), n.transfers, n.Time())
}

func (n *node) Weight(target model.Location, initial time.Time, forward bool) time.Duration {
	if n.computeWeight {
		var duration time.Duration
		if forward {
			duration = n.Time().Sub(initial)
		} else {
			duration = initial.Sub(n.Time())
		}

		// calculating heuristics
		distancePenalty := time.Duration(n.location.Distance(target)) * DISTANCE_PENALTY
		walkPenalty := walkingDuration(n.walkTotal * WALK_PENALTY)
		transferPenalty := time.Duration(n.transfers) * TRANSFER_PENALTY
		penalty := distancePenalty + walkPenalty + transferPenalty

		// update the node
		n.weight = duration + penalty
		n.computeWeight = false
	}

	return n.weight
}

func (n *node) blocked(directedRouteId string) bool {
	return n.blockers.Contains(directedRouteId)
}

func createInitialNode(initial model.Location, t time.Time) *node {
	return &node{
		previous:      nil,
		id:            "INITIAL",
		transfers:     0,
		walkTotal:     0,
		location:      initial,
		time:          t,
		kind:          INITIAL,
		computeWeight: true,
		weight:        0,
		walk:          false,
		transit:       nil,
		blockers:      algorithms.Set{},
	}
}

func createTargetNode(previous *node, target model.Location, forward bool) *node {
	distance := target.Distance(previous.Location())
	duration := walkingDuration(distance)
	if !forward {
		duration = -duration
	}

	return &node{
		previous:      previous,
		id:            "TARGET",
		transfers:     previous.transfers,
		walkTotal:     previous.walkTotal + distance,
		location:      target,
		kind:          TARGET,
		computeWeight: true,
		weight:        0,
		blockers:      nil, // never used
		time:          previous.Time().Add(duration),
		walk:          true, // always walk to destination
		transit:       nil,
	}
}

func createWalkingNode(previous *node, walk closestWalk, forward bool) *node {
	t := previous.Time()
	if forward {
		t = t.Add(walk.duration)
	} else {
		t = t.Add(-walk.duration)
	}

	return &node{
		previous:      previous,
		id:            walk.stopId,
		transfers:     previous.transfers,
		walkTotal:     previous.walkTotal + walk.distance,
		location:      walk.stopLocation,
		kind:          STOP,
		computeWeight: true,
		weight:        0,
		blockers:      previous.blockers,
		time:          t,
		walk:          true,
		transit:       nil,
	}
}

func createTransitNode(previous *node, ft fastestTransit, blockers algorithms.Set) *node {
	return &node{
		previous:      previous,
		id:            ft.stopId,
		transfers:     previous.transfers + 1,
		walkTotal:     previous.walkTotal,
		location:      ft.stopLocation,
		kind:          STOP,
		computeWeight: true,
		weight:        0,
		blockers:      blockers,
		time:          ft.stopArrival,
		walk:          false,
		transit: &transit{
			route: ft.routeId,
			trip:  ft.tripId,
		},
	}
}
