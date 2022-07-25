package travel

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"stop-checker.com/db"
	"stop-checker.com/db/model"
)

type FastestTrip struct {
	Arrival         time.Time
	TransitDuration time.Duration
	WaitDuration    time.Duration
	TripId          string
}

func (trip FastestTrip) Duration() time.Duration {
	return trip.WaitDuration + trip.WaitDuration
}

type Set map[string]struct{}

func (s Set) Add(key string) {
	s[key] = struct{}{}
}

func (s Set) Contains(key string) bool {
	_, ok := s[key]
	return ok
}

func (s Set) String() string {
	keys := []string{}
	for key := range s {
		keys = append(keys, key)
	}
	return fmt.Sprintf("[%s]", strings.Join(keys, ", "))
}

type Edge struct {
	Node            *Node
	Walking         bool // this edge was walked. wait duration is zero.
	TripId          string
	TransitDuration time.Duration
	WaitDuration    time.Duration
}

type Node struct {
	Previous *Edge
	StopId   string
	Arrival  time.Time // time we arrived at this node (includes a date)
	Blocked  Set       // blocked route ids
}

func (n *Node) String() string {
	prev := "ORIGIN"

	if n.Previous != nil {
		prev = fmt.Sprintf("{from:%s}", n.Previous.Node.StopId)
	}

	return fmt.Sprintf("%s -> {stop:%s, arrival:%s, wait:%s, transit:%s, walk:%t, blocked:%s}", prev,
		n.StopId, n.Arrival.Format("2006-01-02@15:04"), n.Previous.WaitDuration, n.Previous.TransitDuration, n.Previous.Walking, n.Blocked.String())
}

type Planner struct {
	ScheduleIndex     *db.ScheduleIndex
	StopLocationIndex *db.StopLocationIndex
	StopRouteIndex    *db.StopRouteIndex
	StopIndex         *db.Index[model.Stop]
	StopTimesFromTrip *db.InvertedIndex[model.StopTime]
}

func (p *Planner) Depart(departure time.Time, origin string, destination string) {
	initial := &Node{
		StopId:  origin,
		Arrival: departure,
		Blocked: Set{},
	}

	t1 := time.Now()
	next := p.Explore(initial)
	dur := time.Since(t1)

	sort.Slice(next, func(i, j int) bool {
		return next[i].Arrival.Before(next[j].Arrival)
	})

	for _, node := range next {
		fmt.Println("--------------------")
		stop, _ := p.StopIndex.Get(node.StopId)
		fmt.Println(stop.Name)
		fmt.Println(node)
	}

	fmt.Println(dur)
}

// Explore nodes
func (p *Planner) Explore(n *Node) []*Node {
	t := p.visitTransitNeighbors(n)
	transitNodes := p.nodesTransitNeighbors(n, t)

	if len(transitNodes) == 0 {
		return nil
	}

	return append(transitNodes, p.walkingNeighbors(transitNodes, t.stops)...)
}

type transitNeighbors struct {
	stops    Set
	blockers map[string]Set         // stopId:set{routeId}
	fastest  map[string]FastestTrip // stopId:FastestTrip
}

func (p *Planner) visits(origin model.StopTime) []model.StopTime {
	all, _ := p.StopTimesFromTrip.Get(origin.TripId)
	connections := []model.StopTime{}

	for _, stopTime := range all {
		if stopTime.StopSeq > origin.StopSeq {
			connections = append(connections, stopTime)
		}
	}

	return connections
}

