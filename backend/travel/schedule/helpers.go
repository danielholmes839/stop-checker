package schedule

import (
	"time"

	"stop-checker.com/db/model"
)

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
	} else if a.Hour() == b.Hour() && a.Minute() <= b.Minute() { // minutes can be equal
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

// to.Sub(from) ignoring date
func StopTimeDiff(from, to time.Time) time.Duration {
	f := from.Hour()*60 + from.Minute()
	t := to.Hour()*60 + to.Minute()

	if t < f {
		t += 60 * 24
	}

	return time.Duration(t-f) * time.Minute
}
