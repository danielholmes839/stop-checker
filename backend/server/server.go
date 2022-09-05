package server

import (
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"stop-checker.com/db"
	"stop-checker.com/features/octranspo"
	"stop-checker.com/features/travel"
	"stop-checker.com/server/graph"
	"stop-checker.com/server/graph/generated"
)

type Server struct {
	Database *db.Database
}

func (s *Server) HandleGraphQL() {
	database, _ := db.NewDatabaseFromFilesystem("./db/data", time.Now())

	resolvers := handler.NewDefaultServer(generated.NewExecutableSchema(
		generated.Config{
			Resolvers: &graph.Resolver{
				Config: graph.Config{
					GOOGLE_MAPS_API_KEY: "AIzaSyB1ha-Cb9kOv0dPi-mBZQ4JHukDRVEJ4ME",
				},
				OCTranspo: octranspo.NewAPI(time.Second*30, &octranspo.Client{
					Endpoint:          "https://api.octranspo1.com/v2.0/GetNextTripsForStopAllRoutes",
					OCTRANSPO_APP_ID:  "13d12d72",
					OCTRANSPO_API_KEY: "508a0741b6945609192422d77f3a1da4",
				}),
				Database: database,
				Planner: travel.NewPlanner(&travel.PlannerConfig{
					StopLocationIndex: database.StopLocationIndex,
					StopRouteIndex:    database.StopRouteIndex,
					StopIndex:         database.Stops,
					ReachIndex:        database.ReachIndex,
				}),
				Scheduler: travel.NewScheduler(&travel.SchedulerConfig{
					StopIndex:       database.Stops,
					ReachIndex:      database.ReachIndex,
					StopTimesByTrip: database.StopTimesByTrip,
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
