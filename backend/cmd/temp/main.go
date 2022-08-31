package main

import (
	"fmt"
	"time"

	"stop-checker.com/db"
)

func main() {
	database, _ := db.NewDatabaseFromFilesystem("./db/data")

	now := time.Now().Local()

	t0 := time.Now()
	scheduleResults := database.ReachIndex.ReachableWithSchedule("AK151", "49-340", now)

	fmt.Println("-------------")
	fmt.Println(time.Since(t0))

	for _, result := range scheduleResults {
		fmt.Println(result.Origin.Name, result.Departure.Format(time.Kitchen), "->", result.Destination.Name, result.Arrival.Format(time.Kitchen))
	}

	t0 = time.Now()
	results := database.ReachIndex.Reachable("AK151", "49-340", true)

	fmt.Println("-------------")
	fmt.Println(time.Since(t0))
	fmt.Println(results)

	for _, stop := range results {
		fmt.Println(stop.Name)
	}

	t0 = time.Now()
	results = database.ReachIndex.Reachable("AF940", "88-340", false)
	fmt.Println("-------------")
	fmt.Println(time.Since(t0))
	fmt.Println(results)

	for _, stop := range results {
		fmt.Println(stop.Name)
	}

	t0 = time.Now()
	scheduleResults = database.ReachIndex.ReachableWithSchedule("AF940", "88-340", now)
	fmt.Println("-------------")
	fmt.Println(time.Since(t0))

	for _, result := range scheduleResults {
		fmt.Println(result.Origin.Name, result.Departure.Format(time.Kitchen), "->", result.Destination.Name, result.Arrival.Format(time.Kitchen))
	}

}
