package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"stop-checker.com/backend/graph/generated"
	"stop-checker.com/backend/graph/types"
	"stop-checker.com/db"
	"stop-checker.com/db/model"
)

// SearchStopText is the resolver for the searchStopText field.
func (r *queryResolver) SearchStopText(ctx context.Context, text string) ([]*model.Stop, error) {
	panic(fmt.Errorf("not implemented"))
}

// SearchStopLocation is the resolver for the searchStopLocation field.
func (r *queryResolver) SearchStopLocation(ctx context.Context, latitude float64, longitude float64, radius float64) ([]*db.StopLocationResult, error) {
	stops := r.StopLocationIndex.Query(model.Location{
		Latitude:  latitude,
		Longitude: longitude,
	}, radius)

	// return the stops
	return ref(stops), nil
}

// ID is the resolver for the id field.
func (r *routeResolver) ID(ctx context.Context, obj *model.Route) (string, error) {
	return obj.ID(), nil
}

// Text is the resolver for the text field.
func (r *routeResolver) Text(ctx context.Context, obj *model.Route) (string, error) {
	return obj.TextColor, nil
}

// Background is the resolver for the background field.
func (r *routeResolver) Background(ctx context.Context, obj *model.Route) (string, error) {
	return obj.Color, nil
}

// Type is the resolver for the type field.
func (r *routeResolver) Type(ctx context.Context, obj *model.Route) (types.RouteType, error) {
	panic(fmt.Errorf("not implemented"))
}

// Sunday is the resolver for the sunday field.
func (r *serviceResolver) Sunday(ctx context.Context, obj *model.Service) (bool, error) {
	return obj.On[time.Sunday], nil
}

// Monday is the resolver for the monday field.
func (r *serviceResolver) Monday(ctx context.Context, obj *model.Service) (bool, error) {
	return obj.On[time.Monday], nil
}

// Tuesday is the resolver for the tuesday field.
func (r *serviceResolver) Tuesday(ctx context.Context, obj *model.Service) (bool, error) {
	return obj.On[time.Tuesday], nil
}

// Wednesday is the resolver for the wednesday field.
func (r *serviceResolver) Wednesday(ctx context.Context, obj *model.Service) (bool, error) {
	return obj.On[time.Wednesday], nil
}

// Thursday is the resolver for the thursday field.
func (r *serviceResolver) Thursday(ctx context.Context, obj *model.Service) (bool, error) {
	return obj.On[time.Thursday], nil
}

// Friday is the resolver for the friday field.
func (r *serviceResolver) Friday(ctx context.Context, obj *model.Service) (bool, error) {
	return obj.On[time.Friday], nil
}

// Saturday is the resolver for the saturday field.
func (r *serviceResolver) Saturday(ctx context.Context, obj *model.Service) (bool, error) {
	return obj.On[time.Saturday], nil
}

// Start is the resolver for the start field.
func (r *serviceResolver) Start(ctx context.Context, obj *model.Service) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

// End is the resolver for the end field.
func (r *serviceResolver) End(ctx context.Context, obj *model.Service) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

// Exceptions is the resolver for the exceptions field.
func (r *serviceResolver) Exceptions(ctx context.Context, obj *model.Service) ([]*model.ServiceException, error) {
	panic(fmt.Errorf("not implemented"))
}

// Date is the resolver for the date field.
func (r *serviceExceptionResolver) Date(ctx context.Context, obj *model.ServiceException) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

// ID is the resolver for the id field.
func (r *stopResolver) ID(ctx context.Context, obj *model.Stop) (string, error) {
	return obj.ID(), nil
}

// Routes is the resolver for the routes field.
func (r *stopResolver) Routes(ctx context.Context, obj *model.Stop) ([]*model.StopRoute, error) {
	stopRoutes := r.StopRouteIndex.Get(obj.ID())
	return ref(stopRoutes), nil
}

// Stop is the resolver for the stop field.
func (r *stopRouteResolver) Stop(ctx context.Context, obj *model.StopRoute) (*model.Stop, error) {
	stop, _ := r.Stops.Get(obj.StopId)
	return &stop, nil
}

