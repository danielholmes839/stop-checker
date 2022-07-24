package main

import (
	"fmt"
	"time"

	"stop-checker.com/db"
	"stop-checker.com/db/gtfs"
	"stop-checker.com/db/model"
)

func main() {
	dataset, _ := gtfs.NewDatasetFromFilesystem("./db/data")

	octranspo := &model.BaseParser{
		TimeZone:        dataset.TimeZone,
		TimeLayout:      "15:04:05",
		DateLayout:      "20060102",
	}

	base := model.NewBaseFromGTFS(dataset, octranspo)
	database := db.NewDatabase(base)

	stopRanker := db.NewStopRanker(database.StopRouteIndex)

	t1 := time.Now()

	stops := stopRanker.Rank(
		database.StopLocationIndex.Query(model.Location{Latitude: 45.423891, Longitude: -75.6898797}, 1000),
	)

	dur := time.Since(t1)

	if len(stops) > 0 {
		fmt.Println("results:", len(stops), "duration:", dur, "furthest:", stops[len(stops)-1].Distance)
	}
}
