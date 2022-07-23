package db

import (
	"fmt"
	"sort"
	"time"

	"stop-checker.com/db/model"
)

type indexesRequiredBySchedule struct {
	trips             *Index[model.Trip]
	services          *Index[model.Service]                  // services by id
	serviceExceptions *InvertedIndex[model.ServiceException] // service exceptions by service ID
}

type ScheduleIndex struct {
	*indexesRequiredBySchedule
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

	return &ScheduleIndex{
		index: index,
		indexesRequiredBySchedule: &indexesRequiredBySchedule{
			trips:    indexes.Trips,
			services: indexes.Services,
			serviceExceptions: NewInvertedIndex(base.ServiceExceptions, func(record model.ServiceException) (key string) {
				return record.ServiceId
			}),
		},
	}
}

func (schedule *ScheduleIndex) Get(stopId, routeId string) *ScheduleResults {
	results, _ := schedule.index.Get(fmt.Sprintf("%s:%s", stopId, routeId))
	return &ScheduleResults{
		results:                   results,
		indexesRequiredBySchedule: schedule.indexesRequiredBySchedule,
	}
}

type ScheduleResults struct {
	*indexesRequiredBySchedule
	results []model.StopTime
}

//
func (s *ScheduleResults) Next(after time.Time, limit int) *ScheduleResults {
	results := s.nextToday(after, limit)

	attempts := 0
	for len(results) != limit && attempts < 14 {
		after = after.AddDate(0, 0, 1)
		after.Truncate(time.Hour * 24)
		results = append(results, s.nextToday(after, limit-len(results))...)
	}

	return &ScheduleResults{
		indexesRequiredBySchedule: s.indexesRequiredBySchedule,
		results:                   results,
	}
}

func (s *ScheduleResults) nextToday(after time.Time, limit int) []model.StopTime {
	results := []model.StopTime{}

	for _, stopTime := range s.results {
		if !ahead(stopTime.Time, after) {
			continue
		}

		trip, _ := s.trips.Get(stopTime.TripId)
		service, _ := s.services.Get(trip.ServiceId)

		// results must be between the start and end dates
		if service.End.Before(after) || service.Start.After(after) {
			continue
		}

		// results must have service on the given day
		if !service.On[after.Weekday()] {
			continue
		}

		// results must not have execpetions
		// - todo... may want a specialized index for serviceid-date

		results = append(results, stopTime)

		if len(results) == limit {
			break
		}
	}

	return results
}

func (s *ScheduleResults) Return() []model.StopTime {
	return s.results
}

func ahead(a, b time.Time) bool {
	// if a is ahead of b
	if a.Hour() >= b.Hour() {
		return true
	} else if a.Hour() == b.Hour() && a.Hour() > b.Hour() {
		return true
	}
	return false
}
