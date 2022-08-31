package main

import (
	"fmt"
	"time"

	"stop-checker.com/db"
)

func main() {
	database, _ := db.NewDatabaseFromFilesystem("./db/data")

	t0 := time.Now()
	results := database.ReachIndex.Reachable("AK151", "49-340", false)

	fmt.Println("-------------")
	fmt.Println(time.Since(t0))
	fmt.Println(results)

	for _, stop := range results {
		fmt.Println(stop.Name)
	}

	t0 = time.Now()
	results = database.ReachIndex.Reachable("AK151", "49-340", true)

	fmt.Println("-------------")
	fmt.Println(time.Since(t0))
	fmt.Println(results)

	for _, stop := range results {
		fmt.Println(stop.Name)
	}

	t0 = time.Now()
	results = database.ReachIndex.Reachable("AF940", "88-340", false)
	fmt.Println(time.Since(t0))
	fmt.Println("-------------")
	fmt.Println(results)

	for _, stop := range results {
		fmt.Println(stop.Name)
	}
}
