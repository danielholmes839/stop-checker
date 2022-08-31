package travel

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"stop-checker.com/db"
)

func assertEqualSchedule(t *testing.T, s1, s2 Schedule) {
	assert.Equal(t, len(s1), len(s2))

	for i := range s1 {
		assert.Equal(t, s1[i], s2[i])
	}
}

func TestPlanner(t *testing.T) {
	database, _ := db.NewDatabaseFromFilesystem("../../db/data")
	planner := NewPlanner(&PlannerConfig{
		StopLocationIndex: database.StopLocationIndex,
		StopRouteIndex:    database.StopRouteIndex,
		StopIndex:         database.Stops,
		ReachIndex:        database.ReachIndex,
	})

	scheduler := NewScheduler(&SchedulerConfig{
		StopIndex:       database.Stops,
		ReachIndex:      database.ReachIndex,
		StopTimesByTrip: database.StopTimesByTrip,
	})

	t.Run("pleasant park -> uottawa", func(t *testing.T) {
		depart, _ := time.Parse("2006-01-02T15:04:00Z", "2022-08-25T11:57:00Z") // 8:12 am EST
		depart = depart.In(time.Local)

		arrive, _ := time.Parse("2006-01-02T15:04:00Z", "2022-08-25T12:18:00Z") // 8:28 am EST
		arrive = arrive.In(time.Local)

		t.Log(depart)
		t.Log(arrive)

		p1, err := planner.Depart(depart, "AK151", "CD998")
		assert.NoError(t, err)

		assert.Equal(t, "AF920", p1[0].Destination) // hurdman B
		assert.Equal(t, "AF990", p1[1].Destination) //
		assert.Equal(t, "CD998", p1[2].Destination)

		s1, _ := scheduler.Depart(depart, p1)
		assert.NoError(t, err)
		assert.Equal(t, depart, s1.Departure())
		assert.Equal(t, arrive, s1.Arrival())

		s2, err := scheduler.Arrive(arrive, p1)
		assert.NoError(t, err)

		assertEqualSchedule(t, s1, s2)
	})

	t.Run("pleasant park -> carling", func(t *testing.T) {
		depart, _ := time.Parse("2006-01-02T15:04:00Z", "2022-08-25T12:00:00Z") // 10:00 am EST
		depart = depart.In(time.Local)

		arrive, _ := time.Parse("2006-01-02T15:04:00Z", "2022-08-25T14:00:00Z") // 10:00 am EST
		arrive = arrive.In(time.Local)

		p1, err := planner.Depart(depart, "AK151", "NO521")
		assert.NoError(t, err)

		s1, err := scheduler.Arrive(arrive, p1)
		assert.NoError(t, err)

		for _, leg := range s1 {
			t.Log(leg)
		}
	})
}
