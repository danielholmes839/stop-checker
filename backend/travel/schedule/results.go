package schedule

import (
	"errors"
	"fmt"
	"time"

	"stop-checker.com/db/model"
)

type Results struct {
	*requiredIndexes
	results []model.StopTime
}

func (r *Results) Next(after time.Time) (Result, error) {
	results := r.After(after, 1)

	if len(results) == 0 {
		return Result{}, errors.New("not found")
	}

	return results[0], nil
}

func (r *Results) Previous(before time.Time) (Result, error) {
	results := r.Before(before, 1)

	if len(results) == 0 {
		return Result{}, errors.New("not found")
	}

	return results[0], nil
}

func (r *Results) Before(before time.Time, limit int) []Result {
	results := r.beforeWithinDay(before, limit)

	before = truncate(before).Add(-time.Second)
	attempts := 0
	for len(results) != limit && attempts < 3 {
		fmt.Println(before)
		attempts++
		results = append(results, r.beforeWithinDay(before, limit-len(results))...)
		before = before.AddDate(0, 0, -1)
	}

	return results
}

// query next N stop times after a specific time
func (r *Results) After(after time.Time, limit int) []Result {
	results := r.afterWithinDay(after, limit)

	attempts := 0
	for len(results) != limit && attempts < 3 {
		attempts++
		after = truncate(after.AddDate(0, 0, 1))
		results = append(results, r.afterWithinDay(after, limit-len(results))...)
	}

	return results
}

// query all stop times on a specific date
func (r *Results) Day(on time.Time) []Result {
	return r.afterWithinDay(truncate(on), -1)
}

func (r *Results) afterWithinDay(t time.Time, limit int) []Result {
	results := []Result{}

	for _, stopTime := range r.results {
		if !after(stopTime.Time, t) {
			continue
		}

		if !r.valid(t, stopTime) {
			continue
		}

		result := NewResultAfter(t, stopTime)
		results = append(results, result)

		if len(results) == limit {
			break
		}
	}

	return results
}

func (r *Results) beforeWithinDay(t time.Time, limit int) []Result {
	results := []Result{}

	for _, stopTime := range reverse(r.results) {
		if !before(stopTime.Time, t) {
			continue
		}

		if !r.valid(t, stopTime) {
			continue
		}

		result := NewResultBefore(t, stopTime)
		results = append(results, result)

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
func (r *Results) valid(t time.Time, stopTime model.StopTime) bool {
	if stopTime.Overflow {
		t = t.Add(time.Hour * -24)
	}

	// lookup the trip and service records
	trip, _ := r.trips.Get(stopTime.TripId)
	service, _ := r.services.Get(trip.ServiceId)

	// results must be between the start and end dates
	if service.End.Before(t) || service.Start.After(t) {
		return false
	}

	// results must have service on the given day
	if !service.On[t.Weekday()] {
		return false
	}

	// results must not have execpetions
	if exception, err := r.serviceExceptions.Get(service.Id, t); err == nil && !exception.Added {
		return false
	}

	return true
}
