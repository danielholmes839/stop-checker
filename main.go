package main

import (
	"fmt"
	"time"

	"stop-checker.com/db"
	"stop-checker.com/db/gtfs"
	"stop-checker.com/db/model"
	"stop-checker.com/travel"
)

func printLegs(legs []*travel.Leg) {
	for _, leg := range legs {
		fmt.Println(leg.String())
	}
}

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

	fmt.Println("ranked:", len(stops))

	scheduler := travel.NewScheduler(&travel.SchedulerConfig{
		StopIndex:         database.Stops,
		StopTimesFromTrip: database.StopTimesFromTrip,
		ScheduleIndex:     database.ScheduleIndex,
	})

	t0 := time.Now()
	route := travel.Route{
		{Origin: "AK151", Destination: "AF920", RouteId: "49-340", Walk: false}, // arch/pleasant park -> hurdman b
		{Origin: "AF920", Destination: "AF990", Walk: true},                     // hurdman b to o train west
		{Origin: "AF990", Destination: "CD998", RouteId: "1-340", Walk: false},  // o train west to uottawa
	}

	departure, _ := time.ParseInLocation("2006-01-02 15:04", "2022-07-25 7:55", dataset.TimeZone)
	before, _ := time.ParseInLocation("2006-01-02 15:04", "2022-07-25 8:30", dataset.TimeZone)

	if legs, err := scheduler.Depart(departure, route); err == nil {
		fmt.Println(time.Since(t0))
		printLegs(legs)
	}

	fmt.Println("---")

	if legs, err := scheduler.Arrive(before, route); err == nil {
		printLegs(legs)
	}

	planner := travel.NewPlanner(&travel.PlannerConfig{
		ScheduleIndex:     database.ScheduleIndex,
		StopLocationIndex: database.StopLocationIndex,
		StopRouteIndex:    database.StopRouteIndex,
		StopIndex:         database.Stops,
		StopTimesFromTrip: database.StopTimesFromTrip,
	})

	printPlan := func(route travel.Route) {
		fmt.Println("------- NEW PLAN -------")
		for _, leg := range route {
			origin, _ := database.Stops.Get(leg.Origin)
			destination, _ := database.Stops.Get(leg.Destination)
			fmt.Println("route:", leg.RouteId, origin.Name, "->", destination.Name, "walk =", leg.Walk)
		}
	}

	t0 = time.Now()
	plan, _ := planner.Depart(departure, "AK151", "CD998")
	fmt.Println(time.Since(t0))
	printPlan(plan)

	t0 = time.Now()
	plan, _ = planner.Depart(departure, "AL050", "CD998")
	fmt.Println(time.Since(t0))
	printPlan(plan)

	t0 = time.Now()
	plan, _ = planner.Depart(departure, "AE640", "CD999")
	fmt.Println(time.Since(t0))
	printPlan(plan)
}
