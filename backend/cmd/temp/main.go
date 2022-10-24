package main

import (
	"time"

	"stop-checker.com/db"
	"stop-checker.com/db/model"
	v2 "stop-checker.com/features/travel/v2"
)

func main() {
	database, _ := db.NewDatabaseFromFilesystem("./data", time.Now())

	planner := v2.NewPlanner(database)

	origin := model.Location{Latitude: 45.397671, Longitude: -75.631022}
	origin2 := model.Location{Latitude: 45.395909, Longitude: -75.631384}
	destination := model.Location{Latitude: 45.421274, Longitude: -75.681846}

	depart, _ := time.Parse("2006-01-02T15:04:00Z", "2022-10-22T11:50:00Z")
	depart = depart.In(time.Local)

	depart2, _ := time.Parse("2006-01-02T15:04:00Z", "2022-10-24T11:50:00Z")
	depart2 = depart2.In(time.Local)

	planner.Depart(origin, destination, depart)

	planner.Depart(origin2, destination, depart2)
}
