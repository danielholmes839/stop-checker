package schedule

import (
	"testing"
	"time"

	"stop-checker.com/db/model"
)

func TestResult(t *testing.T) {
	loc, _ := time.LoadLocation("America/Montreal")
	format := "2006-01-02T15:04:00"

	t.Run("after tc1", func(t *testing.T) {
		base, _ := time.ParseInLocation(format, "2022-08-18T00:00:00", loc)
		stoptime, _ := time.ParseInLocation(format, "0001-01-01T12:00:00", loc)
		result := NewResultAfter(base, model.StopTime{Time: stoptime})

		assertEqualDate(t, base, result.Time)
		assertEqualTime(t, stoptime, result.Time)
	})

	t.Run("after tc2", func(t *testing.T) {
		base, _ := time.ParseInLocation(format, "2022-08-18T23:59:00", loc)
		stoptime, _ := time.ParseInLocation(format, "0001-01-01T00:00:00", loc) // 15 hour difference - date rolls over
		result := NewResultAfter(base, model.StopTime{Time: stoptime})

		assertEqualDate(t, base.Add(time.Hour*24), result.Time)
		assertEqualTime(t, stoptime, result.Time)
	})

	t.Run("before tc1", func(t *testing.T) {
		base, _ := time.ParseInLocation(format, "2022-08-18T12:00:00", loc)
		stoptime, _ := time.ParseInLocation(format, "0001-01-01T11:00:00", loc)
		result := NewResultBefore(base, model.StopTime{Time: stoptime})

		assertEqualDate(t, base, result.Time)
		assertEqualTime(t, stoptime, result.Time)
	})

	t.Run("before tc2", func(t *testing.T) {
		base, _ := time.ParseInLocation(format, "2022-08-18T12:00:00", loc)
		stoptime, _ := time.ParseInLocation(format, "0001-01-01T13:00:00", loc)
		result := NewResultBefore(base, model.StopTime{Time: stoptime})

		assertEqualDate(t, base.Add(-time.Hour*24), result.Time)
		assertEqualTime(t, stoptime, result.Time)
	})

	// t.Run("before tc2", func(t *testing.T) {
	// 	base, _ := time.ParseInLocation(format, "2022-08-18T23:59:00", loc)
	// 	stoptime, _ := time.ParseInLocation(format, "2000-00-00T12:00:00", loc)
	// 	result := NewResultBefore(base, model.StopTime{Time: stoptime})

	// 	assertEqualDate(t, base, result.Time)
	// 	assertEqualTime(t, stoptime, result.Time)
	// })

}
