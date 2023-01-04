package main

import (
	"fmt"
	"time"

	"stop-checker.com/db"
	"stop-checker.com/db/model"
	"stop-checker.com/features/osrm"
	"stop-checker.com/features/travel"
)

func main() {
	database, _ := db.NewDBFromFilesystem("./data")
	cacheData, _ := osrm.ReadCacheData("./data/300m-directions.json")

	cache := osrm.NewCache(cacheData)
	client := osrm.NewClient("http://localhost:5000")

	// kanata := model.Location{
	// 	Latitude:  45.309942,
	// 	Longitude: -75.900594,
	// }

	// kanata2 := model.Location{
	// 	Latitude:  45.347566,
	// 	Longitude: -75.923104,
	// }

	uottawa := model.Location{
		Latitude:  45.421441,
		Longitude: -75.682023,
	}

	neighborhood := model.Location{
		Latitude:  45.397279,
		Longitude: -75.636874,
	}

	// bayview := model.Location{
	// 	Latitude:  45.406019,
	// 	Longitude: -75.721874,
	// }

	t1 := time.Now()

	planner := travel.NewPlanner(database.StopLocationIndex, database.StopRouteIndex, database.ReachIndex, cache, client)

	depart, _ := time.ParseInLocation("2006-01-02T15:04:00Z", "2022-12-30T07:58:00Z", time.Local)

	arrive, _ := time.ParseInLocation("2006-01-02T15:04:00Z", "2022-12-30T08:24:00Z", time.Local)

	fmt.Println("----")
	fmt.Println(depart, arrive)

	plan, _ := planner.Depart(depart, neighborhood, uottawa)

	for _, leg := range plan.Legs {
		fmt.Printf("%#v\n", leg)
	}

	fmt.Println("----")

	plan, _ = planner.Arrive(arrive, neighborhood, uottawa)
	for _, leg := range plan.Legs {
		fmt.Printf("%#v\n", leg)
	}

	scheduler := travel.NewScheduler(client, cache, database.Stops, database.ReachIndex, database.StopTimesByTrip)

	fmt.Println("--------")

	schedule, err := scheduler.Depart(depart.Add(-5*time.Minute), plan)
	if err != nil {
		panic(err)
	}

	for _, leg := range schedule.Legs {
		fmt.Println(leg.String())
	}

	fmt.Println("--------")
	t := time.Now()
	schedule, err = scheduler.Arrive(arrive.Add(10*time.Minute), plan)
	if err != nil {
		panic(err)
	}

	for _, leg := range schedule.Legs {
		fmt.Println(leg.String())
	}

	fmt.Println(time.Since(t))

	t2 := time.Now()
	fmt.Println(t2.Sub(t1))
}
