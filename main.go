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
	tz := dataset.TimeZone

	octranspo := &model.BaseParser{
		TimeZone:   tz,
		TimeLayout: "15:04:05",
		DateLayout: "20060102",
	}

	base := model.NewBaseFromGTFS(dataset, octranspo)

	database := db.NewDatabase(base)

	t1 := time.Now().In(tz)

	fmt.Println(t1)
	results := database.ScheduleIndex.Get("AK145", "49-340")

	fmt.Println("--- next 3")
	for _, stopTime := range results.Next(t1, 3) {
		fmt.Printf("%#v\n", stopTime)
	}

	fmt.Println("--- monday")
	for _, stopTime := range results.Day(t1) {
		fmt.Printf("%#v\n", stopTime)
	}
}
