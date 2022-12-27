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
	Next(t time.Time)
	Previous(t time.Time)
	After(t time.Time, limit int)
	Before(t time.Time, limit int)
}

type Schedules interface {
	Get(stopId, routeId string) (Schedule, error)
}

type StopLocationSearch interface {
	Query(radius float64) []model.StopWithDistance
}

type StopTextSearch interface {
	Query(search string) []model.Stop
}

type Reach interface {
	
}