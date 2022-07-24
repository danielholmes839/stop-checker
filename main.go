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

	base := model.NewBaseFromGTFS(dataset, &model.BaseParser{
		TimeZone:   dataset.TimeZone,
		TimeLayout: "15:04:05",
		DateLayout: "20060102",
	})

	database := db.NewDatabase(base)

	stopLocationIndex := db.NewStopLocationIndex(database.BaseIndex, base, db.ResolutionConfig{
		Level:      9,
		EdgeLength: 174.375668,
	})

	stopRanker := db.NewStopRanker(database.StopRouteIndex)

	t1 := time.Now()

	stops := stopRanker.Rank(
		stopLocationIndex.Query(model.Location{Latitude: 45.423891, Longitude: -75.6898797}, 1000),
		// stopLocationIndex.Query(model.Location{Latitude: 45.3957197, Longitude: -75.632076}, 1000),
	)

	dur := time.Since(t1)

	for _, stop := range stops {
		fmt.Println("---", stop.Name, "---")
		for _, route := range database.StopRouteIndex.Get(stop.Id) {
			fmt.Println(" -", route.Name, route.DirectionId)
		}
	}

	// for _, stop := range stops {
	// 	fmt.Printf("%#v\n", stop)
	// }
	fmt.Println("results:", len(stops), "duration:", dur, "furthest:", stops[len(stops)-1].Distance)
}
