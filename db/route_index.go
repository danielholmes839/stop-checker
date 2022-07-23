package db

import "stop-checker.com/db/model"

type RouteIndex struct {
	index *InvertedIndex[model.Route]
}

func NewRouteIndex(indexes *BaseIndex, base *model.Base) *RouteIndex {
	routes := make(map[string]map[string]struct{})

	// create a map of unique route ids for each stop id
	for _, stopTime := range base.StopTimes {
		trip, _ := indexes.Trips.Get(stopTime.TripId)
		stopId := stopTime.StopId
		routeId := trip.RouteId

		if _, ok := routes[stopId]; !ok {
			routes[stopId] = map[string]struct{}{}
		}
		routes[stopId][routeId] = struct{}{}
	}

	index := &InvertedIndex[model.Route]{
		data: map[string][]model.Route{},
	}

	// add data to the index
	for stopId, routes := range routes {
		index.data[stopId] = []model.Route{}
		for routeId := range routes {
			route, _ := indexes.Routes.Get(routeId)
			index.data[stopId] = append(index.data[stopId], route)
		}
	}

	return &RouteIndex{
		index: index,
	}
}

func (routes *RouteIndex) Get(stopId string) []model.Route {
	results, _ := routes.index.Get(stopId)
	return results
}
