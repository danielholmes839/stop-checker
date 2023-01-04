package resolvers

import (
	"context"

	"stop-checker.com/db/model"
	"stop-checker.com/db/repository"
)

type TripResolvers struct {
	repository.Routes
	repository.Shapes
	repository.Services
	StopTimesByTrip repository.InvertedIndex[model.StopTime]
}

func (r *TripResolvers) ID(ctx context.Context, obj *model.Trip) (string, error) {
	return obj.ID(), nil
}

func (r *TripResolvers) Route(ctx context.Context, obj *model.Trip) (model.Route, error) {
	return r.Routes.Get(obj.RouteId)
}

func (r *TripResolvers) Stoptimes(ctx context.Context, obj *model.Trip) ([]model.StopTime, error) {
	return r.StopTimesByTrip.Get(obj.ID())
}

func (r *TripResolvers) Shape(ctx context.Context, obj *model.Trip) ([]model.Location, error) {
	shapes, _ := r.Shapes.Get(obj.ShapeId)
	locations := make([]model.Location, len(shapes))

	for i, shape := range shapes {
		locations[i] = shape.Location
	}

	return locations, nil
}

func (r *TripResolvers) Service(ctx context.Context, obj *model.Trip) (model.Service, error) {
	return r.Services.Get(obj.ServiceId)
}

func (r *TripResolvers) Direction(ctx context.Context, obj *model.Trip) (string, error) {
	return obj.DirectionId, nil
}
