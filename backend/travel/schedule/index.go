package schedule

import (
	"fmt"

	"stop-checker.com/db"
	"stop-checker.com/db/model"
)

type requiredIndexes struct {
	trips             *db.Index[model.Trip]
	services          *db.Index[model.Service]  // services by id
	serviceExceptions *db.ServiceExceptionIndex // service exceptions by service ID and time
}

type Index struct {
	*requiredIndexes
	index *db.InvertedIndex[model.StopTime]
}

func NewIndex(indexes *db.BaseIndex, base *model.Base) *Index {
	// create the schedule index
	index := db.NewInvertedIndex("schedule", base.StopTimes, func(stopTime model.StopTime) (key string) {
		trip, _ := indexes.Trips.Get(stopTime.TripId)
		return fmt.Sprintf("%s:%s", stopTime.StopId, trip.RouteId)
	})

	return &Index{
		index: index,
		requiredIndexes: &requiredIndexes{
			trips:             indexes.Trips,
			services:          indexes.Services,
			serviceExceptions: indexes.ServiceExeceptions,
		},
	}
}

func (s *Index) Get(stopId, routeId string) *Results {
	results, _ := s.index.Get(fmt.Sprintf("%s:%s", stopId, routeId))
	return &Results{
		requiredIndexes: s.requiredIndexes,
		results:         results,
	}
}
