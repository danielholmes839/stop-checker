package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"stop-checker.com/db"
	"stop-checker.com/db/model"
	"stop-checker.com/features/staticmaps"
	"stop-checker.com/features/travel"
	"stop-checker.com/server/graph/generated"
	"stop-checker.com/server/graph/sdl"
)

// LastUpdatedMinutes is the resolver for the lastUpdatedMinutes field.
func (r *busResolver) LastUpdatedMinutes(ctx context.Context, obj *model.Bus) (int, error) {
	now := model.NewTimeFromDateTime(time.Now().Local())
	diff := int(model.TimeDiff(obj.LastUpdated, now).Minutes())
	return diff, nil
}

// LastUpdatedMessage is the resolver for the lastUpdatedMessage field.
func (r *busResolver) LastUpdatedMessage(ctx context.Context, obj *model.Bus) (string, error) {
	now := model.NewTimeFromDateTime(time.Now().Local())
	diff := int(model.TimeDiff(obj.LastUpdated, now).Minutes())

	if diff <= 1 {
		return "Updated just now", nil
	}

	return fmt.Sprintf("Updated %d minutes ago", diff), nil
}

// Distance is the resolver for the distance field.
func (r *busResolver) Distance(ctx context.Context, obj *model.Bus) (*float64, error) {
	if obj.Location == nil {
		return nil, nil
	}

	distance := obj.Destination.Distance(*obj.Location)

	return &distance, nil
}

// Distance is the resolver for the distance field.
func (r *locationResolver) Distance(ctx context.Context, obj *model.Location, location model.Location) (float64, error) {
	return obj.Distance(location), nil
}

// Stop is the resolver for the stop field.
func (r *queryResolver) Stop(ctx context.Context, id string) (*model.Stop, error) {
	stop, err := r.Stops.Get(id)
	if err != nil {
		return nil, err
	}
	return &stop, nil
}

// StopRoute is the resolver for the stopRoute field.
func (r *queryResolver) StopRoute(ctx context.Context, stop string, route string) (*model.StopRoute, error) {
	for _, stopRoute := range r.StopRouteIndex.Get(stop) {
		if stopRoute.RouteId == route {
			return &stopRoute, nil
		}
	}
	return nil, nil
}

// SearchStopText is the resolver for the searchStopText field.
func (r *queryResolver) SearchStopText(ctx context.Context, text string, page sdl.PageInput) (*sdl.StopSearchPayload, error) {
	results := r.StopTextIndex.Query(text)
	stops := apply(results, func(result *db.StopTextResult) *model.Stop {
		return &result.Stop
	})

	paged := Paginate(stops, page)

	payload := &sdl.StopSearchPayload{
		Page:    paged.Info(),
		Results: paged.Results(),
	}

	return payload, nil
}

// SearchStopLocation is the resolver for the searchStopLocation field.
func (r *queryResolver) SearchStopLocation(ctx context.Context, location model.Location, radius float64, page sdl.PageInput, sorted bool) (*sdl.StopSearchPayload, error) {
	results := r.StopLocationIndex.Query(location, radius)

	if sorted {
		sort.Slice(results, func(i, j int) bool {
			return results[i].Distance < results[j].Distance
		})
	}

	stops := apply(results, func(result db.StopLocationResult) *model.Stop {
		return &result.Stop
	})

	paged := Paginate(stops, page)

	return &sdl.StopSearchPayload{
		Page:    paged.Info(),
		Results: paged.Results(),
	}, nil
}

