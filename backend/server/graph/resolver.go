package graph

import (
	"stop-checker.com/db"
	"stop-checker.com/features/travel"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	*db.Database
	Planner   travel.RoutePlanner
	Scheduler travel.SchedulePlanner
}
