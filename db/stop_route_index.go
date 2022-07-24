package db

import (
	"fmt"

	"stop-checker.com/db/model"
)

type StopRoute struct {
	model.Route
	DirectionId string
}

func (r StopRoute) ID() string {
	return fmt.Sprintf("%s:%s", r.Id, r.DirectionId)
}

type StopRouteIndex struct {
	index *InvertedIndex[StopRoute]
}

func NewStopRouteIndex(indexes *BaseIndex, base *model.Base) *StopRouteIndex {
	routes := make(map[string]map[string]string)

	// create a map of unique route ids for each stop id
	for _, stopTime := range base.StopTimes {
		trip, _ := indexes.Trips.Get(stopTime.TripId)
		stopId := stopTime.StopId
		routeId := trip.RouteId

		if _, ok := routes[stopId]; !ok {
			routes[stopId] = map[string]string{}
		}
		routes[stopId][routeId] = trip.DirectionId
	}

	index := &InvertedIndex[StopRoute]{
		data: map[string][]StopRoute{},
	}

	// add data to the index
	for stopId, routes := range routes {
		index.data[stopId] = []StopRoute{}
		for routeId, directionId := range routes {
			route, _ := indexes.Routes.Get(routeId)
			index.data[stopId] = append(index.data[stopId], StopRoute{
				Route:       route,
				DirectionId: directionId,
			})
		}
	}

	return &StopRouteIndex{
		index: index,
	}
}

func (s *StopRouteIndex) Get(stopId string) []StopRoute {
	results, _ := s.index.Get(stopId)
	return results
}
