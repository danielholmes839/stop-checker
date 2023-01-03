package travel

import (
	"time"

	"stop-checker.com/db/model"
)

type fastestTransit struct {
	tripId       string
	routeId      string
	stopId       string
	stopArrival  time.Time
	stopLocation model.Location
}

func (f *fastestTransit) Faster(t time.Time, mode Mode) bool {
	if mode == DEPART_AT {
		return f.stopArrival.Before(t)
	}
	return f.stopArrival.After(t)
}
