package db

import (
	"fmt"
	"sort"

	"stop-checker.com/db/model"
)

type indexesRequiredBySchedule struct {
	trips             *Index[model.Trip]
	services          *Index[model.Service]  // services by id
	serviceExceptions *ServiceExceptionIndex // service exceptions by service ID and time
}

type ScheduleIndex struct {
	*indexesRequiredBySchedule
	index *InvertedIndex[model.StopTime]
}

func NewScheduleIndex(indexes *BaseIndex, base *model.Dataset) *ScheduleIndex {
	// create the schedule index
	index := NewInvertedIndex("schedule", base.StopTimes, func(stopTime model.StopTime) (key string) {
		trip, _ := indexes.Trips.Get(stopTime.TripId)
		return fmt.Sprintf("%s:%s", stopTime.StopId, trip.RouteId)
	})

	// sort the stop times by arrival time
	for _, schedule := range index.data {
		sort.Slice(schedule, func(i, j int) bool {
			return schedule[i].Time < schedule[j].Time
		})
	}

	return &ScheduleIndex{
		index: index,
		indexesRequiredBySchedule: &indexesRequiredBySchedule{
			trips:             indexes.Trips,
			services:          indexes.Services,
			serviceExceptions: indexes.ServiceExeceptions,
		},
	}
}

func (schedule *ScheduleIndex) Get(stopId, routeId string) *ScheduleResults {
	results, _ := schedule.index.Get(fmt.Sprintf("%s:%s", stopId, routeId))
	return &ScheduleResults{
		indexesRequiredBySchedule: schedule.indexesRequiredBySchedule,
		results:                   results,
	}
}
