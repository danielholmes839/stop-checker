package main

import (
	"fmt"
	"time"

	"stop-checker.com/db"
	"stop-checker.com/db/gtfs"
	"stop-checker.com/db/model"
	"stop-checker.com/travel"
)

func main() {
	dataset, _ := gtfs.NewDatasetFromFilesystem("./db/data")

	octranspo := &model.BaseParser{
		TimeZone:   dataset.TimeZone,
		TimeLayout: "15:04:05",
		DateLayout: "20060102",
	}

	base := model.NewBaseFromGTFS(dataset, octranspo)
	database := db.NewDatabase(base)

	stopRanker := db.NewStopRanker(database.StopRouteIndex)

	stops := stopRanker.Rank(
		database.StopLocationIndex.Query(model.Location{Latitude: 45.423891, Longitude: -75.6898797}, 1000),
	)

	fmt.Println("results:", len(stops))

	// planner := &travel.Planner{
	// 	ScheduleIndex:     database.ScheduleIndex,
	// 	StopIndex:         database.Stops,
	// 	StopLocationIndex: database.StopLocationIndex,
	// 	StopRouteIndex:    database.StopRouteIndex,
	// 	StopTimesFromTrip: database.StopTimesFromTrip,
	// }

	// fmt.Println("travel planner...")

	// t1 := time.Now()
	departure, _ := time.ParseInLocation("2006-01-02 15:04", "2022-07-25 7:55", dataset.TimeZone)
	before, _ := time.ParseInLocation("2006-01-02 15:04", "2022-07-25 8:30", dataset.TimeZone)
	// node, _ := planner.Depart(departure, "AK151", "CD998")
	// for {
	// 	stop, _ := database.Stops.Get(node.StopId)
	// 	fmt.Println("-----------------------")
	// 	fmt.Println(stop.Name)
	// 	fmt.Println(node.String())
	// 	if node.Previous == nil {
	// 		break
	// 	}

	// 	node = node.Previous.Node
	// }

	planner := travel.FixedPlanner{
		StopIndex:         database.Stops,
		StopTimesFromTrip: database.StopTimesFromTrip,
		ScheduleIndex:     database.ScheduleIndex,
	}

	t1 := time.Now()
	legs, err := planner.Depart(departure, []*travel.FixedLeg{
		// arch/pleasant park -> hurdman b
		{
			Origin:      "AK151",
			Destination: "AF920",
			RouteId:     "49-340",
			Walk:        false,
		},

		// hurdman b to o train west
		{
			Origin:      "AF920",
			Destination: "AF990",
			Walk:        true,
		},

		// o train west to uottawa
		{
			Origin:      "AF990",
			Destination: "CD998",
			RouteId:     "1-340",
			Walk:        false,
		},
	})

	fmt.Println(time.Since(t1), "depart at", err)
	for _, leg := range legs {
		fmt.Println(leg.String())
	}

	t1 = time.Now()
	legs, _ = planner.Arrive(before, []*travel.FixedLeg{
		// arch/pleasant park -> hurdman b
		{
			Origin:      "AK151",
			Destination: "AF920",
			RouteId:     "49-340",
			Walk:        false,
		},

		// hurdman b to o train west
		{
			Origin:      "AF920",
			Destination: "AF990",
			Walk:        true,
		},

		// o train west to uottawa
		{
			Origin:      "AF990",
			Destination: "CD998",
			RouteId:     "1-340",
			Walk:        false,
		},
	})

	fmt.Println(time.Since(t1), "arrive by")
	for _, leg := range legs {
		fmt.Println(leg.String())
	}
}
