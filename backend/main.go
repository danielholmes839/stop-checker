package main

import (
	"fmt"
	"time"

	"stop-checker.com/db"
	"stop-checker.com/travel"
	"stop-checker.com/travel/schedule"
)

func printLegs(legs []*travel.Leg) {
	for _, leg := range legs {
		fmt.Println(leg.String())
	}
}

func main() {
	database, base := db.NewDatabaseFromFilesystem("./db/data")
	scheduleIndex := schedule.NewIndex(database.BaseIndex, base)

	scheduler := travel.NewScheduler(&travel.SchedulerConfig{
		StopIndex:         database.Stops,
		StopTimesFromTrip: database.StopTimesFromTrip,
		ScheduleIndex:     scheduleIndex,
	})

	route := travel.Route{
		{Origin: "AK151", Destination: "AF920", RouteId: "49-340", Walk: false}, // arch/pleasant park -> hurdman b
		{Origin: "AF920", Destination: "AF990", Walk: true},                     // hurdman b to o train west
		{Origin: "AF990", Destination: "CD998", RouteId: "1-340", Walk: false},  // o train west to uottawa
	}

	departure, _ := time.ParseInLocation("2006-01-02 15:04", "2022-08-24 07:57", base.TZ()) // 8:12 am EST
	before, _ := time.ParseInLocation("2006-01-02 15:04", "2022-08-24 08:18", base.TZ())    // 8:28 am EST

	legs, _ := scheduler.Depart(departure, route)
	for _, leg := range legs {
		fmt.Println(leg)
	}

	fmt.Println("------------")

	legs, _ = scheduler.Arrive(before, route)
	for _, leg := range legs {
		fmt.Println(leg)
	}

	// planner := travel.NewPlanner(&travel.PlannerConfig{
	// 	ScheduleIndex:     scheduleIndex,
	// 	StopLocationIndex: database.StopLocationIndex,
	// 	StopRouteIndex:    database.StopRouteIndex,
	// 	StopIndex:         database.Stops,
	// 	StopTimesFromTrip: database.StopTimesFromTrip,
	// })

	// printPlan := func(route travel.Route) {
	// 	fmt.Println("------- NEW PLAN -------")
	// 	for _, leg := range route {
	// 		origin, _ := database.Stops.Get(leg.Origin)
	// 		destination, _ := database.Stops.Get(leg.Destination)
	// 		fmt.Println("route:", leg.RouteId, origin.Name, "->", destination.Name, "walk =", leg.Walk)
	// 	}
	// }

	// t0 = time.Now()
	// plan, _ := planner.Depart(departure, "AK151", "CD998")
	// fmt.Println(time.Since(t0))
	// printPlan(plan)

	// t0 = time.Now()
	// plan, _ = planner.Depart(departure, "AL050", "CD998")
	// fmt.Println(time.Since(t0))
	// printPlan(plan)

	// t0 = time.Now()
	// plan, _ = planner.Depart(departure, "AE640", "CD999")
	// fmt.Println(time.Since(t0))
	// printPlan(plan)
}
