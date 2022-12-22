package db

import (
	"time"

	"github.com/rs/zerolog/log"
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
	ShapesByShape   *InvertedIndex[model.Shape] // weird name. but this index groups shapes with the same shape id

	// specialized indexes
	*StopRouteIndex    // get routes by stop id
	*ScheduleIndex     // get schedule by stop and route id
	*StopLocationIndex // get stops by location
	*StopTextIndex     // get stops by text
	*ReachIndex
}

func NewDatabase(base *model.Dataset, timezone *time.Location) *Database {
	t0 := time.Now()

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

	shapesByShape := NewInvertedIndex("shapes", base.Shapes, func(record model.Shape) (key string) {
		return record.ID()
	})

	scheduleIndex := NewScheduleIndex(baseIndex, base)

	database := &Database{
		timezone: timezone,

		// inverted indexes
		StopTimesByTrip: stopTimesByTrip,
		ShapesByShape:   shapesByShape,

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

	log.Info().Dur("duration", time.Since(t0)).Msg("created indexes")
	return database
}

func (db *Database) TZ() *time.Location {
	return db.timezone
}
