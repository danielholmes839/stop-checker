package db

import (
	"sort"
)

type StopRank struct {
	StopLocationResult
	BestRoute     string
	BestRouteRank int
	RouteCount    int
}

type StopRanker struct {
	stopRoutes *StopRouteIndex
}

func NewStopRanker(stopRoutes *StopRouteIndex) *StopRanker {
	return &StopRanker{
		stopRoutes: stopRoutes,
	}
}

type tracker struct {
	position int
	distance float64
}

/* Ranks stops based on a few criteria:
- StopRank.Rank = 0 means this stop is the closest stop for at least one route, 1 would be second closest for at least one route, etc
- StopRank.Rank is used to sort the stops
- if the StopRank.Rank is the same between two stops. The distance from location subtracted by an offset for the number of routes at each stop
	this way stops with more routes are given higher priority

WARNING: does not preserve order by distance. stops passed in must be sorted by distance.
*/
func (ranker *StopRanker) Rank(stops []StopLocationResult) []StopRank {
	// closest stop by route id
	closest := map[string]*tracker{}
	ranked := []StopRank{}

	for _, stop := range stops {
		stopRoutes := ranker.stopRoutes.Get(stop.Id)
		rank := &StopRank{
			StopLocationResult: stop,
			BestRouteRank:      99,
			RouteCount:         len(stopRoutes),
		}

		for _, stopRoute := range stopRoutes {
			id := stopRoute.ID()

			// first time seeing this route
			if _, seen := closest[id]; !seen {
				closest[id] = &tracker{
					distance: stop.Distance,
					position: 0,
				}
				rank.BestRouteRank = 0
				rank.BestRoute = stopRoute.Route.ID()
				break
			}

			// increment the tracker when another stop uses this route
			current := closest[id]
			current.distance = stop.Distance
			current.position++

			if current.position < rank.BestRouteRank {
				rank.BestRouteRank = 0
				rank.BestRoute = stopRoute.Route.ID()
			}
		}

		ranked = append(ranked, *rank)
	}

	// rank the stops
	sort.Slice(ranked, func(i, j int) bool {
		a := ranked[i]
		b := ranked[j]

		if a.BestRoute < b.BestRoute {
			return true
		}

		aDist := a.Distance - ((float64(a.RouteCount) - 1) * 150)
		bDist := b.Distance - ((float64(b.RouteCount) - 1) * 150)
		if a.BestRouteRank == b.BestRouteRank && aDist < bDist {
			return true
		}

		return false
	})

	return ranked
}