package travel

import (
	"time"

	"github.com/rs/zerolog/log"
	"stop-checker.com/db"
	"stop-checker.com/db/model"
	"stop-checker.com/features/travel/algorithms"
)

const MAX_WALK_DISTANCE = 300
const TRANSFER_PENALTY = 5 * time.Minute

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
	walk      bool           // walked to the stop
	routeId   string         // route taken to the stop, empty when walking
	stopId    string         // stop id
	transfers int            // transfers
	arrival   time.Time      // arrival time at this stop
	blockers  algorithms.Set // routes than cannot be taken from this node
}

func (n *node) ID() string {
	return n.stopId
}

func (n *node) Arrival() time.Time {
	return n.arrival
}

func (n *node) Blocked(routeId string) bool {
	return n.blockers.Contains(routeId)
}

type PlannerConfig struct {
	StopLocationIndex *db.StopLocationIndex
	StopRouteIndex    *db.StopRouteIndex
	StopIndex         *db.Index[model.Stop]
	ReachIndex        *db.ReachIndex
}

type Planner struct {
	stopLocationIndex *db.StopLocationIndex
	stopRouteIndex    *db.StopRouteIndex
	stopIndex         *db.Index[model.Stop]
	reachIndex        *db.ReachIndex
}

func NewPlanner(config *PlannerConfig) *Planner {
	return &Planner{
		stopLocationIndex: config.StopLocationIndex,
		stopRouteIndex:    config.StopRouteIndex,
		stopIndex:         config.StopIndex,
		reachIndex:        config.ReachIndex,
	}
}

func (p *Planner) Depart(at time.Time, origin, destination string) (Route, error) {
	initial := &node{
		stopId:    origin,
		arrival:   at,
		walk:      false,
		routeId:   "",
		transfers: 0,
		blockers:  algorithms.Set{},
	}

	t0 := time.Now()

	solution, err := algorithms.Dijkstra(&algorithms.DijkstraConfig[*node]{
		Destination: destination,
		Initial:     []*node{initial},
		Expand:      p.expand,
		Compare: func(a, b *node) bool {
			aPen := time.Duration(a.transfers) * TRANSFER_PENALTY
			bPen := time.Duration(b.transfers) * TRANSFER_PENALTY
			return a.arrival.Add(aPen).Before(b.arrival.Add(bPen))
		},
	})

	if err != nil {
		log.Error().Err(err).
			Str("planner-mode", "depart-at").
			Dur("planner-duration", time.Since(t0)).
			Str("planner-origin", origin).
			Str("planner-destination", destination).
			Msg("failed to create a travel plan")

		return nil, err
	}

	log.Info().
		Str("planner-mode", "depart-at").
		Dur("planner-duration", time.Since(t0)).
		Str("planner-origin", origin).
		Str("planner-destination", destination).
		Msg("successfully created a travel plan")

	return p.route(solution, true), nil
}

func (p *Planner) Arrive(by time.Time, origin, destination string) (Route, error) {
	initial := &node{
		stopId:    destination,
		arrival:   by,
		walk:      false,
		routeId:   "",
		transfers: 0,
		blockers:  algorithms.Set{},
	}

	t0 := time.Now()

	solution, err := algorithms.Dijkstra(&algorithms.DijkstraConfig[*node]{
		Destination: origin, // the target
		Initial:     []*node{initial},
		Expand:      p.expandReverse,
		Compare: func(a, b *node) bool {
			aPen := time.Duration(a.transfers) * -TRANSFER_PENALTY
			bPen := time.Duration(b.transfers) * -TRANSFER_PENALTY
			return a.arrival.Add(aPen).After(b.arrival.Add(bPen))
		},
	})

	if err != nil {
		log.Error().Err(err).
			Str("planner-mode", "arrive-by").
			Dur("planner-duration", time.Since(t0)).
			Str("planner-origin", origin).
			Str("planner-destination", destination).
			Msg("failed to create a travel plan")
		return nil, err
	}

	log.Info().
		Str("planner-mode", "arrive-by").
		Dur("planner-duration", time.Since(t0)).
		Str("planner-origin", origin).
		Str("planner-destination", destination).
		Msg("failed to create a travel plan")

	return p.route(solution, false), nil
}

func (p *Planner) route(solution *algorithms.Path[*node], depart bool) Route {
	route := Route{}

	for solution.Prev != nil {
		if depart {
			route = append(route, &FixedLeg{
				Origin:      solution.Prev.ID(),
				Destination: solution.ID(),
				RouteId:     solution.Node.routeId,
				Walk:        solution.Node.walk,
			})
		} else {
			route = append(route, &FixedLeg{
				Origin:      solution.ID(),
				Destination: solution.Prev.ID(),
				RouteId:     solution.Node.routeId,
				Walk:        solution.Node.walk,
			})
		}

		solution = solution.Prev
	}

	// reverse the route when the solution was found by Depart
	if depart {
		for i, j := 0, len(route)-1; i < j; i, j = i+1, j-1 {
			route[i], route[j] = route[j], route[i]
		}
	}

	return route
}

