package db

import (
	"runtime"
	"time"

	"github.com/rs/zerolog/log"
	"stop-checker.com/db/model"
	"stop-checker.com/gtfs"
)

type DB struct {
	// basic indexes
	Routes             *Index[model.Route]
	ServiceExeceptions *ServiceExceptionIndex // lookup by serviceId and time
	Services           *Index[model.Service]
	Stops              *Index[model.Stop]
	StopTimes          *Index[model.StopTime]
	Trips              *Index[model.Trip]
	Shapes             *InvertedIndex[model.Shape]

	// inverted indexes
	StopTimesByTrip *InvertedIndex[model.StopTime]

	// specialized indexes
	*StopRouteIndex    // get routes by stop id
	*ScheduleIndex     // get schedule by stop and route id
	*StopLocationIndex // get stops by location
	*StopTextIndex     // get stops by text
	*ReachIndex
}

func NewDB(dataset *model.Dataset) *DB {
	t0 := time.Now()

	// GTFS indexes
	routes := NewIndex("routes", dataset.Routes, func(route model.Route) string {
		return route.ID()
	})

	serviceExeceptions := NewServiceExceptionIndex(dataset.ServiceExceptions)

	services := NewIndex("services", dataset.Services, func(service model.Service) string {
		return service.ID()
	})
	stops := NewIndex("stops", dataset.Stops, func(stop model.Stop) string {
		return stop.ID()
	})
	stopTimes := NewIndex("stop-times", dataset.StopTimes, func(stoptime model.StopTime) string {
		return stoptime.ID()
	})
	trips := NewIndex("trips", dataset.Trips, func(trip model.Trip) string {
		return trip.ID()
	})

	stopRoutesIndex := NewStopRouteIndex(trips, dataset.StopTimes)

	stopTimesByTrip := NewInvertedIndex("stop-times-by-trip", dataset.StopTimes, func(record model.StopTime) (key string) {
		return record.TripId
	})

	shapes := NewInvertedIndex("shapes", dataset.Shapes, func(record model.Shape) (key string) {
		return record.ID()
	})

	scheduleIndex := NewScheduleIndex(dataset.StopTimes, &indexesRequiredBySchedule{
		trips:             trips,
		services:          services,
		serviceExceptions: serviceExeceptions,
	})

	stopsByCode := NewInvertedIndex("stops-by-code", dataset.Stops, func(stop model.Stop) (key string) {
		return stop.Code
	})

	database := &DB{
		// GTFS indexes
		Routes:             routes,
		ServiceExeceptions: serviceExeceptions,
		Services:           services,
		Stops:              stops,
		StopTimes:          stopTimes,
		Trips:              trips,
		Shapes:             shapes,

		StopTimesByTrip: stopTimesByTrip,

		// specialized indexes
		StopRouteIndex: stopRoutesIndex,
		ScheduleIndex:  scheduleIndex,
		StopLocationIndex: NewStopLocationIndex(dataset.Stops, Resolution{
			Level:      9,
			EdgeLength: 174.375668,
		}),
		StopTextIndex: NewStopTextIndex(stopsByCode, stopRoutesIndex, dataset.Stops),
		ReachIndex:    NewReachIndex(trips, stops, dataset.Trips, stopTimesByTrip, scheduleIndex.indexesRequiredBySchedule),
	}

	log.Info().Dur("duration", time.Since(t0)).Msg("initialized database")
	return database
}

func NewDBFromFilesystem(path string) (*DB, *model.Dataset) {
	// input
	input, err := gtfs.FileInput(path)
	if err != nil {
		panic(err)
	}

	// read the dataset
	reader := &gtfs.CSVReader{}
	raw, err := reader.ReadDataset(input)
	if err != nil {
		panic(err)
	}

	// parse the dataset
	parser := &gtfs.CSVParser{
		ParserFilter: gtfs.NewCutoffFilter(time.Now()),
		TZ:           time.Local,
		TimeLayout:   "15:04:05",
		DateLayout:   "20060102",
	}

	dataset := parser.ParseDataset(raw)

	// indexes
	database := NewDB(dataset)
	runtime.GC()

	return database, dataset
}
