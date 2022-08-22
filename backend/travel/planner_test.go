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
		departure, _ := time.ParseInLocation("2006-01-02 15:04", "2022-08-24 07:57", base.TZ()) // 8:12 am EST
		arrive, _ := time.ParseInLocation("2006-01-02 15:04", "2022-08-24 08:18", base.TZ())    // 8:28 am EST

		p1, err := planner.Depart(departure, "AK151", "CD998")
		assert.NoError(t, err)
		
		assert.Equal(t, Route{
			&FixedLeg{Origin: "AK151", Destination: "AF920", RouteId: "49-340", Walk: false}, // arch/pleasant park -> hurdman b
			&FixedLeg{Origin: "AF920", Destination: "AF990", Walk: true},                     // hurdman b to o train west
			&FixedLeg{Origin: "AF990", Destination: "CD998", RouteId: "1-340", Walk: false},  // o train west to uottawa
		}, p1)

		s1, err := scheduler.Depart(departure, p1)
		assert.NoError(t, err)

		assert.Equal(t, departure, s1.Departure())
		assert.Equal(t, arrive, s1.Arrival())

		s2, err := scheduler.Arrive(arrive, p1)
		assert.NoError(t, err)

		assertEqualSchedule(t, s1, s2)
	})
}
