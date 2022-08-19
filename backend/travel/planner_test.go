package travel

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"stop-checker.com/db"
	"stop-checker.com/travel/schedule"
)

func assertEqualSchedule(t *testing.T, s1, s2 Schedule) {
	assert.Equal(t, len(s1), len(s2))

	for i := range s1 {
		assert.Equal(t, s1[i], s2[i])
	}
}

func TestPlanner(t *testing.T) {
	database, base := db.NewDatabaseFromFilesystem("../db/data")
	scheduleIndex := schedule.NewIndex(database.BaseIndex, base)

	planner := NewPlanner(&PlannerConfig{
		ScheduleIndex:     scheduleIndex,
		StopLocationIndex: database.StopLocationIndex,
		StopRouteIndex:    database.StopRouteIndex,
		StopIndex:         database.Stops,
		StopTimesFromTrip: database.StopTimesFromTrip,
	})

	scheduler := NewScheduler(&SchedulerConfig{
		StopIndex:         database.Stops,
		StopTimesFromTrip: database.StopTimesFromTrip,
		ScheduleIndex:     scheduleIndex,
	})

	t.Run("pleasant park -> uottawa", func(t *testing.T) {
		depart, _ := time.ParseInLocation("2006-01-02 15:04", "2022-08-25 08:12", base.TZ()) // 8:12 am EST
		arrive, _ := time.ParseInLocation("2006-01-02 15:04", "2022-08-25 08:28", base.TZ()) // 8:28 am EST

		p1, err := planner.Depart(depart, "AK151", "CD998")
		assert.NoError(t, err)

		assert.Equal(t, "AF920", p1[0].Destination) // hurdman B
		assert.Equal(t, "AF990", p1[1].Destination) //
		assert.Equal(t, "CD998", p1[2].Destination)

		s1, err := scheduler.Depart(depart, p1)
		assert.NoError(t, err)

		assert.Equal(t, depart, s1.Departure())
		assert.Equal(t, arrive, s1.Arrival())

		// s2, err := scheduler.Arrive(arrive, p1)
		// assert.NoError(t, err)

		// assertEqualSchedule(t, s1, s2)
	})
}
