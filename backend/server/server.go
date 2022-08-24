package server

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"stop-checker.com/db"
	"stop-checker.com/features/travel"
	"stop-checker.com/server/graph"
	"stop-checker.com/server/graph/generated"
)

type Server struct {
	Database *db.Database
}

func (s *Server) HandleGraphQL() {
	database, base := db.NewDatabaseFromFilesystem("./db/data")

	resolvers := handler.NewDefaultServer(generated.NewExecutableSchema(
		generated.Config{
			Resolvers: &graph.Resolver{
				Database: database,
				Timezone: base.TZ(),
				Planner: travel.NewPlanner(&travel.PlannerConfig{
					ScheduleIndex:     database.ScheduleIndex,
					StopLocationIndex: database.StopLocationIndex,
					StopRouteIndex:    database.StopRouteIndex,
					StopIndex:         database.Stops,
					StopTimesFromTrip: database.StopTimesFromTrip,
				}),
				Scheduler: travel.NewScheduler(&travel.SchedulerConfig{
					StopIndex:         database.Stops,
					StopTimesFromTrip: database.StopTimesFromTrip,
					ScheduleIndex:     database.ScheduleIndex,
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
