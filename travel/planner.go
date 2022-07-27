package travel

import (
	"math"
	"time"

	"stop-checker.com/db"
	"stop-checker.com/db/model"
	"stop-checker.com/travel/dijkstra"
)

type Planner struct {
	ScheduleIndex     *db.ScheduleIndex
	StopLocationIndex *db.StopLocationIndex
	StopRouteIndex    *db.StopRouteIndex
	StopIndex         *db.Index[model.Stop]
	TripIndex         *db.Index[model.Trip]
	StopTimesFromTrip *db.InvertedIndex[model.StopTime]
}

func (p *Planner) Depart(at time.Time, origin, destination string) (Route, error) {
	initial := &node{
		stopId:   origin,
		arrival:  at,
		duration: time.Duration(0),
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

	// reverse
	for i, j := 0, len(route)-1; i < j; i, j = i+1, j-1 {
		route[i], route[j] = route[j], route[i]
	}

	return route
}

func (p *Planner) expand(n *node) []*node {
	transit, _ := p.expandTransit(n)
	walking := p.expandWalk(n)
	return append(transit, walking...)
}

func (p *Planner) expandWalk(origin *node) []*node {
	stop, _ := p.StopIndex.Get(origin.ID())

	originRoutesIndex := dijkstra.Set{}
	originRoutes := p.StopRouteIndex.Get(origin.ID())

	for _, originRoute := range originRoutes {
		originRoutesIndex.Add(originRoute.ID())
	}

	// closest walk for each route
	closest := map[string]closestWalk{}

	// for all stops within a 250m radius
	for _, neighbor := range p.StopLocationIndex.Query(stop.Location, 200) {
		neighborRoutes := p.StopRouteIndex.Get(neighbor.ID())

		for _, route := range neighborRoutes {
			directedRouteId := route.ID()
			if origin.Blocked(directedRouteId) {
				continue
			}

			if originRoutesIndex.Contains(directedRouteId) {
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
			walk:     true,
			routeId:  "",
			stopId:   c.stopId,
			arrival:  origin.arrival.Add(c.duration),
			blockers: origin.blockers,
		})
	}

	return connections
}

func (p *Planner) expandTransit(n *node) ([]*node, dijkstra.Set) {
	// origin
	origin := n.ID()
	originArrival := n.Arrival()

	stops := dijkstra.Set{n.stopId: struct{}{}}
	blockers := map[string]dijkstra.Set{}
	fastest := map[string]fastestTransit{}

	// expand on routes
	for _, route := range p.StopRouteIndex.Get(n.ID()) {
		if n.Blocked(route.ID()) {
			continue
		}

		// lookup the trip and "tripOrigin" stop time
		tripOrigin, _ := p.ScheduleIndex.Get(origin, route.Route.Id).Next(originArrival)

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
					wait:    waitDuration,
					transit: transitDuration,
					routeId: route.Id,
				}
			}

			// add blocked route
			if _, ok := blockers[tripDestination.StopId]; !ok {
				blockers[tripDestination.StopId] = dijkstra.Set{}
			}

			stopBlockers := blockers[tripDestination.StopId]
			stopBlockers.Add(route.ID())
			stops.Add(tripDestination.StopId)
		}
	}

	// build the connections
	connections := []*node{}

	for stopId, trip := range fastest {
		// stop, _ := p.StopIndex.Get(stopId)
		connection := &node{
			stopId:   stopId,
			arrival:  trip.arrival,
			duration: trip.transit,
			blockers: blockers[stopId],
			routeId:  trip.routeId,
			walk:     false,
		}

		connections = append(connections, connection)
	}

	return connections, stops
}

func (p *Planner) expandTrip(origin model.StopTime) []model.StopTime {
	// all stop times after the origin stop time
	all, _ := p.StopTimesFromTrip.Get(origin.TripId)
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
