package graph

import (
	"stop-checker.com/db"
	"stop-checker.com/features/octranspo"
	"stop-checker.com/features/travel"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Config struct {
	GOOGLE_MAPS_API_KEY string
}
type Resolver struct {
	Config
	*db.Database
	OCTranspo *octranspo.API
	Planner   travel.RoutePlanner
	Scheduler travel.SchedulePlanner
}
