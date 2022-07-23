package db

import (
	"fmt"
	"time"

	"stop-checker.com/db/model"
)

type Database struct {
	// basic indexes
	Routes             *Index[model.Route]
	ServiceExeceptions *Index[model.ServiceException]
	Services           *Index[model.Service]
	Stops              *Index[model.Stop]
	StopTimes          *Index[model.StopTime]
	Trips              *Index[model.Trip]

	// inverted indexes

	// specialized indexes
	*RouteIndex
	*ScheduleIndex
}

func NewDatabase(base *model.Base) *Database {
	// record start time
	now := time.Now()
	defer func() {
		fmt.Println("created database indexes in", time.Since(now))
	}()

	return &Database{
		// basic indexes
		Routes:             NewIndex(base.Routes),
		ServiceExeceptions: NewIndex(base.ServiceExceptions),
		Services:           NewIndex(base.Services),
		Stops:              NewIndex(base.Stops),
		StopTimes:          NewIndex(base.StopTimes),
		Trips:              NewIndex(base.Trips),

		// specialized indexes
		RouteIndex:    NewRouteIndex(base),
		ScheduleIndex: NewScheduleIndex(base),
	}
}
