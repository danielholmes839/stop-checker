package db

import (
	"fmt"
	"sort"

	"stop-checker.com/db/model"
)

type ScheduleIndex struct {
	index *InvertedIndex[model.StopTime]
}

func NewScheduleIndex(base *model.Base) *ScheduleIndex {
	tripIndex := NewIndex(base.Trips)

	// create the schedule index
	scheduleIndex := &ScheduleIndex{
		index: NewInvertedIndex(base.StopTimes, func(stopTime model.StopTime) (key string) {
			trip, _ := tripIndex.Get(stopTime.TripId)
			return fmt.Sprintf("%s:%s", stopTime.StopId, trip.RouteId)
		}),
	}

	// sort the stop times by arrival time
	data := scheduleIndex.index.data
	for _, schedule := range data {
		sort.Slice(schedule, func(i, j int) bool {
			return schedule[i].Time.Before(schedule[j].Time)
		})
	}

	return scheduleIndex
}

func (schedule *ScheduleIndex) Get(stopId, routeId string) []model.StopTime {
	results, _ := schedule.index.Get(fmt.Sprintf("%s:%s", stopId, routeId))
	return results
}
