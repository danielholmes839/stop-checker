package db

import (
	"stop-checker.com/db/model"
)

type StopRouteIndex struct {
	index *InvertedIndex[model.StopRoute]
}

type stopRouteInfo struct {
	directionId string
	headsign    string
}

func NewStopRouteIndex(trips *Index[model.Trip], stopTimes []model.StopTime) *StopRouteIndex {
	stopRoutes := make(map[string]map[string]stopRouteInfo)

	// create a map of unique route ids for each stop id
	for _, stopTime := range stopTimes {
		trip, _ := trips.Get(stopTime.TripId)
		stopId := stopTime.StopId
		routeId := trip.RouteId

		if _, ok := stopRoutes[stopId]; !ok {
			stopRoutes[stopId] = map[string]stopRouteInfo{}
		}
		stopRoutes[stopId][routeId] = stopRouteInfo{
			directionId: trip.DirectionId,
			headsign:    trip.Headsign,
		}
	}

	index := &InvertedIndex[model.StopRoute]{
		data: map[string][]model.StopRoute{},
	}

	// add data to the index
	for stopId, routes := range stopRoutes {
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
