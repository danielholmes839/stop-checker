package main

import (
	"fmt"
	"time"

	"stop-checker.com/db"
	"stop-checker.com/db/gtfs"
	"stop-checker.com/db/model"
	"stop-checker.com/travel"
)

func main() {
	dataset, _ := gtfs.NewDatasetFromFilesystem("./db/data")

	octranspo := &model.BaseParser{
		TimeZone:   dataset.TimeZone,
		TimeLayout: "15:04:05",
		DateLayout: "20060102",
	}

	base := model.NewBaseFromGTFS(dataset, octranspo)
	database := db.NewDatabase(base)

	stopRanker := db.NewStopRanker(database.StopRouteIndex)

	stops := stopRanker.Rank(
		database.StopLocationIndex.Query(model.Location{Latitude: 45.423891, Longitude: -75.6898797}, 1000),
	)

	fmt.Println("results:", len(stops))

	planner := &travel.Planner{
		ScheduleIndex:     database.ScheduleIndex,
		StopRouteIndex:    database.StopRouteIndex,
		StopTimesFromTrip: database.StopTimesFromTrip,
	}

	fmt.Println("travel planner...")

	planner.Depart(time.Now().Add(time.Hour*24), "AK151", "CD998")
}
