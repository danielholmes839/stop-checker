package travel

import (
	"math"
	"time"

	"stop-checker.com/db/model"
)

const MAX_WALK = 300.0
const MAX_WALK_INITIAL = 1000.0
const MAX_WALK_TARGET = 1000.0

const TRANSFER_PENALTY = 5 * time.Minute

// value multiplied by distance remaining to the stop (2.5 minute penalty per km away)
// typical A* heuristic
const DISTANCE_PENALTY = (2*time.Minute + 30*time.Second) / 1000

// walking constants
const WALK_SPEED = 1.3 // meters per second
const WALK_PENALTY = 0.50

type Kind int

const (
	INITIAL Kind = iota
	TARGET
	STOP
)

type Mode int

const (
	DEPART_AT Mode = iota
	ARRIVE_BY
)

func walkingDuration(distance float64) time.Duration {
	duration := time.Duration(math.Round(distance/(WALK_SPEED*60))) * time.Minute
	if duration < time.Minute {
		return time.Minute
	}
	return duration
}

func walkingDurationContinuous(distance float64) time.Duration {
	duration := time.Duration(distance) * time.Minute / (WALK_SPEED * 60)
	if duration < time.Minute {
		return time.Minute
	}
	return duration
}

type walkingDirections interface {
	GetDirections(origin, destination model.Location) (model.Path, error)
}

type walkingDirectionsCache interface {
	GetDirections(originId, destinationId string) (model.Path, error)
}
