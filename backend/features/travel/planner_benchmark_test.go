package travel

import (
	"math/rand"
	"testing"
	"time"

	"stop-checker.com/db"
)

func BenchmarkPlanComplete(b *testing.B) {
	database, base := db.NewDatabaseFromFilesystem("../db/data")
	planner := NewPlanner(&PlannerConfig{
		ScheduleIndex:     database.ScheduleIndex,
		StopLocationIndex: database.StopLocationIndex,
		StopRouteIndex:    database.StopRouteIndex,
		StopIndex:         database.Stops,
		StopTimesFromTrip: database.StopTimesByTrip,
	})

	departure, _ := time.ParseInLocation("2006-01-02 15:04", "2022-07-25 11:55", database.TZ())

	b.ResetTimer()

	// relatively short distance 1 transfer (1 bus, 1 train)
	b.Run("pleasant park -> uottawa", func(b *testing.B) {
		var r Route
		for i := 0; i < b.N; i++ {
			r, _ = planner.Depart(departure, "AK151", "CD998")
		}
		b.Log(r != nil)
	})

	// far 2 transfers (3 buses)
	b.Run("pleasant park -> carling", func(b *testing.B) {
		var r Route
		for i := 0; i < b.N; i++ {
			r, _ = planner.Depart(departure, "AK151", "NO521")
		}
		b.Log(r != nil)
	})

	// random
	b.Run("random -> random", func(b *testing.B) {
		var r Route
		for n := 0; n < b.N; n++ {
			var i, j int

			for i != j {
				i = rand.Intn(len(base.Stops))
				j = rand.Intn(len(base.Stops))
			}

			origin := base.Stops[i].ID()
			destination := base.Stops[j].ID()

			r, _ = planner.Depart(departure, origin, destination)
		}
		b.Log(r != nil)
	})
}
