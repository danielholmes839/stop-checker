package db

import (
	"errors"
	"time"

	"stop-checker.com/db/model"
)

type ScheduleResult struct {
	model.StopTime
	time.Time
}

type ScheduleResults struct {
	*indexesRequiredBySchedule
	results []model.StopTime
}

func (s *ScheduleResults) Next(after time.Time) (ScheduleResult, error) {
	results := s.After(after, 1)

	if len(results) == 0 {
		return ScheduleResult{}, errors.New("not found")
	}

	return results[0], nil
}

func (s *ScheduleResults) Previous(before time.Time) (ScheduleResult, error) {
	results := s.Before(before, 1)

	if len(results) == 0 {
		return ScheduleResult{}, errors.New("not found")
	}

	return results[0], nil
}

func (s *ScheduleResults) Before(before time.Time, limit int) []ScheduleResult {
	results := s.beforeWithinDay(before, limit)

	before = truncate(before).Add(-time.Second)
	attempts := 0
	for len(results) != limit && attempts < 7 {
		attempts++
		results = append(results, s.beforeWithinDay(before, limit-len(results))...)
		before = before.AddDate(0, 0, -1)
	}

	return results
}

// query next N stop times after a specific time
func (s *ScheduleResults) After(after time.Time, limit int) []ScheduleResult {
	results := s.afterWithinDay(after, limit)

	attempts := 0
	for len(results) != limit && attempts < 7 {
		attempts++
		after = truncate(after.AddDate(0, 0, 1))
		results = append(results, s.afterWithinDay(after, limit-len(results))...)
	}

	return results
}

// query all stop times on a specific date
func (s *ScheduleResults) Day(on time.Time) []ScheduleResult {
	return s.afterWithinDay(truncate(on), -1)
}

func (s *ScheduleResults) afterWithinDay(t time.Time, limit int) []ScheduleResult {
	results := []ScheduleResult{}
	year, month, day := t.Date()

	for _, stopTime := range s.results {
		if !stopTime.After(t) {
			continue
		}

		if !s.valid(t, stopTime) {
			continue
		}

		dt := time.Date(year, month, day, stopTime.Hour(), stopTime.Minute(), 0, 0, time.Local)

		results = append(results, ScheduleResult{
			StopTime: stopTime,
			Time:     dt,
		})

		if len(results) == limit {
			break
		}
	}

	return results
}

func (s *ScheduleResults) beforeWithinDay(t time.Time, limit int) []ScheduleResult {
	results := []ScheduleResult{}
	year, month, day := t.Date()

	for _, stopTime := range reverse(s.results) {
		if !stopTime.Before(t) {
			continue
		}

		if !s.valid(t, stopTime) {
			continue
		}

		results = append(results, ScheduleResult{
			StopTime: stopTime,
			Time:     time.Date(year, month, day, stopTime.Hour(), stopTime.Minute(), 0, 0, time.Local),
		})

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
func (s *ScheduleResults) valid(date time.Time, stopTime model.StopTime) bool {
	if stopTime.Overflow {
		/* stop times can overflow to the next day. service day of
		2022-08-28 and time of 26:00 means 2AM on 2022-08-29 so to check
		if the stoptime will happen on the 29th we actually check the 28th*/
		date = date.Add(time.Hour * -24)
	}

	trip, _ := s.trips.Get(stopTime.TripId)
	service, _ := s.services.Get(trip.ServiceId)

	// results must be between the start and end dates
	if service.End.Before(date) || service.Start.After(date) {
		return false
	}

	// results must have service on the week day
	if !service.On[date.Weekday()] {
		return false
	}

	// results must not have execpetions
	if exception, err := s.serviceExceptions.Get(service.Id, date); err == nil && !exception.Added {
		return false
	}

	return true
}

// truncate time leaving only the date
func truncate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func reverse(data []model.StopTime) []model.StopTime {
	reversed := make([]model.StopTime, len(data))
	for i, stopTime := range data {
		reversed[len(data)-(i+1)] = stopTime
	}
	return reversed
}
