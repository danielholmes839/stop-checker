package services

import (
	"time"

	"stop-checker.com/db/model"
	"stop-checker.com/features/staticmaps"
)

type TravelPlanner interface {
	Depart(at time.Time, origin, destination model.Location) (*model.TravelPlan, error)
	Arrive(by time.Time, origin, destination model.Location) (*model.TravelPlan, error)
}

type TravelScheduler interface {
	Depart(at time.Time, plan *model.TravelPlan) (*model.TravelSchedule, error)
	Arrive(by time.Time, plan *model.TravelPlan) (*model.TravelSchedule, error)
}

type OCTranspo interface {
	StopData(stop model.Stop) (map[string][]model.Bus, error)
	StopRouteData(stop model.Stop, routeName string, routeDirection string) ([]model.Bus, error)
}

type StaticMapEncoder interface {
	Encode(m *staticmaps.Map) string
}
