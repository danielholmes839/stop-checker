package model

import (
	"fmt"
	"sort"
	"time"
)

func StopTimeSort(stopTimes []StopTime) {
	sort.Slice(stopTimes, func(i, j int) bool {
		return stopTimes[i].Time < stopTimes[j].Time
	})
}

type Time int

func NewTimeFromDateTime(t time.Time) Time {
	return NewTime(t.Hour(), t.Minute())
}

func NewTime(hours, minutes int) Time {
	return Time(hours*60 + minutes)
}

func (t Time) Hour() int {
	return int(t / 60)
}

func (t Time) Minute() int {
	return int(t % 60)
}

// t is after dt
func (t Time) After(dt time.Time) bool {
	return NewTimeFromDateTime(dt) <= t
}

// t is before dt
func (t Time) Before(dt time.Time) bool {
	return NewTimeFromDateTime(dt) >= t
}

func TimeDiff(t0, t1 Time) time.Duration {
	if t1 < t0 {
		t1 += 60 * 24
	}

	return time.Duration(t1-t0) * time.Minute
}

func (t Time) String() string {
	h := t.Hour()
	m := t.Minute()

	meridiem := "AM"

	if h >= 12 {
		meridiem = "PM"
	}

	h %= 12
	if h == 0 {
		h = 12
	}

	padding := ""
	if m < 10 {
		padding = "0"
	}

	return fmt.Sprintf("%d:%s%d %s", h, padding, m, meridiem)
}
