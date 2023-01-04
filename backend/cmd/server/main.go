package main

import (
	"time"

	"stop-checker.com/application"
	"stop-checker.com/db"
	"stop-checker.com/features/octranspo"
	"stop-checker.com/features/osrm"
	"stop-checker.com/features/staticmaps"
	"stop-checker.com/features/travel"
)

func main() {
	application.ReadConfig()
	config := application.GetConfig()

	// db setup
	database, _ := db.NewDBFromFilesystem(config.DATA_GTFS)

	// osrm setup
	directionsCacheData, err := osrm.ReadCacheData(config.DATA_DIRECTIONS)
	if err != nil {
		panic(err)
	}
	directionsCache := osrm.NewCache(directionsCacheData)
	directions := osrm.NewClient(config.OSRM_ENDPOINT)

	// octranspo
	octranspoAPI := octranspo.NewAPI(time.Second*30, &octranspo.Client{
		Endpoint:          config.OCTRANSPO_ENDPOINT,
		OCTRANSPO_APP_ID:  config.OCTRANSPO_APP_ID,
		OCTRANSPO_API_KEY: config.OCTRANSPO_API_KEY,
	})

	// map encoder
	mapEncoder := &staticmaps.GoogleMapEncoder{
		Key: config.GOOGLE_MAPS_API_KEY,
	}

	// travel setup
	planner := travel.NewPlanner(
		database.StopLocationIndex,
		database.StopRouteIndex,
		database.ReachIndex,
		directionsCache,
		directions,
	)

	scheduler := travel.NewScheduler(
		directions,
		directionsCache,
		database.Stops,
		database.ReachIndex,
		database.StopTimesByTrip,
	)

	server := application.NewServer(
		&application.ServerConfig{
			CORS:       config.SERVER_ENABLE_CORS,
			Playground: config.SERVER_ENABLE_PLAYGROUND,
		},
		&application.ServerDependencies{
			Stops:              database.Stops,
			StopRoutes:         database.StopRouteIndex,
			Routes:             database.Routes,
			Trips:              database.Trips,
			Schedules:          database.ScheduleIndex,
			Services:           database.Services,
			ServiceExceptions:  database.ServiceExeceptions,
			Shapes:             database.Shapes,
			Reach:              database.ReachIndex,
			StopLocationSearch: database.StopLocationIndex,
			StopTextSearch:     database.StopTextIndex,
			StopTimesByTrip:    database.StopTimesByTrip,
			TravelPlanner:      planner,
			TravelScheduler:    scheduler,
			OCTranspo:          octranspoAPI,
			StaticMapEncoder:   mapEncoder,
		})

	server.Run(":3001")
}
