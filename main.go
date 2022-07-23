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
		TimeLayout: "15:04:05",
		DateLayout: "20060102",
	})

	database := db.NewDatabase(base)

	route, _ := database.Routes.Get("49-340")
	fmt.Printf("%#v\n", route)

	// routes := database.RouteIndex.Get("AK145")
	now := time.Now().Truncate(time.Hour * 24)
	fmt.Println("time:", now.UTC())

	t1 := time.Now()
	arrivals := database.ScheduleIndex.
		Get("AK145", "49-340").Next(now, 20).Return()

	t2 := time.Now()

	fmt.Println(t2.Sub(t1))
	fmt.Printf("records: #%d first: %#v\n", len(arrivals), arrivals[0])

}
