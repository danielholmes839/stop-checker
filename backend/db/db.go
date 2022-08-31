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
	timezone *time.Location

	// basic indexes
	*BaseIndex

	// inverted indexes
	StopTimesByTrip *InvertedIndex[model.StopTime]

	// specialized indexes
	*StopRouteIndex    // get routes by stop id
	*ScheduleIndex     // get schedule by stop and route id
	*StopLocationIndex // get stops by location
	*StopTextIndex     // get stops by text
	*ReachIndex
}

func NewDatabase(base *model.Base, timezone *time.Location) *Database {
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

	stopRoutesIndex := NewStopRouteIndex(baseIndex, base)
	stopTimesByTrip := NewInvertedIndex("stop time", base.StopTimes, func(record model.StopTime) (key string) {
		return record.TripId
	})
	scheduleIndex := NewScheduleIndex(baseIndex, base)

	return &Database{
		timezone: timezone,

		// inverted indexes
		StopTimesByTrip: stopTimesByTrip,

		// specialized indexes
		BaseIndex:      baseIndex,
		StopRouteIndex: stopRoutesIndex,
		ScheduleIndex:  scheduleIndex,
		StopLocationIndex: NewStopLocationIndex(baseIndex, base, ResolutionConfig{
			Level:      9,
			EdgeLength: 174.375668,
		}),
		StopTextIndex: NewStopTextIndex(base.Stops, stopRoutesIndex),
		ReachIndex:    NewReachIndex(baseIndex, base, stopTimesByTrip, scheduleIndex.indexesRequiredBySchedule),
	}
}

func (db *Database) TZ() *time.Location {
	return db.timezone
}
