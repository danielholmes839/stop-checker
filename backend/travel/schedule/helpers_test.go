package schedule

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func assertEqualDate(t *testing.T, t1, t2 time.Time) {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()

	assert.Equal(t, y1, y2)
	assert.Equal(t, m1, m2)
	assert.Equal(t, d1, d2)
}

func assertEqualTime(t *testing.T, t1, t2 time.Time) {
	h1, m1 := t1.Hour(), t1.Minute()
	h2, m2 := t2.Hour(), t2.Minute()

	assert.Equal(t, h1, h2)
	assert.Equal(t, m1, m2)
}

func assertMidnight(t *testing.T, t1 time.Time) {
	assert.Equal(t, 0, t1.Minute())
	assert.Equal(t, 0, t1.Hour())
}

func TestTruncate(t *testing.T) {
	loc, _ := time.LoadLocation("America/Montreal")

	t.Run("midnight", func(t *testing.T) {
		midnight, _ := time.ParseInLocation("2006-01-02T15:04:00", "2022-08-18T00:00:00", loc)
		truncated := truncate(midnight)

		assertEqualDate(t, midnight, truncated)
		assertMidnight(t, truncated)
	})

	t.Run("noon", func(t *testing.T) {
		noon, _ := time.ParseInLocation("2006-01-02T15:04:00", "2022-08-18T00:12:00", loc)
		truncated := truncate(noon)

		assertEqualDate(t, noon, truncated)
		assertMidnight(t, truncated)
	})
}

func TestStopTimeDiff(t *testing.T) {
	loc, _ := time.LoadLocation("America/Montreal")
	format := "2006-01-02T15:04:00"

	t.Run("t1 < t2 (same day)", func(t *testing.T) {
		t1, _ := time.ParseInLocation(format, "2022-08-18T12:00:00", loc)
		t2, _ := time.ParseInLocation(format, "2022-08-18T13:30:00", loc)

		duration := StopTimeDiff(t1, t2)
		assert.Equal(t, time.Minute*90, duration)
	})

	t.Run("t1 < t2 (next day)", func(t *testing.T) {
		t1, _ := time.ParseInLocation(format, "2022-08-18T12:00:00", loc)
		t2, _ := time.ParseInLocation(format, "2022-08-19T13:30:00", loc)

		duration := StopTimeDiff(t1, t2)
		assert.Equal(t, time.Minute*90, duration)
	})

	t.Run("t1 > t2 (same day)", func(t *testing.T) {
		t1, _ := time.ParseInLocation(format, "2000-01-01T23:00:00", loc)
		t2, _ := time.ParseInLocation(format, "2022-08-19T00:30:00", loc)

		duration := StopTimeDiff(t1, t2)
		assert.Equal(t, time.Minute*90, duration)
	})

	t.Run("t1 > t2 (different day)", func(t *testing.T) {
		t1, _ := time.ParseInLocation(format, "2023-09-01T23:00:00", loc)
		t2, _ := time.ParseInLocation(format, "2022-08-19T00:30:00", loc)

		duration := StopTimeDiff(t1, t2)
		assert.Equal(t, time.Minute*90, duration)
	})

	t.Run("t1 == t2 (different day)", func(t *testing.T) {
		t1, _ := time.ParseInLocation(format, "2020-08-18T12:00:00", loc)
		t2, _ := time.ParseInLocation(format, "2025-08-19T12:00:00", loc)

		duration := StopTimeDiff(t1, t2)
		assert.Equal(t, time.Duration(0), duration)
	})
}
