package resolvers

import (
	"context"

	"stop-checker.com/application/services"
	"stop-checker.com/db/model"
	"stop-checker.com/db/repository"
	"stop-checker.com/features/staticmaps"
)

type StopRouteResolversReach interface {
	repository.Reachable
	repository.ReachableBetween
}

type StopRouteResolvers struct {
	repository.Stops
	repository.Routes
	repository.Schedules
	Reach StopRouteResolversReach
	services.OCTranspo
	services.StaticMapEncoder
}

func (r *StopRouteResolvers) Stop(ctx context.Context, obj *model.StopRoute) (model.Stop, error) {
	return r.Stops.Get(obj.StopId)
}

func (r *StopRouteResolvers) Route(ctx context.Context, obj *model.StopRoute) (model.Route, error) {
	return r.Routes.Get(obj.RouteId)
}

func (r *StopRouteResolvers) Direction(ctx context.Context, obj *model.StopRoute) (string, error) {
	return obj.DirectionId, nil
}

func (r *StopRouteResolvers) Schedule(ctx context.Context, obj *model.StopRoute) (repository.Schedule, error) {
	return r.Schedules.Get(obj.StopId, obj.RouteId), nil
}

func (r *StopRouteResolvers) ScheduleReaches(ctx context.Context, obj *model.StopRoute, destination string) (repository.Schedule, error) {
	schedule, _ := r.Reach.ReachableBetweenWithSchedule(obj.StopId, destination, obj.RouteId)
	return schedule, nil
}

func (r *StopRouteResolvers) Reaches(ctx context.Context, obj *model.StopRoute, forward bool) ([]model.Stop, error) {
	return r.Reach.Reachable(obj.StopId, obj.RouteId, !forward), nil
}

func (r *StopRouteResolvers) LiveMap(ctx context.Context, obj *model.StopRoute) (*string, error) {
	stop, _ := r.Stops.Get(obj.StopId)
	route, _ := r.Routes.Get(obj.RouteId)
	buses, _ := r.OCTranspo.StopRouteData(stop, route.Name, obj.DirectionId)

	// create the map
	m, err := staticmaps.NewStopRouteMap(800, 400, stop.Location, buses)
	if err != nil {
		return nil, nil
	}

	url := r.StaticMapEncoder.Encode(m)
	return &url, nil
}

func (r *StopRouteResolvers) LiveBuses(ctx context.Context, obj *model.StopRoute) ([]model.Bus, error) {
	stop, _ := r.Stops.Get(obj.StopId)
	route, _ := r.Routes.Get(obj.RouteId)
	buses, _ := r.OCTranspo.StopRouteData(stop, route.Name, obj.DirectionId)
	return buses, nil
}
