package repository

import (
	"time"

	"stop-checker.com/db/model"
)

type InvertedIndex[R any] interface {
	Get(key string) ([]R, error)
}

type Routes interface {
	Get(routeId string) (model.Route, error)
}

type Services interface {
	Get(serviceId string) (model.Service, error)
}

type ServiceExceptions interface {
	Get(serviceId string, date time.Time) (model.ServiceException, error)
}

type Stops interface {
	Get(stopId string) (model.Stop, error)
}

type StopTimes interface {
	Get(stopId string) (model.Stop, error)
}

type Trips interface {
	Get(tripId string) (model.Trip, error)
}

type Shapes interface {
	Get(shapeId string) ([]model.Shape, error)
}

type StopRoutes interface {
	Get(stopId string) []model.StopRoute
}

type Schedule interface {
	Next(t time.Time) (model.ScheduleResult, error)
	Previous(t time.Time) (model.ScheduleResult, error)
	After(t time.Time, limit int) []model.ScheduleResult
	Before(t time.Time, limit int) []model.ScheduleResult
	Day(on time.Time) []model.ScheduleResult
}

type Schedules interface {
	Get(stopId, routeId string) (Schedule, error)
}

type StopLocationSearch interface {
	Query(origin model.Location, radius float64) []model.StopWithDistance
}

type StopTextSearch interface {
	Query(search string) []model.Stop
}

type Reach interface {
	ReachableForwardWithNext(originId, routeId string, after time.Time) []model.ReachableSchedule
	ReachableBackwardWithPrevious(originId, routeId string, after time.Time) []model.ReachableSchedule
}

type ReachBetween interface {
	ReachableBetweenWithSchedule(originId, destinationId, routeId string) (Schedule, Schedule)
}
