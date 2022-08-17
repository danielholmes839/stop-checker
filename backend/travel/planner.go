package travel

import (
	"math"
	"time"

	"stop-checker.com/db"
	"stop-checker.com/db/model"
	"stop-checker.com/travel/dijkstra"
)

type PlannerConfig struct {
	ScheduleIndex     *db.ScheduleIndex
	StopLocationIndex *db.StopLocationIndex
	StopRouteIndex    *db.StopRouteIndex
	StopIndex         *db.Index[model.Stop]
	StopTimesFromTrip *db.InvertedIndex[model.StopTime]
}

type Planner struct {
	scheduleIndex     *db.ScheduleIndex
	stopLocationIndex *db.StopLocationIndex
	stopRouteIndex    *db.StopRouteIndex
	stopIndex         *db.Index[model.Stop]
	stopTimesFromTrip *db.InvertedIndex[model.StopTime]
}

func NewPlanner(config *PlannerConfig) *Planner {
	return &Planner{
		scheduleIndex:     config.ScheduleIndex,
		stopLocationIndex: config.StopLocationIndex,
		stopRouteIndex:    config.StopRouteIndex,
		stopIndex:         config.StopIndex,
		stopTimesFromTrip: config.StopTimesFromTrip,
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

	// closest walk for each route
	closest := map[string]closestWalk{}

	// for all stops within a 250m radius
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

	// calculate walking distance and duration

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

	blockers := map[string]dijkstra.Set{}  // blocked routes key:stopid, set of routeid
	fastest := map[string]fastestTransit{} // fastest transit option key:stopid

	// expand on routes
	for _, route := range p.stopRouteIndex.Get(n.ID()) {
		if n.Blocked(route.DirectedID()) {
			continue
		}

		// lookup the trip and "tripOrigin" stop time
		tripOrigin, _ := p.scheduleIndex.Get(origin, route.RouteId).Next(originArrival)

		// calculate the time spent waiting for the trip
		waitDuration := stopTimeDiffDuration(originArrival, tripOrigin.Time)

		for _, tripDestination := range p.expandTrip(tripOrigin) {
			// calculate the time spent in transit and the destination arrival time
			transitDuration := stopTimeDiffDuration(tripOrigin.Time, tripDestination.Time)
			tripDestinationArrival := n.arrival.Add(waitDuration + transitDuration)

			// the current fastest trip
			current, seen := fastest[tripDestination.StopId]

			// the current fastest trip should be replaced
			if !seen || tripDestinationArrival.Before(current.arrival) {
				fastest[tripDestination.StopId] = fastestTransit{
					arrival: tripDestinationArrival,
					routeId: route.RouteId,
				}
			}

			// add key to blocked route
			if _, ok := blockers[tripDestination.StopId]; !ok {
				blockers[tripDestination.StopId] = dijkstra.Set{}
			}

			stopBlockers := blockers[tripDestination.StopId]
			stopBlockers.Add(route.DirectedID())
		}
	}

	// build the connections
	connections := []*node{}

	for stopId, trip := range fastest {
		connection := &node{
			stopId:    stopId,
			arrival:   trip.arrival,
			transfers: n.transfers + 1,
			blockers:  blockers[stopId],
			routeId:   trip.routeId,
			walk:      false,
		}

		connections = append(connections, connection)
	}

	return connections
}

func (p *Planner) expandTrip(origin model.StopTime) []model.StopTime {
	// all stop times after the origin stop time
	all, _ := p.stopTimesFromTrip.Get(origin.TripId)
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

func walkingDuration(distance float64) time.Duration {
	duration := time.Duration(math.Round(distance*1.4/60)) * time.Minute
	if duration < time.Minute {
		return time.Minute
	}
	return duration
}
