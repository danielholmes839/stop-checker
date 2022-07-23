package db

import (
	"fmt"
	"time"

	"stop-checker.com/db/model"
)

type BaseIndex struct {
	Routes             *Index[model.Route]
	ServiceExeceptions *Index[model.ServiceException]
	Services           *Index[model.Service]
	Stops              *Index[model.Stop]
	StopTimes          *Index[model.StopTime]
	Trips              *Index[model.Trip]
}

type Database struct {
	// basic indexes
	*BaseIndex

	// specialized indexes
	*RouteIndex    // get routes by stop id
	*ScheduleIndex // get schedule by stop and route id
}

func NewDatabase(base *model.Base) *Database {
	// record start time
	now := time.Now()
	defer func() {
		fmt.Println("created database indexes in", time.Since(now))
	}()

	baseIndex := &BaseIndex{
		// basic indexes
		Routes:             NewIndex(base.Routes),
		ServiceExeceptions: NewIndex(base.ServiceExceptions),
		Services:           NewIndex(base.Services),
		Stops:              NewIndex(base.Stops),
		StopTimes:          NewIndex(base.StopTimes),
		Trips:              NewIndex(base.Trips),
	}

	return &Database{
		BaseIndex:     baseIndex,
		RouteIndex:    NewRouteIndex(baseIndex, base),
		ScheduleIndex: NewScheduleIndex(baseIndex, base),
	}
}
