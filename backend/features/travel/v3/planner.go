package v3

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"stop-checker.com/db/model"
	"stop-checker.com/db/repository"
	"stop-checker.com/features/travel/algorithms"
)

type fastestTransit struct {
	tripId       string
	routeId      string
	stopId       string
	stopArrival  time.Time
	stopLocation model.Location
}

func (f *fastestTransit) Faster(t time.Time, mode Mode) bool {
	if mode == DEPART_AT {
		return f.stopArrival.Before(t)
	}
	return f.stopArrival.After(t)
}

type Planner struct {
	stopLocationIndex repository.StopLocationSearch
	stopRouteIndex    repository.StopRoutes
	reachIndex        repository.Reach
	directionsCache   WalkingDirectionsCache
	directions        WalkingDirections
}

func NewPlanner(
	stopLocationIndex repository.StopLocationSearch,
	stopRouteIndex repository.StopRoutes,
	reachIndex repository.Reach,
	directionsCache WalkingDirectionsCache,
	directions WalkingDirections,
) *Planner {
	return &Planner{
		stopLocationIndex: stopLocationIndex,
		stopRouteIndex:    stopRouteIndex,
		reachIndex:        reachIndex,
		directionsCache:   directionsCache,
		directions:        directions,
	}
}

func (p *Planner) Depart(at time.Time, origin, destination model.Location) {
	solution, err := p.explore(at, origin, destination, DEPART_AT)
	if err != nil {
		panic(err)
	}

	for solution != nil {
		fmt.Println(solution.id, solution.time)
		fmt.Println(solution.transit, solution.walking)
		solution = solution.prev
	}
}

func (p *Planner) explore(t time.Time, initial, target model.Location, mode Mode) (*node, error) {
	// initial node
	initialNode := createInitialNode(t, initial)

	// priority queue
	pq := algorithms.NewPriorityQueue(func(a, b *node) bool {
		return a.Weight(target, t, mode) < b.Weight(target, t, mode)
	})
	pq.Push(p.exploreInitial(initialNode, mode)...)

	// visited
	explored := algorithms.Set{}

	for !pq.Empty() {
		current := pq.Pop()

		// duplicate node
		if explored.Contains(current.ID()) {
			continue
		}
		explored.Add(current.ID())
		fmt.Printf("%f,%f,blue,square\n", current.location.Latitude, current.location.Longitude)

		// target node
		if current.kind == TARGET {
			fmt.Println(explored.Size())
			return current, nil
		}

		// explore nodes by walking
		pq.Push(p.exploreWalking(current, mode)...)

		// explore nodes by transit
		pq.Push(p.exploreTransit(current, mode)...)

		distance := current.location.Distance(target)
		duration := walkingDuration(distance)
		if mode == ARRIVE_BY {
			duration = -duration
		}

		if distance < MAX_WALK_TARGET {
			pq.Push(createTargetNode(current, &targetNodeParams{
				location: target,
				arrival:  current.time.Add(duration),
				distance: distance,
			}))
		}
	}

	return nil, errors.New("no solution")
}

func (p *Planner) exploreWalking(current *node, mode Mode) []*node {
	nodes := []*node{}

	// skip don't walk two nodes in a row
	if current.kind == STOP && current.transit == nil {
		return nodes
	}

	for _, neighbor := range p.stopLocationIndex.Query(current.location, MAX_WALK) {
		// don't walk to the same stop
		if neighbor.ID() == current.ID() {
			continue
		}

		directions := p.getWalkingDirections(current, neighbor)

		// calculate arrival time
		var arrival time.Time
		if mode == DEPART_AT {
			arrival = current.time.Add(walkingDuration(directions.Distance))
		} else {
			arrival = current.time.Add(-walkingDuration(directions.Distance))
		}

		nodes = append(nodes, createWalkingNode(current, &walkingNodeParams{
			id:       neighbor.ID(),
			location: neighbor.Location,
			arrival:  arrival,
			distance: directions.Distance,
		}))
	}

	return nodes
}

