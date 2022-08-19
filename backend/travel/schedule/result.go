package schedule

import (
	"time"

	"stop-checker.com/db/model"
)

type Result struct {
	model.StopTime
	// DateTime compoenent
	time.Time
}

// NewResultAfter - stoptime is after base
func NewResultAfter(base time.Time, stopTime model.StopTime) Result {
	duration := StopTimeDiff(base, stopTime.Time)

	return Result{
		StopTime: stopTime,
		Time:     base.Add(duration),
	}
}

// NewResultBefore - stoptime is before base
func NewResultBefore(base time.Time, stopTime model.StopTime) Result {
	duration := StopTimeDiff(stopTime.Time, base)

	return Result{
		StopTime: stopTime,
		Time:     base.Add(-duration),
	}
}
