package main

import (
	"fmt"
	"time"

	"stop-checker.com/db"
)

func main() {
	database, _ := db.NewDatabaseFromFilesystem("./db/data")

	t0 := time.Now()
	scheduleResults := database.ReachIndex.ReachableDepartAt("AK151", "49-340", time.Now().Local())

	fmt.Println("-------------")
	fmt.Println(time.Since(t0))
	// fmt.Println(results)

	for _, result := range scheduleResults {
		fmt.Println(result.Time, result.Name)
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
	scheduleResults = database.ReachIndex.ReachableDepartAt("AF940", "88-340", time.Now().Local())
	fmt.Println("-------------")
	fmt.Println(time.Since(t0))
	fmt.Println(results)

	for _, result := range scheduleResults {
		fmt.Println(result.Time, result.Name)
	}
}