func (p *Planner) getWalkingDirections(current *node, neighbor model.StopWithDistance) model.Path {
	directions, err := p.directionsCache.GetDirections(current.ID(), neighbor.ID())

	if err != nil {
		// should not happen if MAX_WALK and the cache are configured properly
		log.Warn().
			Float64("distance", neighbor.Distance).
			Str("current-id", current.ID()).
			Str("neighbor-id", neighbor.ID()).
			Msg("failed to get walking directions from cache")

		return model.Path{
			Distance: neighbor.Distance,
			Path:     []model.Location{current.location, neighbor.Location},
		}
	}

	return directions
}

func (p *Planner) exploreTransit(current *node, mode Mode) []*node {
	blockers := algorithms.Set{}
	fastest := map[string]fastestTransit{} // fastest transit {stopid: fastest}

	for _, stopRoute := range p.stopRouteIndex.Get(current.ID()) {
		// ignore blocked stop routes by current node
		if current.Blocked(stopRoute.DirectedID()) {
			continue
		}
		blockers.Add(stopRoute.DirectedID())

		// reachable stops
		for _, reachable := range p.exploreTransitRoute(current, stopRoute, mode) {
			current, seen := fastest[reachable.stopId]
			if !seen || !current.Faster(reachable.stopArrival, mode) {
				fastest[reachable.stopId] = reachable
			}
		}
	}

	return p.getTransitNodes(current, fastest, blockers)
}

func (p *Planner) exploreTransitRoute(current *node, stopRoute model.StopRoute, mode Mode) []fastestTransit {
	var results []model.ReachableSchedule
	if mode == DEPART_AT {
		results = p.reachIndex.ReachableForwardWithNext(current.ID(), stopRoute.RouteId, current.time)
	} else {
		results = p.reachIndex.ReachableBackwardWithPrevious(current.ID(), stopRoute.RouteId, current.time)
	}

	reachable := make([]fastestTransit, len(results))

	for i, result := range results {
		if mode == DEPART_AT {
			/* since we're going forward (depart at):
			TODO: comments
			*/
			reachable[i] = fastestTransit{
				stopId:       result.Destination.Id,
				stopArrival:  result.Arrival,
				stopLocation: result.Destination.Location,
				tripId:       result.Trip.Id,
				routeId:      result.Trip.RouteId,
			}
		} else {
			/* since we're going backward (arrive by):
			TODO: comments
			*/
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

func (p *Planner) getTransitNodes(current *node, fastest map[string]fastestTransit, blockers algorithms.Set) []*node {
	nodes := []*node{}
	for _, f := range fastest {
		nodes = append(nodes, createTransitNode(current, &transitNodeParams{
			id:       f.stopId,
			location: f.stopLocation,
			arrival:  f.stopArrival,
			transit: &transit{
				tripId:  f.tripId,
				routeId: f.routeId,
			},
			blockers: blockers,
		}))
	}
	return nodes
}

func (p *Planner) exploreInitial(initial *node, mode Mode) []*node {
	t0 := time.Now()
	neighbors := p.stopLocationIndex.Query(initial.location, MAX_WALK_INITIAL)
	nodes := make([]*node, len(neighbors))
	wg := sync.WaitGroup{}

	for i, neighbor := range neighbors {
		wg.Add(1)

		go func(i int, neighbor model.StopWithDistance) {
			// directions
			directions, err := p.directions.GetDirections(initial.location, neighbor.Location)
			if err != nil {
				fmt.Println(err)
			}

			// calculate arrival time
			var arrival time.Time
			if mode == DEPART_AT {
				arrival = initial.time.Add(walkingDuration(directions.Distance))
			} else {
				arrival = initial.time.Add(-walkingDuration(directions.Distance))
			}

			nodes[i] = createWalkingNode(initial, &walkingNodeParams{
				id:       neighbor.ID(),
				location: neighbor.Location,
				arrival:  arrival,
				distance: directions.Distance,
			})

			wg.Done()
		}(i, neighbor)

	}

	wg.Wait()

	log.Info().Dur("duration", time.Since(t0)).Int("neighbors", len(neighbors)).Msg("initial node")

	return nodes
}