// Route is the resolver for the route field.
func (r *stopRouteResolver) Route(ctx context.Context, obj *model.StopRoute) (*model.Route, error) {
	route, _ := r.Routes.Get(obj.RouteId)
	return &route, nil
}

// Direction is the resolver for the direction field.
func (r *stopRouteResolver) Direction(ctx context.Context, obj *model.StopRoute) (string, error) {
	return obj.DirectionId, nil
}

// Schedule is the resolver for the schedule field.
func (r *stopRouteResolver) Schedule(ctx context.Context, obj *model.StopRoute) (*db.ScheduleResults, error) {
	return r.ScheduleIndex.Get(obj.StopId, obj.RouteId), nil
}

// Next is the resolver for the next field.
func (r *stopRouteScheduleResolver) Next(ctx context.Context, obj *db.ScheduleResults) ([]*model.StopTime, error) {
	now := time.Now().In(r.Timezone)
	stopTimes := obj.After(now, 3)
	return ref(stopTimes), nil
}

// Trip is the resolver for the trip field.
func (r *stopTimeResolver) Trip(ctx context.Context, obj *model.StopTime) (*model.Trip, error) {
	trip, _ := r.Trips.Get(obj.TripId)
	return &trip, nil
}

// Sequence is the resolver for the sequence field.
func (r *stopTimeResolver) Sequence(ctx context.Context, obj *model.StopTime) (int, error) {
	return obj.StopSeq, nil
}

// ID is the resolver for the id field.
func (r *tripResolver) ID(ctx context.Context, obj *model.Trip) (string, error) {
	return obj.ID(), nil
}

// Route is the resolver for the route field.
func (r *tripResolver) Route(ctx context.Context, obj *model.Trip) (*model.Route, error) {
	route, _ := r.Routes.Get(obj.RouteId)
	return &route, nil
}

// StopTimes is the resolver for the stopTimes field.
func (r *tripResolver) StopTimes(ctx context.Context, obj *model.Trip) ([]*model.StopTime, error) {
	stopTimes, _ := r.StopTimesFromTrip.Get(obj.Id)
	return ref(stopTimes), nil
}

// Service is the resolver for the service field.
func (r *tripResolver) Service(ctx context.Context, obj *model.Trip) (*model.Service, error) {
	service, _ := r.Services.Get(obj.ServiceId)
	return &service, nil
}

// Direction is the resolver for the direction field.
func (r *tripResolver) Direction(ctx context.Context, obj *model.Trip) (string, error) {
	return obj.DirectionId, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Route returns generated.RouteResolver implementation.
func (r *Resolver) Route() generated.RouteResolver { return &routeResolver{r} }

// Service returns generated.ServiceResolver implementation.
func (r *Resolver) Service() generated.ServiceResolver { return &serviceResolver{r} }

// ServiceException returns generated.ServiceExceptionResolver implementation.
func (r *Resolver) ServiceException() generated.ServiceExceptionResolver {
	return &serviceExceptionResolver{r}
}

// Stop returns generated.StopResolver implementation.
func (r *Resolver) Stop() generated.StopResolver { return &stopResolver{r} }

// StopRoute returns generated.StopRouteResolver implementation.
func (r *Resolver) StopRoute() generated.StopRouteResolver { return &stopRouteResolver{r} }

// StopRouteSchedule returns generated.StopRouteScheduleResolver implementation.
func (r *Resolver) StopRouteSchedule() generated.StopRouteScheduleResolver {
	return &stopRouteScheduleResolver{r}
}

// StopTime returns generated.StopTimeResolver implementation.
func (r *Resolver) StopTime() generated.StopTimeResolver { return &stopTimeResolver{r} }

// Trip returns generated.TripResolver implementation.
func (r *Resolver) Trip() generated.TripResolver { return &tripResolver{r} }

type queryResolver struct{ *Resolver }
type routeResolver struct{ *Resolver }
type serviceResolver struct{ *Resolver }
type serviceExceptionResolver struct{ *Resolver }
type stopResolver struct{ *Resolver }
type stopRouteResolver struct{ *Resolver }
type stopRouteScheduleResolver struct{ *Resolver }
type stopTimeResolver struct{ *Resolver }
type tripResolver struct{ *Resolver }