// TravelPlanner is the resolver for the travelPlanner field.
func (r *queryResolver) TravelPlanner(ctx context.Context, origin string, destination string, options sdl.TravelOptions) (*sdl.TravelSchedulePayload, error) {
	if _, err := r.Stops.Get(origin); err != nil {
		sErr := "origin not found"
		return &sdl.TravelSchedulePayload{
			Schedule: nil,
			Error:    &sErr,
		}, nil
	}

	if _, err := r.Stops.Get(destination); err != nil {
		sErr := "destination and origin must be different"
		return &sdl.TravelSchedulePayload{
			Schedule: nil,
			Error:    &sErr,
		}, nil
	}

	if origin == destination {
		sErr := "destination not found"
		return &sdl.TravelSchedulePayload{
			Schedule: nil,
			Error:    &sErr,
		}, nil
	}

	if options.Datetime == nil {
		t := time.Now().In(r.TZ())
		options.Datetime = &t
	}

	// determine the route
	route, err := sdl.PlannerWrapper(r.Planner, origin, destination, options)
	if err != nil {
		sErr := "failed to plan a route"
		return &sdl.TravelSchedulePayload{
			Schedule: nil,
			Error:    &sErr,
		}, nil
	}

	// determine the schedule
	schedule, err := sdl.ScheduleWrapper(r.Scheduler, route, options)
	if err != nil {
		sErr := "failed to plan a schedule"
		return &sdl.TravelSchedulePayload{
			Schedule: nil,
			Error:    &sErr,
		}, nil
	}

	return &sdl.TravelSchedulePayload{Schedule: schedule, Error: nil}, nil
}

// TravelPlannerFixedRoute is the resolver for the travelPlannerFixedRoute field.
func (r *queryResolver) TravelPlannerFixedRoute(ctx context.Context, input []*sdl.TravelLegInput, options sdl.TravelOptions) (*sdl.TravelSchedulePayload, error) {
	fixed := sdl.NewTravelRoute(input)
	schedule, _ := sdl.ScheduleWrapper(r.Scheduler, fixed, options)
	return &sdl.TravelSchedulePayload{Schedule: schedule, Error: nil}, nil
}