func (p *Planner) visitTransitNeighbors(n *Node) *transitNeighbors {
	stops := Set{}
	blockers := map[string]Set{}
	fastest := map[string]FastestTrip{}

	for _, route := range p.StopRouteIndex.Get(n.StopId) {
		// skip blocked routes
		if n.Blocked.Contains(route.Id) {
			continue
		}

		// route next stop time
		next, err := p.ScheduleIndex.Get(n.StopId, route.Id).Next(n.Arrival)
		if err != nil {
			continue
		}

		// route wait duration
		waitDuration := stopTimeDiffDuration(n.Arrival, next.Time)

		// stop times reached by the next.tripId
		for _, visit := range p.visits(next) {
			current, seen := fastest[visit.StopId]
			transitDuration := stopTimeDiffDuration(next.Time, visit.Time)
			totalDuration := waitDuration + transitDuration

			// update the fastest trip
			if !seen || totalDuration < current.Duration() {
				fastest[visit.StopId] = FastestTrip{
					Arrival:         n.Arrival.Add(totalDuration),
					TransitDuration: transitDuration,
					WaitDuration:    waitDuration,
					TripId:          next.TripId,
				}
			}

			// update the stop blockers
			if _, ok := blockers[visit.StopId]; !ok {
				blockers[visit.StopId] = Set{}
			}
			stopBlockers := blockers[visit.StopId]
			stopBlockers.Add(route.Id)

			// update the stops seen
			stops.Add(visit.StopId)
		}
	}

	return &transitNeighbors{
		stops:    stops,
		blockers: blockers,
		fastest:  fastest,
	}
}

func (p *Planner) nodesTransitNeighbors(prev *Node, t *transitNeighbors) []*Node {
	connections := []*Node{}

	for stopId, trip := range t.fastest {
		node := &Node{
			Previous: &Edge{
				Node:            prev,
				TripId:          trip.TripId,
				TransitDuration: trip.TransitDuration,
				WaitDuration:    trip.WaitDuration,
			},
			Blocked: t.blockers[stopId],
			StopId:  stopId,
			Arrival: trip.Arrival,
		}
		connections = append(connections, node)
	}

	return connections
}

func (p *Planner) visitWalkingNeighbors(prev *Node, visited Set) []*Node {
	// create the nodes
	neighbors := []*Node{}
	stop, _ := p.StopIndex.Get(prev.StopId)

	for _, neighbor := range p.StopLocationIndex.Query(stop.Location, 250) {
		// the neighbor can't be accessible by transit
		if visited.Contains(neighbor.ID()) {
			continue
		}

		// neighbor information
		neighborRoutes := p.StopRouteIndex.Get(neighbor.ID())
		neighborDistance := neighbor.Location.Distance(stop.Location)

		// the neighbor must contain a new route
		neighborNewRoutes := false
		for _, route := range neighborRoutes {
			if !prev.Blocked.Contains(route.ID()) {
				neighborNewRoutes = true
				break
			}
		}

		if !neighborNewRoutes {
			continue
		}

		// neighbor node
		transitDuration := time.Duration(neighborDistance) * time.Second

		neighborNode := &Node{
			Previous: &Edge{
				Node:            prev,
				Walking:         true,
				TransitDuration: transitDuration,
			},
			StopId:  neighbor.Id,
			Arrival: prev.Arrival.Add(transitDuration),
			Blocked: prev.Blocked, // not sure about this
		}

		neighbors = append(neighbors, neighborNode)
	}

	return neighbors
}

func (p *Planner) walkingNeighbors(transitNodes []*Node, stops Set) []*Node {
	walkingNodes := map[string]*Node{}

	for _, transitNode := range transitNodes {
		for _, walkingNode := range p.visitWalkingNeighbors(transitNode, stops) {
			stopId := walkingNode.StopId

			// no existing walking node
			existing, ok := walkingNodes[stopId]
			if !ok {
				walkingNodes[stopId] = walkingNode
				continue
			}

			// closer than existing walking node
			if existing.Arrival.After(walkingNode.Arrival) {
				walkingNodes[stopId] = walkingNode
			}
		}
	}

	// add the walking nodes
	nodes := []*Node{}
	for _, walkingNode := range walkingNodes {
		nodes = append(nodes, walkingNode)
	}
	return nodes
}

func stopTimeDiffDuration(from, to time.Time) time.Duration {
	f := from.Hour()*60 + from.Minute()
	t := to.Hour()*60 + to.Minute()

	if t < f {
		t += 60 * 24
	}

	return time.Duration(t-f) * time.Minute
}
