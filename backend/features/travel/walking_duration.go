package travel

import (
	"math"
	"time"
)

func walkingDuration(distance float64) time.Duration {
	duration := time.Duration(math.Round(distance*1.4/60)) * time.Minute
	if duration < time.Minute {
		return time.Minute
	}
	return duration
}
