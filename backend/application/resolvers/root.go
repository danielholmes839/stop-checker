package resolvers

import "stop-checker.com/application/schema"

type Root struct {
	schema.BusResolver
	schema.LocationResolver
	schema.QueryResolver
	schema.RouteResolver
	schema.ScheduleResolver
	schema.ScheduleResultResolver
	schema.ServiceResolver
	schema.StopResolver
	schema.StopRouteResolver
	schema.StopTimeResolver
	schema.TransitResolver
	schema.TravelScheduleResolver
	schema.TravelScheduleLegResolver
	schema.TravelScheduleNodeResolver
	schema.TripResolver
}

func (r *Root) Bus() schema.BusResolver {
	return r.BusResolver
}

func (r *Root) Location() schema.LocationResolver {
	return r.LocationResolver
}

func (r *Root) Query() schema.QueryResolver {
	return r.QueryResolver
}

func (r *Root) Route() schema.RouteResolver {
	return r.RouteResolver
}

func (r *Root) Schedule() schema.ScheduleResolver {
	return r.ScheduleResolver
}

func (r *Root) ScheduleResult() schema.ScheduleResultResolver {
	return r.ScheduleResultResolver
}

func (r *Root) Service() schema.ServiceResolver {
	return r.ServiceResolver
}

func (r *Root) Stop() schema.StopResolver {
	return r.StopResolver
}

func (r *Root) StopRoute() schema.StopRouteResolver {
	return r.StopRouteResolver
}

func (r *Root) StopTime() schema.StopTimeResolver {
	return r.StopTimeResolver
}

func (r *Root) Transit() schema.TransitResolver {
	return r.TransitResolver
}

func (r *Root) TravelSchedule() schema.TravelScheduleResolver {
	return r.TravelScheduleResolver
}

func (r *Root) TravelScheduleLeg() schema.TravelScheduleLegResolver {
	return r.TravelScheduleLegResolver
}

func (r *Root) TravelScheduleNode() schema.TravelScheduleNodeResolver {
	return r.TravelScheduleNodeResolver
}

func (r *Root) Trip() schema.TripResolver {
	return r.TripResolver
}
