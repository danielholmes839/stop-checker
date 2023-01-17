package travel

import (
	"testing"
	"time"

	"stop-checker.com/db"
	"stop-checker.com/db/model"
	"stop-checker.com/features/osrm"
)

func newTestPlanner() *Planner {
	database, _ := db.NewDBFromFilesystem("../../../data")
	cacheData, _ := osrm.ReadCacheData("../../../data/300m-directions.json")
	cache := osrm.NewCache(cacheData)
	client := osrm.NewClient("http://localhost:5000")
	return NewPlanner(database.StopLocationIndex, database.StopRouteIndex, database.ReachIndex, cache, client, &PlannerMetricsEmpty{})
}
func BenchmarkPlanner(b *testing.B) {
	planner := newTestPlanner()
	depart, _ := time.ParseInLocation("2006-01-02T15:04:00Z", "2022-12-30T12:55:00Z", time.Local)

	kanata2 := model.Location{
		Latitude:  45.347566,
		Longitude: -75.923104,
	}

	home := model.Location{
		Latitude:  45.39835798229993,
		Longitude: -75.63314640440906,
	}

	b.ResetTimer()

	b.Run("benchmark", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			planner.Depart(depart, home, kanata2)
		}
	})
}
