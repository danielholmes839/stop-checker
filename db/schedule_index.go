package db

import (
	"errors"
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

func (s *ScheduleResults) Next(after time.Time) (model.StopTime, error) {
	results := s.After(after, 1)

	if len(results) == 0 {
		return model.StopTime{}, errors.New("not found")
	}

	return results[0], nil
}

func (s *ScheduleResults) Previous(before time.Time) (model.StopTime, error) {
	results := s.Before(before, 1)

	if len(results) == 0 {
		return model.StopTime{}, errors.New("not found")
	}

	return results[0], nil
}

func (s *ScheduleResults) Before(before time.Time, limit int) []model.StopTime {
	results := s.beforeWithinDay(before, limit)

	before = truncate(before).Add(-time.Second)
	attempts := 0
	for len(results) != limit && attempts < 3 {
		fmt.Println(before)
		attempts++
		results = append(results, s.beforeWithinDay(before, limit-len(results))...)
		before = before.AddDate(0, 0, -1)
	}

	return results
}

// query next N stop times after a specific time
func (s *ScheduleResults) After(after time.Time, limit int) []model.StopTime {
	results := s.afterWithinDay(after, limit)

	attempts := 0
	for len(results) != limit && attempts < 3 {
		attempts++
		after = truncate(after.AddDate(0, 0, 1))
		results = append(results, s.afterWithinDay(after, limit-len(results))...)
	}

	return results
}

// query all stop times on a specific date
func (s *ScheduleResults) Day(on time.Time) []model.StopTime {
	return s.afterWithinDay(truncate(on), -1)
}

func (s *ScheduleResults) afterWithinDay(t time.Time, limit int) []model.StopTime {
	results := []model.StopTime{}

	for _, stopTime := range s.results {
		if !after(stopTime.Time, t) {
			continue
		}

		if !s.valid(t, stopTime) {
			continue
		}

		results = append(results, stopTime)

		if len(results) == limit {
			break
		}
	}

	return results
}

func (s *ScheduleResults) beforeWithinDay(t time.Time, limit int) []model.StopTime {
	results := []model.StopTime{}

	for _, stopTime := range reverse(s.results) {
		if !before(stopTime.Time, t) {
			continue
		}

		if !s.valid(t, stopTime) {
			continue
		}

		results = append(results, stopTime)

		if len(results) == limit {
			break
		}
	}

	return results
}

/* valid stopTime given a specific date
- the date is between the service start and end date
- the service is running on the time's day of the week
- there on no service exception on the time's date
*/
func (s *ScheduleResults) valid(t time.Time, stopTime model.StopTime) bool {
	if stopTime.Overflow {
		t = t.Add(time.Hour * -24)
	}

	trip, _ := s.trips.Get(stopTime.TripId)
	service, _ := s.services.Get(trip.ServiceId)

	// results must be between the start and end dates
	if service.End.Before(t) || service.Start.After(t) {
		return false
	}

	// results must have service on the given day
	if !service.On[t.Weekday()] {
		return false
	}

	// results must not have execpetions
	if exception, exists := s.serviceExceptions.Get(service.Id, t); exists && !exception.Added {
		return false
	}

	return true
}

// if a is ahead of b (hours and minutes)
func after(a, b time.Time) bool {
	if a.Hour() > b.Hour() {
		return true
	} else if a.Hour() == b.Hour() && a.Minute() >= b.Minute() { // minutes can be equal
		return true
	}
	return false
}

// if a is before of b (hours and minutes)
func before(a, b time.Time) bool {
	if a.Hour() < b.Hour() {
		return true
	} else if a.Hour() == b.Hour() && a.Minute() >= b.Minute() { // minutes can be equal
		return true
	}
	return false
}

// truncate time leaving only the date
func truncate(t time.Time) time.Time {
	_, offset := t.Zone()
	sub := time.Duration(offset) * time.Second
	return t.Add(-sub).Truncate(time.Hour * 24).Add(-sub)
}

func reverse(data []model.StopTime) []model.StopTime {
	reversed := make([]model.StopTime, len(data))
	for i, stopTime := range data {
		reversed[len(data)-(i+1)] = stopTime
	}
	return reversed
}