func (p *Planner) expand(n *node) []*node {
	transit := p.expandTransit(n)
	walking := p.expandWalk(n, false)
	return append(transit, walking...)
}

func (p *Planner) expandReverse(n *node) []*node {
	transit := p.expandTransitReverse(n)
	walking := p.expandWalk(n, true)
	return append(transit, walking...)
}

func (p *Planner) expandWalk(origin *node, reverse bool) []*node {
	if origin.walk {
		return []*node{}
	}

	stop, _ := p.stopIndex.Get(origin.ID())

	originRoutes := algorithms.Set{}
	originRoutesList := p.stopRouteIndex.Get(origin.ID())

	for _, originRoute := range originRoutesList {
		originRoutes.Add(originRoute.DirectedID())
	}

	// closest walk for each route key:routeid
	closest := map[string]closestWalk{}

	// for all stops within a 200m radius
	for _, neighbor := range p.stopLocationIndex.Query(stop.Location, MAX_WALK_DISTANCE) {
		neighborRoutes := p.stopRouteIndex.Get(neighbor.ID())

		for _, route := range neighborRoutes {
			directedRouteId := route.DirectedID()
			if origin.Blocked(directedRouteId) {
				continue
			}

			if originRoutes.Contains(directedRouteId) {
				continue
			}

			distance := neighbor.Location.Distance(stop.Location)
			duration := walkingDuration(distance)

			current, seen := closest[directedRouteId]
			if !seen || current.distance > distance {
				closest[directedRouteId] = closestWalk{
					distance: distance,
					stopId:   neighbor.ID(),
					duration: duration,
				}
			}
		}
	}

	connections := []*node{}

	for _, c := range closest {
		arrival := origin.arrival.Add(c.duration)
		if reverse {
			arrival = origin.arrival.Add(-c.duration)
		}
		connections = append(connections, &node{
			walk:      true,
			routeId:   "",
			stopId:    c.stopId,
			transfers: origin.transfers,
			arrival:   arrival,
			blockers:  origin.blockers,
		})
	}

	return connections
}

func (p *Planner) expandTransit(n *node) []*node {
	// origin
	origin := n.ID()
	originArrival := n.Arrival()

	blockers := algorithms.Set{}           // blocked routes key:stopid, set of directed routeid
	fastest := map[string]fastestTransit{} // fastest transit option key:stopid

	// expand on routes
	for _, stopRoute := range p.stopRouteIndex.Get(n.ID()) {
		if n.Blocked(stopRoute.DirectedID()) {
			continue
		}

		// add stop route to blocked stop routes
		blockers.Add(stopRoute.DirectedID())

		// set fastest transit for each stop
		for _, result := range p.reachIndex.ReachableForwardWithNext(origin, stopRoute.RouteId, originArrival) {
			destinationId := result.Destination.Id
			current, ok := fastest[destinationId]

			if !ok || result.Arrival.Before(current.arrival) {
				fastest[destinationId] = fastestTransit{
					routeId: stopRoute.RouteId,
					arrival: result.Arrival,
				}
			}
		}
	}

	// build the connections
	connections := []*node{}

	for stopId, trip := range fastest {
		connection := &node{
			stopId:    stopId,
			arrival:   trip.arrival,
			transfers: n.transfers + 1,
			blockers:  blockers,
			routeId:   trip.routeId,
			walk:      false,
		}

		connections = append(connections, connection)
	}

	return connections
}

func (p *Planner) expandTransitReverse(n *node) []*node {
	// origin
	destination := n.ID()
	destinationArrival := n.Arrival()

	blockers := algorithms.Set{}           // blocked routes key:stopid, set of directed routeid
	fastest := map[string]fastestTransit{} // fastest transit option key:stopid

	// expand on routes
	for _, stopRoute := range p.stopRouteIndex.Get(n.ID()) {
		// add stop route to blocked stop routes
		blockers.Add(stopRoute.DirectedID())

		if n.Blocked(stopRoute.DirectedID()) {
			continue
		}

		// set fastest transit for each stop
		for _, result := range p.reachIndex.ReachableBackwardWithPrevious(destination, stopRoute.RouteId, destinationArrival) {
			originId := result.Origin.Id
			current, ok := fastest[originId]

			// you can get from this origin to the destination stop faster
			if !ok || result.Departure.After(current.arrival) {
				fastest[originId] = fastestTransit{
					routeId: stopRoute.RouteId,
					arrival: result.Departure,
				}
			}
		}
	}

	// build the connections
	connections := []*node{}

	for stopId, trip := range fastest {
		connection := &node{
			stopId:    stopId,
			arrival:   trip.arrival,
			transfers: n.transfers + 1,
			blockers:  blockers,
			routeId:   trip.routeId,
			walk:      false,
		}

		connections = append(connections, connection)
	}

	return connections
}
