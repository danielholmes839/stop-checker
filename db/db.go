package db

import "stop-checker.com/db/model"

type InMemoryDB struct {
	StopTimesByTripId *InvertedIndex[model.StopTime]
}

func NewInMemoryDB(base *model.Base) *InMemoryDB {
	return &InMemoryDB{
		StopTimesByTripId: NewInvertedIndex(base.StopTimes, func(stopTime model.StopTime) (key string) {
			return stopTime.TripID
		}),
	}
}
