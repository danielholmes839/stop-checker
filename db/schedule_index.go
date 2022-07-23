package db

import (
	"fmt"
	"sort"

	"stop-checker.com/db/model"
)

type ScheduleIndex struct {
	index *InvertedIndex[model.StopTime]
}

func NewScheduleIndex(indexes *BaseIndex, base *model.Base) *ScheduleIndex {
	// create the schedule index
	index := NewInvertedIndex(base.StopTimes, func(stopTime model.StopTime) (key string) {
		trip, _ := indexes.Trips.Get(stopTime.TripId)
		return fmt.Sprintf("%s:%s", stopTime.StopId, trip.RouteId)
	})

	// sort the stop times by arrival time
	for _, schedule := range index.data {
		sort.Slice(schedule, func(i, j int) bool {
			return schedule[i].Time.Before(schedule[j].Time)
		})
	}

	return &ScheduleIndex{index: index}
}

func (schedule *ScheduleIndex) Get(stopId, routeId string) []model.StopTime {
	results, _ := schedule.index.Get(fmt.Sprintf("%s:%s", stopId, routeId))
	return results
}
