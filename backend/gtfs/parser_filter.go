package gtfs

import (
	"time"

	"stop-checker.com/db/model"
)

// ParserFilter. filter records when true
type ParserFilter interface {
	FilterService(service model.Service) bool
	FilterServiceException(serviceException model.ServiceException) bool
	FilterTrip(trip model.Trip) bool
	FilterRoute(route model.Route) bool
	FilterStopTime(stoptime model.StopTime) bool
	FilterStop(stop model.Stop) bool
	FilterShape(shape model.Shape) bool
}

type EmptyFilter struct {
}

func (e *EmptyFilter) FilterService(service model.Service) bool {
	return false
}

func (e *EmptyFilter) FilterServiceException(serviceException model.ServiceException) bool {
	return false
}

func (e *EmptyFilter) FilterTrip(trip model.Trip) bool {
	return false
}

func (e *EmptyFilter) FilterRoute(route model.Route) bool {
	return false
}

func (e *EmptyFilter) FilterStopTime(stoptime model.StopTime) bool {
	return false
}

func (e *EmptyFilter) FilterStop(stop model.Stop) bool {
	return false
}

func (e *EmptyFilter) FilterShape(shape model.Shape) bool {
	return false
}

// CutoffFilter should
type CutoffFilter struct {
	cutoff        time.Time
	validRoutes   map[string]struct{}
	validTrips    map[string]struct{}
	validServices map[string]struct{}
	validShapes   map[string]struct{}
}

func NewCutoffFilter(cutoff time.Time) *CutoffFilter {
	return &CutoffFilter{
		cutoff:        cutoff,
		validRoutes:   map[string]struct{}{},
		validTrips:    map[string]struct{}{},
		validServices: map[string]struct{}{},
		validShapes:   map[string]struct{}{},
	}
}

func (c *CutoffFilter) FilterService(service model.Service) bool {
	valid := service.End.After(c.cutoff)
	if valid {
		c.validServices[service.Id] = struct{}{}
	}
	return !valid
}

func (c *CutoffFilter) FilterServiceException(serviceException model.ServiceException) bool {
	_, ok := c.validServices[serviceException.ServiceId]
	return !ok
}

func (c *CutoffFilter) FilterTrip(trip model.Trip) bool {
	_, ok := c.validServices[trip.ServiceId]
	if ok {
		c.validTrips[trip.Id] = struct{}{}
		c.validShapes[trip.ShapeId] = struct{}{}
		c.validRoutes[trip.RouteId] = struct{}{}
	}
	return !ok
}

func (c *CutoffFilter) FilterRoute(route model.Route) bool {
	_, ok := c.validRoutes[route.Id]
	return !ok
}

func (c *CutoffFilter) FilterStopTime(stoptime model.StopTime) bool {
	_, ok := c.validTrips[stoptime.TripId]
	return !ok
}

func (c *CutoffFilter) FilterStop(stop model.Stop) bool {
	return false
}

func (c *CutoffFilter) FilterShape(shape model.Shape) bool {
	_, ok := c.validShapes[shape.Id]
	return !ok
}