// TravelRoute is the resolver for the travelRoute field.
func (r *queryResolver) TravelRoute(ctx context.Context, input []*sdl.TravelLegInput) (*sdl.TravelRoutePayload, error) {
	fixed := sdl.NewTravelRoute(input)
	return &sdl.TravelRoutePayload{
		Route: fixed,
		Error: nil,
	}, nil
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
func (r *routeResolver) Type(ctx context.Context, obj *model.Route) (sdl.RouteType, error) {
	if obj.Type == 0 {
		return sdl.RouteTypeTrain, nil
	}

	return sdl.RouteTypeBus, nil
}

// Datetime is the resolver for the datetime field.
func (r *scheduleResultResolver) Datetime(ctx context.Context, obj *db.ScheduleResult) (*time.Time, error) {
	return &obj.Time, nil
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

// ScheduleReaches is the resolver for the scheduleReaches field.
func (r *stopRouteResolver) ScheduleReaches(ctx context.Context, obj *model.StopRoute, destination string) (*db.ScheduleResults, error) {
	schedule, _ := r.ReachableBetweenWithSchedule(obj.StopId, destination, obj.RouteId)
	return schedule, nil
}

// Reaches is the resolver for the reaches field.
func (r *stopRouteResolver) Reaches(ctx context.Context, obj *model.StopRoute, forward bool) ([]*model.Stop, error) {
	stops := r.ReachIndex.Reachable(obj.StopId, obj.RouteId, !forward)
	return ref(stops), nil
}

// LiveMap is the resolver for the liveMap field.
func (r *stopRouteResolver) LiveMap(ctx context.Context, obj *model.StopRoute) (*string, error) {
	stop, _ := r.Stops.Get(obj.StopId)
	route, _ := r.Routes.Get(obj.RouteId)

	buses, err := r.OCTranspo.StopRouteData(stop, route.Name)
	if err != nil {
		return nil, nil
	}

	m, err := staticmaps.NewStopRouteMap(800, 300, stop.Location, buses)
	if err != nil {
		return nil, nil
	}

	url := m.Encode(r.GOOGLE_MAPS_API_KEY)
	return &url, nil
}

// LiveBuses is the resolver for the liveBuses field.
func (r *stopRouteResolver) LiveBuses(ctx context.Context, obj *model.StopRoute) ([]*model.Bus, error) {
	stop, _ := r.Stops.Get(obj.StopId)
	route, _ := r.Routes.Get(obj.RouteId)

	buses, err := r.OCTranspo.StopRouteData(stop, route.Name)
	if err != nil {
		return nil, nil
	}

	return ref(buses), nil
}

// Next is the resolver for the next field.
func (r *stopRouteScheduleResolver) Next(ctx context.Context, obj *db.ScheduleResults, limit int, after time.Time) ([]*db.ScheduleResult, error) {
	t := time.Now().In(r.TZ())
	if !after.IsZero() {
		t = after
	}

	stopTimes := obj.After(t, limit)
	return ref(stopTimes), nil
}

// On is the resolver for the on field.
func (r *stopRouteScheduleResolver) On(ctx context.Context, obj *db.ScheduleResults, date time.Time) ([]*db.ScheduleResult, error) {
	stopTimes := obj.Day(date)
	return ref(stopTimes), nil
}

// Stop is the resolver for the stop field.
func (r *stopTimeResolver) Stop(ctx context.Context, obj *model.StopTime) (*model.Stop, error) {
	stop, _ := r.Stops.Get(obj.StopId)
	return &stop, nil
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

// Departure is the resolver for the departure field.
func (r *transitResolver) Departure(ctx context.Context, obj *model.Transit) (*model.StopTime, error) {
	stopTimes, _ := r.StopTimesByTrip.Get(obj.TripId)
	for _, stopTime := range stopTimes {
		if stopTime.StopId == obj.OriginId {
			return &stopTime, nil
		}
	}
	return nil, nil
}

// Arrival is the resolver for the arrival field.
func (r *transitResolver) Arrival(ctx context.Context, obj *model.Transit) (*model.StopTime, error) {
	stopTimes, _ := r.StopTimesByTrip.Get(obj.TripId)
	for _, stopTime := range stopTimes {
		if stopTime.StopId == obj.DestinationId {
			return &stopTime, nil
		}
	}
	return nil, nil
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

// StopRoute is the resolver for the stopRoute field.
func (r *travelRouteLegResolver) StopRoute(ctx context.Context, obj *travel.FixedLeg) (*model.StopRoute, error) {
	if obj.Walk {
		return nil, nil
	}

	stopRoutes := r.StopRouteIndex.Get(obj.Origin)
	for _, stopRoute := range stopRoutes {
		if stopRoute.RouteId == obj.RouteId {
			return &stopRoute, nil
		}
	}

	return nil, errors.New("invalid route for origin stop")
}

// Distance is the resolver for the distance field.
func (r *travelRouteLegResolver) Distance(ctx context.Context, obj *travel.FixedLeg) (float64, error) {
	origin, _ := r.Stops.Get(obj.Origin)
	destination, _ := r.Stops.Get(obj.Destination)
	return origin.Distance(destination.Location), nil
}

// Legs is the resolver for the legs field.
func (r *travelScheduleResolver) Legs(ctx context.Context, obj travel.Schedule) ([]*travel.Leg, error) {
	return obj, nil
}

// Origin is the resolver for the origin field.
func (r *travelScheduleResolver) Origin(ctx context.Context, obj travel.Schedule) (*model.Stop, error) {
	originId := obj[0].Origin
	origin, _ := r.Stops.Get(originId)
	return &origin, nil
}

// Destination is the resolver for the destination field.
func (r *travelScheduleResolver) Destination(ctx context.Context, obj travel.Schedule) (*model.Stop, error) {
	destinationId := obj[len(obj)-1].Destination
	destination, _ := r.Stops.Get(destinationId)
	return &destination, nil
}

// Duration is the resolver for the duration field.
func (r *travelScheduleResolver) Duration(ctx context.Context, obj travel.Schedule) (int, error) {
	duration := obj.Duration()
	return int(duration.Minutes()), nil
}

// Origin is the resolver for the origin field.
func (r *travelScheduleLegResolver) Origin(ctx context.Context, obj *travel.Leg) (*model.Stop, error) {
	stop, _ := r.Stops.Get(obj.Origin)
	return &stop, nil
}

// Destination is the resolver for the destination field.
func (r *travelScheduleLegResolver) Destination(ctx context.Context, obj *travel.Leg) (*model.Stop, error) {
	stop, _ := r.Stops.Get(obj.Destination)
	return &stop, nil
}

// Transit is the resolver for the transit field.
func (r *travelScheduleLegResolver) Transit(ctx context.Context, obj *travel.Leg) (*model.Transit, error) {
	if obj.Transit == nil {
		return nil, nil // walking
	}

	trip, err := r.Trips.Get(obj.Transit.TripId)
	if err != nil {
		panic(err)
	}

	return &model.Transit{
		RouteId:       trip.RouteId,
		TripId:        trip.Id,
		OriginId:      obj.Origin,
		DestinationId: obj.Destination,
	}, nil
}

// Shape is the resolver for the shape field.
func (r *travelScheduleLegResolver) Shape(ctx context.Context, obj *travel.Leg) ([]*model.Location, error) {
	origin, _ := r.Stops.Get(obj.Origin)
	destination, _ := r.Stops.Get(obj.Destination)

	if obj.Walk {
		return []*model.Location{
			&origin.Location,
			&destination.Location,
		}, nil
	}

	trip, _ := r.Trips.Get(obj.Transit.TripId)
	shapes, _ := r.ShapesByShape.Get(trip.ShapeId)

	originClosestShapeIndex := 0
	originClosestShapeDistance := shapes[0].Distance(origin.Location)

	destinationClosestShapeIndex := 0
	destinationClosestShapeDistance := shapes[0].Distance(destination.Location)

	for i := 1; i < len(shapes); i++ {
		originDistance := shapes[i].Distance(origin.Location)
		destinationDistance := shapes[i].Distance(destination.Location)

		if originDistance < originClosestShapeDistance {
			originClosestShapeDistance = originDistance
			originClosestShapeIndex = i
		}

		if destinationDistance < destinationClosestShapeDistance {
			destinationClosestShapeDistance = destinationDistance
			destinationClosestShapeIndex = i
		}
	}

	locations := apply(shapes, func(s model.Shape) *model.Location {
		return &s.Location
	})

	return locations[originClosestShapeIndex : destinationClosestShapeIndex+1], nil
}

// Distance is the resolver for the distance field.
func (r *travelScheduleLegResolver) Distance(ctx context.Context, obj *travel.Leg) (float64, error) {
	origin, _ := r.Stops.Get(obj.Origin)
	destination, _ := r.Stops.Get(obj.Destination)
	return origin.Distance(destination.Location), nil
}

// Duration is the resolver for the duration field.
func (r *travelScheduleLegResolver) Duration(ctx context.Context, obj *travel.Leg) (int, error) {
	return int(obj.Duration.Minutes()), nil
}

// Arrival is the resolver for the arrival field.
func (r *travelScheduleLegResolver) Arrival(ctx context.Context, obj *travel.Leg) (*time.Time, error) {
	arrival := obj.Departure.Add(obj.Duration)
	return &arrival, nil
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

// Stoptimes is the resolver for the stoptimes field.
func (r *tripResolver) Stoptimes(ctx context.Context, obj *model.Trip) ([]*model.StopTime, error) {
	stoptimes, _ := r.StopTimesByTrip.Get(obj.ID())
	return ref(stoptimes), nil
}

// Shape is the resolver for the shape field.
func (r *tripResolver) Shape(ctx context.Context, obj *model.Trip) ([]*model.Location, error) {
	panic(fmt.Errorf("not implemented: Shape - shape"))
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

// Bus returns generated.BusResolver implementation.
func (r *Resolver) Bus() generated.BusResolver { return &busResolver{r} }

// Location returns generated.LocationResolver implementation.
func (r *Resolver) Location() generated.LocationResolver { return &locationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Route returns generated.RouteResolver implementation.
func (r *Resolver) Route() generated.RouteResolver { return &routeResolver{r} }

// ScheduleResult returns generated.ScheduleResultResolver implementation.
func (r *Resolver) ScheduleResult() generated.ScheduleResultResolver {
	return &scheduleResultResolver{r}
}

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

type busResolver struct{ *Resolver }
type locationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type routeResolver struct{ *Resolver }
type scheduleResultResolver struct{ *Resolver }
type serviceResolver struct{ *Resolver }
type stopResolver struct{ *Resolver }
type stopRouteResolver struct{ *Resolver }
type stopRouteScheduleResolver struct{ *Resolver }
type stopTimeResolver struct{ *Resolver }
type transitResolver struct{ *Resolver }
type travelRouteLegResolver struct{ *Resolver }
type travelScheduleResolver struct{ *Resolver }
type travelScheduleLegResolver struct{ *Resolver }
type tripResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *travelRouteLegResolver) Route(ctx context.Context, obj *travel.FixedLeg) (*model.Route, error) {
	panic("not implemented")
}
