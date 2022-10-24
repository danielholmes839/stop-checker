package v2

import (
	"errors"
	"time"

	"stop-checker.com/db"
	"stop-checker.com/db/model"
	"stop-checker.com/features/travel/algorithms"
)

type closestWalk struct {
	distance     float64
	stopId       string
	stopLocation model.Location
	duration     time.Duration
}

type fastestTransit struct {
	tripId       string
	routeId      string
	stopId       string
	stopArrival  time.Time
	stopLocation model.Location
}

func (f *fastestTransit) Faster(t time.Time, forward bool) bool {
	if forward {
		return f.stopArrival.Before(t)
	}
	return f.stopArrival.After(t)
}

type Planner struct {
	stopLocationIndex *db.StopLocationIndex
	stopRouteIndex    *db.StopRouteIndex
	stopIndex         *db.Index[model.Stop]
	reachIndex        *db.ReachIndex
}

func NewPlanner(database *db.Database) *Planner {
	return &Planner{
		stopLocationIndex: database.StopLocationIndex,
		stopRouteIndex:    database.StopRouteIndex,
		stopIndex:         database.Stops,
		reachIndex:        database.ReachIndex,
	}
}

func (p *Planner) Depart(origin, destination model.Location, at time.Time) (*node, error) {
	return p.plan(origin, destination, at, true)
}

func (p *Planner) Arrive(origin, destination model.Location, by time.Time) (*node, error) {
	return p.plan(origin, destination, by, false)
}

func (p *Planner) plan(origin, destination model.Location, t time.Time, forward bool) (*node, error) {
	// target and initial locations
	target, initial := p.getTargetAndInitial(origin, destination, forward)

	explored := algorithms.Set{}
	pq := algorithms.NewPriorityQueue(func(a, b *node) bool {
		return a.Weight(target, t, forward) < b.Weight(target, t, forward)
	})
	pq.Push(createInitialNode(initial, t))

	for !pq.Empty() {
		current := pq.Pop()

		// check if we reached a target node
		if current.kind == TARGET {
			return current, nil
		}

		// check if we reached a node thats already explored
		if explored.Contains(current.ID()) {
			continue
		}
		explored.Add(current.ID())

		// We're exploring a potential last stop so add another node to walk to the final destination.
		if target.Distance(current.Location()) < MAX_WALK {
			pq.Push(createTargetNode(current, target, forward))
		}

		// Explore by walking
		for _, node := range p.exploreWalk(current, forward) {
			pq.Push(node)
		}

		// Explore using transit
		for _, node := range p.exploreTransit(current, forward) {
			pq.Push(node)
		}
	}

	return nil, errors.New("failed to create travel route")
}

func (p *Planner) exploreTransit(current *node, forward bool) []*node {
	if current.kind != STOP {
		return []*node{}
	}

	blockers := algorithms.Set{}           // blocked set of directed route ids
	fastest := map[string]fastestTransit{} // fastest transit option key:stopid

	for _, stopRoute := range p.stopRouteIndex.Get(current.ID()) {
		if current.blocked(stopRoute.DirectedID()) {
			continue
		}
		blockers.Add(stopRoute.DirectedID())

		// reachable stops
		for _, reachable := range p.getReachable(current, stopRoute, forward) {
			current, seen := fastest[reachable.stopId]
			if !seen || !current.Faster(reachable.stopArrival, forward) {
				fastest[reachable.stopId] = reachable
			}
		}
	}

	return p.getTransitNodes(current, fastest, blockers)
}

func (p *Planner) exploreWalk(current *node, forward bool) []*node {
	if current.walk {
		return []*node{}
	}

	closest := map[string]closestWalk{}

	for _, neighbor := range p.getNeighbors(current) {
		if neighbor.ID() == current.ID() {
			continue
		}

		for _, route := range p.stopRouteIndex.Get(neighbor.ID()) {
			directedRouteId := route.DirectedID()
			if current.blocked(directedRouteId) {
				continue
			}

			walk, seen := closest[directedRouteId]
			if !seen || walk.distance > neighbor.Distance {
				closest[directedRouteId] = closestWalk{
					distance:     neighbor.Distance,
					stopId:       neighbor.ID(),
					stopLocation: neighbor.Location,
					duration:     walkingDuration(neighbor.Distance),
				}
			}
		}
	}

	return p.getWalkNodes(current, closest, forward)
}

func (p *Planner) getReachable(current *node, stopRoute model.StopRoute, forward bool) []fastestTransit {
	var results []db.ReachableSchedule
	if forward {
		results = p.reachIndex.ReachableForwardWithNext(current.ID(), stopRoute.RouteId, current.Time())
	} else {
		results = p.reachIndex.ReachableBackwardWithPrevious(current.ID(), stopRoute.RouteId, current.Time())
	}

	reachable := make([]fastestTransit, len(results))

	for i, result := range results {
		if forward {
			reachable[i] = fastestTransit{
				stopId:       result.Destination.Id,
				stopArrival:  result.Arrival,
				stopLocation: result.Destination.Location,
				tripId:       result.Trip.Id,
				routeId:      result.Trip.RouteId,
			}
		} else {
			reachable[i] = fastestTransit{
				stopId:       result.Origin.Id,
				stopArrival:  result.Departure,
				stopLocation: result.Origin.Location,
				tripId:       result.Trip.Id,
				routeId:      result.Trip.RouteId,
			}
		}
	}

	return reachable
}

func (p *Planner) getNeighbors(current *node) []db.StopLocationResult {
	var neighbors []db.StopLocationResult
	if current.kind == INITIAL {
		neighbors = p.stopLocationIndex.Query(current.Location(), MAX_WALK)
	} else {
		neighbors = p.stopLocationIndex.Query(current.Location(), MAX_WALK_EXPLORE)
	}
	return neighbors
}

func (p *Planner) getTransitNodes(current *node, fastest map[string]fastestTransit, blockers algorithms.Set) []*node {
	nodes := []*node{}
	for _, transit := range fastest {
		nodes = append(nodes, createTransitNode(current, transit, blockers))
	}
	return nodes
}

func (p *Planner) getWalkNodes(current *node, closest map[string]closestWalk, forward bool) []*node {
	nodes := []*node{}
	added := algorithms.Set{}

	for _, walk := range closest {
		// don't add duplicates
		if added.Contains(walk.stopId) {
			continue
		}
		added.Add(walk.stopId)

		nodes = append(nodes, createWalkingNode(current, walk, forward))
	}
	return nodes
}

func (p *Planner) getTargetAndInitial(origin, destination model.Location, forward bool) (initial, target model.Location) {
	if forward {
		initial = origin
		target = destination
	} else {
		initial = destination
		target = origin
	}
	return initial, target
}
