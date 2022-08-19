package server

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"stop-checker.com/db"
	"stop-checker.com/server/graph"
	"stop-checker.com/server/graph/generated"
	"stop-checker.com/travel"
	"stop-checker.com/travel/schedule"
)

type Server struct {
	Database *db.Database
}

func (s *Server) HandleGraphQL() {
	database, base := db.NewDatabaseFromFilesystem("./db/data")
	scheduleIndex := schedule.NewIndex(database.BaseIndex, base)

	resolvers := handler.NewDefaultServer(generated.NewExecutableSchema(
		generated.Config{
			Resolvers: &graph.Resolver{
				Database:      database,
				Timezone:      base.TZ(),
				ScheduleIndex: scheduleIndex,
				Planner: travel.NewPlanner(&travel.PlannerConfig{
					ScheduleIndex:     scheduleIndex,
					StopLocationIndex: database.StopLocationIndex,
					StopRouteIndex:    database.StopRouteIndex,
					StopIndex:         database.Stops,
					StopTimesFromTrip: database.StopTimesFromTrip,
				}),
				Scheduler: travel.NewScheduler(&travel.SchedulerConfig{
					StopIndex:         database.Stops,
					StopTimesFromTrip: database.StopTimesFromTrip,
					ScheduleIndex:     scheduleIndex,
				}),
			},
		},
	))

	http.Handle("/graphql", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		resolvers.ServeHTTP(w, r)
	}))

	http.Handle("/graphql-playground", playground.Handler("stop-checker", "/graphql"))
}
