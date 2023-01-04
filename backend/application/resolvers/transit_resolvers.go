package resolvers

import (
	"context"
	"time"

	"stop-checker.com/db/model"
	"stop-checker.com/db/repository"
)

type TransitResolvers struct {
	repository.Routes
	repository.Trips
	repository.Shapes
}

func (r *TransitResolvers) Route(ctx context.Context, obj *model.Transit) (model.Route, error) {
	return r.Routes.Get(obj.RouteId)
}

func (r *TransitResolvers) Trip(ctx context.Context, obj *model.Transit) (model.Trip, error) {
	return r.Trips.Get(obj.TripId)
}

func (r *TransitResolvers) Departure(ctx context.Context, obj *model.Transit) (time.Time, error) {
	return obj.OriginDeparture, nil
}

func (r *TransitResolvers) Duration(ctx context.Context, obj *model.Transit) (int, error) {
	return int(obj.TripDuration.Minutes()), nil
}

func (r *TransitResolvers) Wait(ctx context.Context, obj *model.Transit) (int, error) {
	return int(obj.WaitDuration.Minutes()), nil
}
