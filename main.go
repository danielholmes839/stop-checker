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
	now := time.Now().UTC().Truncate(time.Hour * 24)
	fmt.Println("time:", now.UTC())

	t1 := time.Now().UTC()
	results := database.ScheduleIndex.Get("AK145", "49-340")

	fmt.Println("--- next 3")
	for _, stopTime := range results.Next(t1, 3) {
		fmt.Printf("%#v\n", stopTime)
	}

	fmt.Println("--- monday")
	for _, stopTime := range results.Day(t1.Add(time.Hour * 48)) {
		fmt.Printf("%#v\n", stopTime)
	}
}
