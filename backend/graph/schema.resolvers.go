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
	"stop-checker.com/travel"
)

// SearchStopText is the resolver for the searchStopText field.
func (r *queryResolver) SearchStopText(ctx context.Context, text string) ([]*model.Stop, error) {
	results := r.StopTextIndex.Search(text)
	stops := make([]*model.Stop, len(results))

	for i, result := range results {
		stops[i] = &result.Stop
	}
	return stops, nil
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

// TravelRoutePlanner is the resolver for the travelRoutePlanner field.
func (r *queryResolver) TravelRoutePlanner(ctx context.Context, input types.TravelRoutePlannerInput) (*types.TravelRoutePayload, error) {
	route, err := r.Planner.Depart(time.Now(), input.Origin, input.Destination)
	if err != nil {
		return nil, err
	}
	return &types.TravelRoutePayload{
		Route:  route,
		Errors: []*types.Error{},
	}, nil
}

// TravelSchedulePlanner is the resolver for the travelSchedulePlanner field.
func (r *queryResolver) TravelSchedulePlanner(ctx context.Context, input types.TravelSchedulePlannerInput) (*types.TravelSchedulePayload, error) {
	legs := travel.Route{}

	for _, leg := range input.Legs {
		legs = append(legs, &travel.FixedLeg{
			Origin:      leg.Origin,
			Destination: leg.Destination,
			RouteId:     *leg.Route,
			Walk:        leg.Route == nil,
		})
	}

	if input.Arrival != nil {
		schedule, err := r.Scheduler.Arrive(*input.Arrival, legs)
		if err != nil {
			return nil, err
		}

		return &types.TravelSchedulePayload{Schedule: schedule, Errors: []*types.Error{}}, nil
	}

	departure := time.Now().In(r.Timezone)

	if input.Departure != nil {
		departure = *input.Departure
	}

	schedule, err := r.Scheduler.Depart(departure, legs)
	if err != nil {
		return nil, err
	}

	return &types.TravelSchedulePayload{Schedule: schedule, Errors: []*types.Error{}}, nil
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
	if obj.Type == 0 {
		return types.RouteTypeTrain, nil
	}

	return types.RouteTypeBus, nil
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

// Exceptions is the resolver for the exceptions field.
func (r *serviceResolver) Exceptions(ctx context.Context, obj *model.Service) ([]*model.ServiceException, error) {
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

// On is the resolver for the on field.
func (r *stopRouteScheduleResolver) On(ctx context.Context, obj *db.ScheduleResults, date time.Time) ([]*model.StopTime, error) {
	stopTimes := obj.Day(date)
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

// Route is the resolver for the route field.
func (r *transitResolver) Route(ctx context.Context, obj *model.Transit) (*model.Route, error) {
	route, _ := r.Routes.Get(obj.RouteId)
	return &route, nil
}

// Trip is the resolver for the trip field.
func (r *transitResolver) Trip(ctx context.Context, obj *model.Transit) (*model.Trip, error) {
	trip, _ := r.Trips.Get(obj.TripId)
	return &trip, nil
}

// Legs is the resolver for the legs field.
func (r *travelRouteResolver) Legs(ctx context.Context, obj travel.Route) ([]*travel.FixedLeg, error) {
	return obj, nil
}

// Origin is the resolver for the origin field.
func (r *travelRouteLegResolver) Origin(ctx context.Context, obj *travel.FixedLeg) (*model.Stop, error) {
	origin, _ := r.Stops.Get(obj.Origin)
	return &origin, nil
}

// Destination is the resolver for the destination field.
func (r *travelRouteLegResolver) Destination(ctx context.Context, obj *travel.FixedLeg) (*model.Stop, error) {
	destination, _ := r.Stops.Get(obj.Destination)
	return &destination, nil
}

// Distance is the resolver for the distance field.
func (r *travelRouteLegResolver) Distance(ctx context.Context, obj *travel.FixedLeg) (float64, error) {
	panic(fmt.Errorf("not implemented"))
}

// Route is the resolver for the route field.
func (r *travelRouteLegResolver) Route(ctx context.Context, obj *travel.FixedLeg) (*model.Route, error) {
	route, err := r.Routes.Get(obj.RouteId)
	return &route, err
}

// Legs is the resolver for the legs field.
func (r *travelScheduleResolver) Legs(ctx context.Context, obj travel.Schedule) ([]*travel.Leg, error) {
	panic(fmt.Errorf("not implemented"))
}

// Departure is the resolver for the departure field.
func (r *travelScheduleResolver) Departure(ctx context.Context, obj travel.Schedule) (*time.Time, error) {
	panic(fmt.Errorf("not implemented"))
}

// Arrival is the resolver for the arrival field.
func (r *travelScheduleResolver) Arrival(ctx context.Context, obj travel.Schedule) (*time.Time, error) {
	panic(fmt.Errorf("not implemented"))
}

// Duration is the resolver for the duration field.
func (r *travelScheduleResolver) Duration(ctx context.Context, obj travel.Schedule) (int, error) {
	panic(fmt.Errorf("not implemented"))
}

// Origin is the resolver for the origin field.
func (r *travelScheduleLegResolver) Origin(ctx context.Context, obj *travel.Leg) (*model.Stop, error) {
	panic(fmt.Errorf("not implemented"))
}

// Destination is the resolver for the destination field.
func (r *travelScheduleLegResolver) Destination(ctx context.Context, obj *travel.Leg) (*model.Stop, error) {
	panic(fmt.Errorf("not implemented"))
}

// Distance is the resolver for the distance field.
func (r *travelScheduleLegResolver) Distance(ctx context.Context, obj *travel.Leg) (float64, error) {
	panic(fmt.Errorf("not implemented"))
}

// Transit is the resolver for the transit field.
func (r *travelScheduleLegResolver) Transit(ctx context.Context, obj *travel.Leg) (*model.Transit, error) {
	panic(fmt.Errorf("not implemented"))
}

// Duration is the resolver for the duration field.
func (r *travelScheduleLegResolver) Duration(ctx context.Context, obj *travel.Leg) (int, error) {
	panic(fmt.Errorf("not implemented"))
}

// Arrival is the resolver for the arrival field.
func (r *travelScheduleLegResolver) Arrival(ctx context.Context, obj *travel.Leg) (*time.Time, error) {
	panic(fmt.Errorf("not implemented"))
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
	fmt.Println(service, obj)
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

// Transit returns generated.TransitResolver implementation.
func (r *Resolver) Transit() generated.TransitResolver { return &transitResolver{r} }

// TravelRoute returns generated.TravelRouteResolver implementation.
func (r *Resolver) TravelRoute() generated.TravelRouteResolver { return &travelRouteResolver{r} }

// TravelRouteLeg returns generated.TravelRouteLegResolver implementation.
func (r *Resolver) TravelRouteLeg() generated.TravelRouteLegResolver {
	return &travelRouteLegResolver{r}
}

// TravelSchedule returns generated.TravelScheduleResolver implementation.
func (r *Resolver) TravelSchedule() generated.TravelScheduleResolver {
	return &travelScheduleResolver{r}
}

// TravelScheduleLeg returns generated.TravelScheduleLegResolver implementation.
func (r *Resolver) TravelScheduleLeg() generated.TravelScheduleLegResolver {
	return &travelScheduleLegResolver{r}
}

// Trip returns generated.TripResolver implementation.
func (r *Resolver) Trip() generated.TripResolver { return &tripResolver{r} }

type queryResolver struct{ *Resolver }
type routeResolver struct{ *Resolver }
type serviceResolver struct{ *Resolver }
type stopResolver struct{ *Resolver }
type stopRouteResolver struct{ *Resolver }
type stopRouteScheduleResolver struct{ *Resolver }
type stopTimeResolver struct{ *Resolver }
type transitResolver struct{ *Resolver }
type travelRouteResolver struct{ *Resolver }
type travelRouteLegResolver struct{ *Resolver }
type travelScheduleResolver struct{ *Resolver }
type travelScheduleLegResolver struct{ *Resolver }
type tripResolver struct{ *Resolver }
