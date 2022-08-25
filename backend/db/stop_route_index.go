package db

import (
	"stop-checker.com/model"
)

type StopRouteIndex struct {
	index *InvertedIndex[model.StopRoute]
}

type stopRouteInfo struct {
	directionId, headsign string
}

func NewStopRouteIndex(indexes *BaseIndex, base *model.Base) *StopRouteIndex {
	routes := make(map[string]map[string]stopRouteInfo)

	// create a map of unique route ids for each stop id
	for _, stopTime := range base.StopTimes {
		trip, _ := indexes.Trips.Get(stopTime.TripId)
		stopId := stopTime.StopId
		routeId := trip.RouteId

		if _, ok := routes[stopId]; !ok {
			routes[stopId] = map[string]stopRouteInfo{}
		}
		routes[stopId][routeId] = stopRouteInfo{
			directionId: trip.DirectionId,
			headsign:    trip.Headsign,
		}
	}

	index := &InvertedIndex[model.StopRoute]{
		data: map[string][]model.StopRoute{},
	}

	// add data to the index
	for stopId, routes := range routes {
		index.data[stopId] = []model.StopRoute{}
		for routeId, info := range routes {
			index.data[stopId] = append(index.data[stopId], model.StopRoute{
				RouteId:     routeId,
				StopId:      stopId,
				DirectionId: info.directionId,
				Headsign:    info.headsign,
			})
		}
	}

	return &StopRouteIndex{
		index: index,
	}
}

func (s *StopRouteIndex) Get(stopId string) []model.StopRoute {
	results, _ := s.index.Get(stopId)
	return results
}
