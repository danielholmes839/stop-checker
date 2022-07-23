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

	geoindex := db.NewStopGeoIndex(database.BaseIndex, base, db.ResolutionConfig{
		Level:      9,
		EdgeLength: 174.375668,
	})

	t1 := time.Now()
	stops := geoindex.Query(model.Location{Latitude: 45.423891, Longitude: -75.6898797}, 1000)
	dur := time.Since(t1)

	fmt.Println("results:", len(stops), "duration:", dur, "furthest:", stops[len(stops)-1].Distance)
}
