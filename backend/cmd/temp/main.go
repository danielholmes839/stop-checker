package main

import (
	"fmt"
	"time"

	"stop-checker.com/db"
	"stop-checker.com/db/model"
	v2 "stop-checker.com/features/travel/v2"
)

func main() {
	database, _ := db.NewDatabaseFromFilesystem("./data", time.Now())

	planner := v2.NewRoutePlanner(database)

	origin := model.Location{Latitude: 45.397289, Longitude: -75.631956}
	destination := model.Location{Latitude: 45.421069, Longitude: -75.682363}

	depart, _ := time.Parse("2006-01-02T15:04:00Z", "2022-10-24T11:50:00Z")
	depart = depart.In(time.Local)

	arrive, _ := time.Parse("2006-01-02T15:04:00Z", "2022-10-24T12:24:00Z")
	arrive = arrive.In(time.Local)

	planner.Depart(origin, destination, depart)

	fmt.Println("--------------")
	planner.Depart(origin, destination, depart)

	fmt.Println("--------------")
	planner.Arrive(origin, destination, arrive)
}
