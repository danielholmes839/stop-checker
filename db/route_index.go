package db

import "stop-checker.com/db/model"

type RouteIndex struct {
	index *InvertedIndex[model.Route]
}

func NewRouteIndex(base *model.Base) *RouteIndex {
	tripIndex := NewIndex(base.Trips, func(trip model.Trip) (key string) {
		return trip.ID
	})

	routeIndex := NewIndex(base.Routes, func(route model.Route) (key string) {
		return route.ID
	})

	routes := make(map[string]map[string]struct{})

	for _, stopTime := range base.StopTimes {
		trip, _ := tripIndex.Get(stopTime.TripID)
		stopId := stopTime.StopID
		routeId := trip.RouteID

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
			route, _ := routeIndex.Get(routeId)
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
