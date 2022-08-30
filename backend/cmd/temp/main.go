package main

import (
	"fmt"

	"stop-checker.com/db"
)

func main() {
	database, base := db.NewDatabaseFromFilesystem("./db/data")

	//                    route     headsign direction
	routeHeadsigns := map[string]map[string]bool{}

	for _, trip := range base.Trips {
		if _, ok := routeHeadsigns[trip.RouteId]; !ok {
			routeHeadsigns[trip.RouteId] = map[string]bool{}
		}
		stoptimes, _ := database.StopTimesByTrip.Get(trip.Id)

		key := fmt.Sprintf("hash:%s:%s:%s:%d", trip.DirectionId, stoptimes[0].StopId, stoptimes[len(stoptimes)-1].StopId, len(stoptimes))
		routeHeadsigns[trip.RouteId][key] = true
	}

	counter := 0
	for route, headsigns := range routeHeadsigns {
		if len(headsigns) <= 2 {
			continue
		}
		counter++
		fmt.Println("--------------------")
		fmt.Println("route:", route)
		fmt.Printf("headsigns: %#v\n", headsigns)
	}

	fmt.Printf("%d/%d routes", counter, len(base.Routes))

}
