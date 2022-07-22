package main

import (
	"fmt"

	"stop-checker.com/db"
	"stop-checker.com/db/gtfs"
	"stop-checker.com/db/model"
)

func main() {
	dataset, err := gtfs.NewDatasetFromFilesystem("./data")
	if err != nil {
		panic(err)
	}

	base := model.NewBaseFromGTFS(dataset, &model.BaseParser{
		TimeLayout: "15:04:05",
		DateLayout: "20060102",
	})

	fmt.Printf("%#v\n", dataset.Agencies[0])
	fmt.Printf("%#v\n", dataset.CalendarDates[0])
	fmt.Printf("%#v\n", dataset.Calendars[0])
	fmt.Printf("%#v\n", dataset.Routes[0])
	fmt.Printf("%#v\n", dataset.StopTimes[0])
	fmt.Printf("%#v\n", dataset.Stops[0])
	fmt.Printf("%#v\n", dataset.Trips[0])

	fmt.Printf("%#v\n", base.Agency)
	fmt.Printf("%#v\n", base.CalendarDates[0])
	fmt.Printf("%#v\n", base.Calendars[0])
	fmt.Printf("%#v\n", base.Routes[0])
	fmt.Printf("%#v\n", base.StopTimes[0])
	fmt.Printf("%#v\n", base.Stops[0])
	fmt.Printf("%#v\n", base.Trips[0])

	db := db.NewInMemoryDB(base)

	if stopTimes, ok := db.StopTimesByTripId.Get(base.Trips[0].ID); ok {
		fmt.Println(stopTimes)
	}
}
