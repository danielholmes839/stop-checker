package main

import (
	"fmt"

	"stop-checker.com/db"
	"stop-checker.com/db/gtfs"
	"stop-checker.com/db/model"
)

func main() {
	dataset, _ := gtfs.NewDatasetFromFilesystem("./db/data")

	base := model.NewBaseFromGTFS(dataset, &model.BaseParser{
		TimeLayout: "15:04:05",
		DateLayout: "20060102",
	})

	database := db.NewDatabase(base)

	route, _ := database.Routes.Get("49-340")
	fmt.Printf("%#v\n", route)

	// routes := database.RouteIndex.Get("AK145")
	// arrivals := database.ScheduleIndex.Get("AK145", "49-340")

}
