package travel

import (
	"time"

	"stop-checker.com/db"
	"stop-checker.com/db/model"
	"stop-checker.com/features/travel/dijkstra"
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
		blockers:  dijkstra.Set{},
	}

	solution, err := dijkstra.Algorithm(&dijkstra.Config[*node]{
		Destination: destination,
		Initial:     initial,
		Expand:      p.expand,
		Compare: func(a, b *node) bool {
			return a.arrival.Before(b.arrival)
		},
	})

	if err != nil {
		return nil, err
	}

	return p.route(solution), nil
}

func (p *Planner) route(solution *dijkstra.Path[*node]) Route {
	route := Route{}

	for solution.Prev != nil {
		route = append(route, &FixedLeg{
			Origin:      solution.Prev.ID(),
			Destination: solution.ID(),
			RouteId:     solution.Node.routeId,
			Walk:        solution.Node.walk,
		})

		solution = solution.Prev
	}

	// reverse the route
	for i, j := 0, len(route)-1; i < j; i, j = i+1, j-1 {
		route[i], route[j] = route[j], route[i]
	}

	return route
}

func (p *Planner) expand(n *node) []*node {
	transit := p.expandTransit(n)
	walking := p.expandWalk(n)
	return append(transit, walking...)
}

func (p *Planner) expandWalk(origin *node) []*node {
	stop, _ := p.stopIndex.Get(origin.ID())

	originRoutes := dijkstra.Set{}
	originRoutesList := p.stopRouteIndex.Get(origin.ID())

	for _, originRoute := range originRoutesList {
		originRoutes.Add(originRoute.DirectedID())
	}

	// closest walk for each route key:routeid
	closest := map[string]closestWalk{}

	// for all stops within a 200m radius
	for _, neighbor := range p.stopLocationIndex.Query(stop.Location, 200) {
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
		connections = append(connections, &node{
			walk:      true,
			routeId:   "",
			stopId:    c.stopId,
			transfers: origin.transfers,
			arrival:   origin.arrival.Add(c.duration),
			blockers:  origin.blockers,
		})
	}

	return connections
}

func (p *Planner) expandTransit(n *node) []*node {
	// origin
	origin := n.ID()
	originArrival := n.Arrival()

	blockers := dijkstra.Set{}             // blocked routes key:stopid, set of directed routeid
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
			current, ok := fastest[result.Destination.Id]

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
