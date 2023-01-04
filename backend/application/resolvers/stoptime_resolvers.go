package resolvers

import (
	"context"

	"stop-checker.com/db/model"
	"stop-checker.com/db/repository"
)

type StopTimeResolvers struct {
	repository.Trips
	repository.Stops
}

func (r *StopTimeResolvers) Stop(ctx context.Context, obj *model.StopTime) (model.Stop, error) {
	return r.Stops.Get(obj.StopId)
}

func (r *StopTimeResolvers) Trip(ctx context.Context, obj *model.StopTime) (model.Trip, error) {
	return r.Trips.Get(obj.TripId)
}
