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
	tripIndex := NewIndex(base.Trips, func(trip model.Trip) (key string) {
		return trip.ID
	})

	// create the schedule index
	scheduleIndex := &ScheduleIndex{
		index: NewInvertedIndex(base.StopTimes, func(stopTime model.StopTime) (key string) {
			trip, _ := tripIndex.Get(stopTime.TripID)
			return fmt.Sprintf("%s:%s", stopTime.StopID, trip.RouteID)
		}),
	}

	// sort the stop times by arrival time
	data := scheduleIndex.index.data
	for _, schedule := range data {
		sort.Slice(schedule, func(i, j int) bool {
			return schedule[i].Arrival.Before(schedule[j].Arrival)
		})
	}

	return scheduleIndex
}

func (schedule *ScheduleIndex) Get(stopId, routeId string) []model.StopTime {
	results, _ := schedule.index.Get(fmt.Sprintf("%s:%s", stopId, routeId))
	return results
}
