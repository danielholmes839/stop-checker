package main

import (
	"fmt"
	"time"

	"stop-checker.com/db"
	"stop-checker.com/features/travel"
)

func main() {
	database, _ := db.NewDatabaseFromFilesystem("./db/data", time.Now())

	now := time.Now().Local()

	t0 := time.Now()
	scheduleResults := database.ReachIndex.ReachableForwardWithNext("AK151", "49-340", now)

	fmt.Println("-------------")
	fmt.Println(time.Since(t0))

	for _, result := range scheduleResults {
		fmt.Println(result.Origin.Name, result.Departure.Format(time.Kitchen), "->", result.Destination.Name, result.Arrival.Format(time.Kitchen))
	}

	t0 = time.Now()
	scheduleResults = database.ReachIndex.ReachableBackwardWithPrevious("AK151", "49-340", now)

	fmt.Println("-------------")
	fmt.Println(time.Since(t0))

	for _, result := range scheduleResults {
		fmt.Println(result.Origin.Name, result.Departure.Format(time.Kitchen), "->", result.Destination.Name, result.Arrival.Format(time.Kitchen))
	}

	fmt.Println("---------------------------")

	planner := travel.NewPlanner(&travel.PlannerConfig{
		StopLocationIndex: database.StopLocationIndex,
		StopRouteIndex:    database.StopRouteIndex,
		StopIndex:         database.Stops,
		ReachIndex:        database.ReachIndex,
	})

	arrive, _ := time.Parse("2006-01-02T15:04:00Z", "2022-08-25T12:18:00Z") // 8:28 am EST
	arrive = arrive.In(time.Local)

	t0 = time.Now()
	p1, err := planner.Arrive(arrive, "AK151", "CD998")
	fmt.Println(time.Since(t0))
	fmt.Println(err)
	fmt.Println(p1)

	// t0 = time.Now()
	// results := database.ReachIndex.Reachable("AK151", "49-340", true)

	// fmt.Println("-------------")
	// fmt.Println(time.Since(t0))
	// fmt.Println(results)

	// for _, stop := range results {
	// 	fmt.Println(stop.Name)
	// }

	// t0 = time.Now()
	// results = database.ReachIndex.Reachable("AF940", "88-340", false)
	// fmt.Println("-------------")
	// fmt.Println(time.Since(t0))
	// fmt.Println(results)

	// for _, stop := range results {
	// 	fmt.Println(stop.Name)
	// }

	// t0 = time.Now()
	// scheduleResults = database.ReachIndex.ReachableForwardWithNext("AF940", "88-340", now)
	// fmt.Println("-------------")
	// fmt.Println(time.Since(t0))

	// for _, result := range scheduleResults {
	// 	fmt.Println(result.Origin.Name, result.Departure.Format(time.Kitchen), "->", result.Destination.Name, result.Arrival.Format(time.Kitchen))
	// }

}
