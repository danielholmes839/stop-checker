package server

import (
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/zerolog/log"
	"stop-checker.com/db"
	"stop-checker.com/features/octranspo"
	"stop-checker.com/features/travel"
	"stop-checker.com/server/graph"
	"stop-checker.com/server/graph/generated"
)

type Server struct {
	Database *db.Database
}

func (s *Server) Listen(config Config) {
	t0 := time.Now()
	log.Info().Msg("server starting")

	database, _ := db.NewDatabaseFromFilesystem(config.SERVER_DATASET, time.Now())

	resolvers := handler.NewDefaultServer(generated.NewExecutableSchema(
		generated.Config{
			Resolvers: &graph.Resolver{
				Config: graph.Config{
					GOOGLE_MAPS_API_KEY: config.GOOGLE_MAPS_API_KEY,
				},
				OCTranspo: octranspo.NewAPI(time.Second*30, &octranspo.Client{
					Endpoint:          config.OCTRANSPO_ENDPOINT,
					OCTRANSPO_APP_ID:  config.OCTRANSPO_APP_ID,
					OCTRANSPO_API_KEY: config.OCTRANSPO_API_KEY,
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

	// server endpoint
	http.Handle("/graphql", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if config.SERVER_ENABLE_CORS {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		}
		resolvers.ServeHTTP(w, r)
	}))

	if config.SERVER_ENABLE_PLAYGROUND {
		http.Handle("/graphql-playground", playground.Handler("stop-checker", "/graphql"))
	}

	// server ready
	log.Info().Dur("total-duration", time.Since(t0)).Msg("server ready")
	err := http.ListenAndServe(config.SERVER_PORT, nil)

	// server shutdown
	log.Error().Err(err).Msg("server shutdown")
}
