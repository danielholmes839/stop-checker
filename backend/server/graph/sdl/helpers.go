package sdl

import "stop-checker.com/features/travel"

func NewTravelRoute(route []*TravelLegInput) travel.Route {
	travelRoute := travel.Route{}
	for _, r := range route {
		routeId := ""
		if r.Route != nil {
			routeId = *r.Route
		}

		travelRoute = append(travelRoute, &travel.FixedLeg{
			Origin:      r.Origin,
			Destination: r.Destination,
			RouteId:     routeId,
			Walk:        r.Route == nil,
		})
	}
	return travelRoute
}

func PlannerWrapper(planner travel.RoutePlanner, origin, destination string, options TravelOptions) (travel.Route, error) {
	if options.Mode == ScheduleModeArriveBy {
		return planner.Arrive(*options.Datetime, origin, destination)
	}
	return planner.Depart(*options.Datetime, origin, destination)
}

func ScheduleWrapper(scheduler travel.SchedulePlanner, route travel.Route, options TravelOptions) (travel.Schedule, error) {
	if options.Mode == ScheduleModeArriveBy {
		return scheduler.Arrive(*options.Datetime, route)
	}
	return scheduler.Depart(*options.Datetime, route)
}
