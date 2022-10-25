package v2

import (
	"math"
	"time"
)

const TRANSFER_PENALTY = time.Minute * 5

// value multiplied by distance remaining to the stop (2.5 minute penalty per km away)
const DISTANCE_PENALTY = (5 * time.Minute) / (2 * 1000)

// the maximum distance to walk from origin to first stop, or last stop to destination
const MAX_WALK = 500.0

// the maximum distance to walk between stops
const MAX_WALK_EXPLORE = 300.0

const WALK_SPEED = 1.4 // meters per second
const WALK_PENALTY = 0.5

type kind int

const (
	INITIAL kind = iota
	TARGET
	STOP
)

func walkingDuration(distance float64) time.Duration {
	duration := time.Duration(math.Round(distance/(WALK_SPEED*60))) * time.Minute
	if duration < time.Minute {
		return time.Minute
	}
	return duration
}
