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

	before := time.Duration(0)
	after := time.Duration(0)
	day := time.Duration(0)

	n := time.Duration(100_000)

	for i := 0; i < int(n); i++ {
		t1 := time.Now().In(tz)
		_ = database.ScheduleIndex.Get("AK145", "49-340").Before(t1, 1)
		before = before + time.Since(t1)

		t1 = time.Now().In(tz)
		_ = database.ScheduleIndex.Get("AK145", "49-340").After(t1, 3)
		after = after + time.Since(t1)

		t1 = time.Now().In(tz)
		_ = database.ScheduleIndex.Get("AK145", "49-340").Day(t1)
		day = day + time.Since(t1)
	}

	fmt.Println(before/n, after/n, day/n)
}
