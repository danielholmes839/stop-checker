package db

import (
	"fmt"
	"time"

	"stop-checker.com/db/model"
)

type BaseIndex struct {
	Routes             *Index[model.Route]
	ServiceExeceptions *ServiceExceptionIndex // lookup by serviceId and time
	Services           *Index[model.Service]
	Stops              *Index[model.Stop]
	StopTimes          *Index[model.StopTime]
	Trips              *Index[model.Trip]
}

type Database struct {
	// basic indexes
	*BaseIndex

	// inverted indexes
	StopTimesFromTrip *InvertedIndex[model.StopTime]

	// specialized indexes
	*StopRouteIndex    // get routes by stop id
	*ScheduleIndex     // get schedule by stop and route id
	*StopLocationIndex // get stops by location
	*StopTextIndex     // get stops by text
}

func NewDatabase(base *model.Base) *Database {
	// record start time
	now := time.Now()
	defer func() {
		fmt.Println("created database indexes in", time.Since(now))
	}()

	baseIndex := &BaseIndex{
		// basic indexes
		Routes:             NewIndex("route", base.Routes),
		ServiceExeceptions: NewServiceExceptionIndex(base.ServiceExceptions),
		Services:           NewIndex("service", base.Services),
		Stops:              NewIndex("stop", base.Stops),
		StopTimes:          NewIndex("stop time", base.StopTimes),
		Trips:              NewIndex("trip", base.Trips),
	}

	return &Database{
		// inverted indexes
		StopTimesFromTrip: NewInvertedIndex("stop time", base.StopTimes, func(record model.StopTime) (key string) {
			return record.TripId
		}),

		// specialized indexes
		BaseIndex:      baseIndex,
		StopRouteIndex: NewStopRouteIndex(baseIndex, base),
		ScheduleIndex:  NewScheduleIndex(baseIndex, base),
		StopLocationIndex: NewStopLocationIndex(baseIndex, base, ResolutionConfig{
			Level:      9,
			EdgeLength: 174.375668,
		}),
		StopTextIndex: NewStopTextIndex(base.Stops),
	}
}
