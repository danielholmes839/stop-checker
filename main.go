package main

import (
	"fmt"
	"time"

	"stop-checker.com/db"
	"stop-checker.com/db/gtfs"
	"stop-checker.com/db/model"
)

func main() {
	loadingDataset := time.Now()
	dataset, err := gtfs.NewDatasetFromFilesystem("./db/data")
	if err != nil {
		panic(err)
	}

	fmt.Println("loaded dataset in", time.Since(loadingDataset))

	loadingBase := time.Now()

	base := model.NewBaseFromGTFS(dataset, &model.BaseParser{
		TimeLayout: "15:04:05",
		DateLayout: "20060102",
	})

	fmt.Println("loaded base in", time.Since(loadingBase))

	database := db.NewDatabase(base)

	t1 := time.Now()
	routes := database.RouteIndex.Get("AK145")
	fmt.Println(routes)

	t2 := time.Now()
	arrivals := database.ScheduleIndex.Get("AK145", "49-340")

	for _, arrival := range arrivals {
		fmt.Println(arrival.Time.String())
	}

	t3 := time.Now()

	fmt.Println("loaded routes in:", t2.Sub(t1))
	fmt.Println("loaded stop times in:", t3.Sub(t2))
}
