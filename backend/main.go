package main

import (
	"fmt"
	"time"

	"stop-checker.com/db"
	"stop-checker.com/db/model"
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

	stopRanker := db.NewStopRanker(database.StopRouteIndex)

	stops := stopRanker.Rank(
		database.StopLocationIndex.Query(model.Location{Latitude: 45.423891, Longitude: -75.6898797}, 1000),
	)

	fmt.Println("ranked:", len(stops))

	scheduleIndex := schedule.NewIndex(database.BaseIndex, base)

	scheduler := travel.NewScheduler(&travel.SchedulerConfig{
		StopIndex:         database.Stops,
		StopTimesFromTrip: database.StopTimesFromTrip,
		ScheduleIndex:     scheduleIndex,
	})

	t0 := time.Now()
	route := travel.Route{
		{Origin: "AK151", Destination: "AF920", RouteId: "49-340", Walk: false}, // arch/pleasant park -> hurdman b
		{Origin: "AF920", Destination: "AF990", Walk: true},                     // hurdman b to o train west
		{Origin: "AF990", Destination: "CD998", RouteId: "1-340", Walk: false},  // o train west to uottawa
	}

	departure, _ := time.ParseInLocation("2006-01-02 15:04", "2022-07-25 7:55", base.TZ())
	before, _ := time.ParseInLocation("2006-01-02 15:04", "2022-07-25 8:30", base.TZ())

	if legs, err := scheduler.Depart(departure, route); err == nil {
		fmt.Println(time.Since(t0))
		printLegs(legs)
	}

	fmt.Println("---")

	if legs, err := scheduler.Arrive(before, route); err == nil {
		printLegs(legs)
	}

	planner := travel.NewPlanner(&travel.PlannerConfig{
		ScheduleIndex:     scheduleIndex,
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

	tripsByRoute := db.NewInvertedIndex("trips by route", base.Trips, func(record model.Trip) (key string) {
		return record.RouteId
	})

	for _, route := range base.Routes {
		trips, _ := tripsByRoute.Get(route.ID())
		lengths := map[int]model.Trip{}

		for _, trip := range trips {
			stoptimes, _ := database.StopTimesFromTrip.Get(trip.Id)
			lengths[len(stoptimes)] = trip
		}

		if len(lengths) > 2 {
			fmt.Println("------------")
			fmt.Println(route.ID())
			for length, trip := range lengths {
				fmt.Println(trip.Headsign, trip.DirectionId, length)
			}
		}

	}

}
