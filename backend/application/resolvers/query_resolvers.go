package resolvers

import (
	"context"
	"sort"
	"time"

	"stop-checker.com/application/schema"
	"stop-checker.com/application/services"
	"stop-checker.com/db/model"
	"stop-checker.com/db/repository"
)

type QueryResolver struct {
	repository.Stops
	repository.StopRoutes
	repository.StopLocationSearch
	repository.StopTextSearch
	*QueryTravelPlanner
}

func (r *QueryResolver) Stop(ctx context.Context, id string) (*model.Stop, error) {
	return nullable(r.Stops.Get(id)), nil
}

func (r *QueryResolver) StopRoute(ctx context.Context, stopId string, routeId string) (*model.StopRoute, error) {
	for _, stopRoute := range r.StopRoutes.Get(stopId) {
		if stopRoute.RouteId == routeId {
			return &stopRoute, nil
		}
	}
	return nil, nil
}

func (r *QueryResolver) SearchStopText(ctx context.Context, text string, page schema.PageInput) (schema.StopSearchPayload, error) {
	results := r.StopTextSearch.Query(text)
	paged := Paginate(results, page)

	return schema.StopSearchPayload{
		Page:    paged.Info(),
		Results: paged.Results(),
	}, nil
}

func (r *QueryResolver) SearchStopLocation(ctx context.Context, location model.Location, radius float64, page schema.PageInput, sorted bool) (schema.StopSearchPayload, error) {
	if radius > 5000 {
		radius = 5000
	}

	stopsWithDistance := r.StopLocationSearch.Query(location, radius)
	stops := make([]model.Stop, len(stopsWithDistance))

	if sorted {
		sort.Slice(stopsWithDistance, func(i, j int) bool {
			return stopsWithDistance[i].Distance < stopsWithDistance[j].Distance
		})
	}

	for i, stopWithDistance := range stopsWithDistance {
		stops[i] = stopWithDistance.Stop
	}

	results := Paginate(stops, page)

	return schema.StopSearchPayload{
		Page:    results.Info(),
		Results: results.Results(),
	}, nil
}

type QueryTravelPlanner struct {
	Planner   services.TravelPlanner
	Scheduler services.TravelScheduler
}

func (r *QueryTravelPlanner) TravelPlanner(ctx context.Context, origin model.Location, destination model.Location, options schema.TravelPlannerOptions) (schema.TravelSchedulePayload, error) {
	// create the travel plan
	var plan *model.TravelPlan
	var err error

	if options.Datetime == nil {
		now := time.Now()
		options.Datetime = &now
	}

	if options.Mode == schema.ScheduleModeArriveBy {
		plan, err = r.Planner.Arrive(*options.Datetime, origin, destination)
	} else {
		plan, err = r.Planner.Depart(*options.Datetime, origin, destination)
	}

	if err != nil {
		return schema.TravelSchedulePayload{
			Schedule: nil,
			Error:    ref("failed to create a travel plan"),
		}, nil
	}

	// create the travel schedule
	return r.schedule(plan, options), nil
}

func (r *QueryTravelPlanner) TravelPlannerFixedRoute(ctx context.Context, plan model.TravelPlan, options schema.TravelPlannerOptions) (schema.TravelSchedulePayload, error) {
	if options.Datetime == nil {
		now := time.Now()
		options.Datetime = &now
	}

	return r.schedule(&plan, options), nil
}

func (r *QueryTravelPlanner) TravelPlannerFixedRoutes(ctx context.Context, plans []model.TravelPlan, options schema.TravelPlannerOptions) ([]schema.TravelSchedulePayload, error) {
	if options.Datetime == nil {
		now := time.Now()
		options.Datetime = &now
	}

	schedules := []schema.TravelSchedulePayload{}
	for _, plan := range plans {
		schedules = append(schedules, r.schedule(&plan, options))
	}
	return schedules, nil
}

func (r *QueryTravelPlanner) schedule(plan *model.TravelPlan, options schema.TravelPlannerOptions) schema.TravelSchedulePayload {
	// create the travel schedule
	var err error
	var schedule *model.TravelSchedule

	if options.Mode == schema.ScheduleModeArriveBy {
		schedule, err = r.Scheduler.Arrive(*options.Datetime, plan)
	} else {
		schedule, err = r.Scheduler.Depart(*options.Datetime, plan)
	}

	if err != nil {
		return schema.TravelSchedulePayload{
			Schedule: nil,
			Error:    ref("failed to create a travel schedule"),
		}
	}

	return schema.TravelSchedulePayload{
		Schedule: schedule,
		Error:    nil,
	}
}
