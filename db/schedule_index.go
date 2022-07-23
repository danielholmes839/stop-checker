package db

import (
	"fmt"
	"sort"
	"time"

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

type ScheduleResults struct {
	*indexesRequiredBySchedule
	results []model.StopTime
}

func (s *ScheduleResults) Next(after time.Time, limit int) []model.StopTime {
	results := s.nextWithinDay(after, limit)

	attempts := 0
	for len(results) != limit && attempts < 14 {
		attempts++
		after = truncate(after.AddDate(0, 0, 1))
		results = append(results, s.nextWithinDay(after, limit-len(results))...)
	}

	return results
}

func (s *ScheduleResults) Day(on time.Time) []model.StopTime {
	return s.nextWithinDay(truncate(on), -1)
}

func (s *ScheduleResults) nextWithinDay(after time.Time, limit int) []model.StopTime {
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
		if exception, exists := s.serviceExceptions.Get(service.Id, after); exists && !exception.Added {
			continue
		}

		results = append(results, stopTime)

		if len(results) == limit {
			break
		}
	}

	return results
}

func ahead(a, b time.Time) bool {
	// if a is ahead of b
	if a.Hour() >= b.Hour() {
		return true
	} else if a.Hour() == b.Hour() && a.Minute() > b.Minute() {
		return true
	}
	return false
}

func truncate(t time.Time) time.Time {
	// truncate time leaving only the date
	_, offset := t.Zone()
	sub := time.Duration(offset) * time.Second
	return t.Add(-sub).Truncate(time.Hour * 24).Add(-sub)
}
