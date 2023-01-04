package application

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"stop-checker.com/application/resolvers"
	"stop-checker.com/application/schema"
	"stop-checker.com/application/services"
	"stop-checker.com/db/model"
	"stop-checker.com/db/repository"
)

type ServerDependencies struct {
	repository.Stops
	repository.StopRoutes
	repository.Routes
	repository.Trips
	repository.Schedules
	repository.Services
	repository.ServiceExceptions
	repository.Shapes
	repository.Reach
	repository.StopLocationSearch
	repository.StopTextSearch
	StopTimesByTrip repository.InvertedIndex[model.StopTime]
	services.TravelPlanner
	services.TravelScheduler
	services.OCTranspo
	services.StaticMapEncoder
}

type ServerConfig struct {
	CORS       bool
	Playground bool
}

type Server struct {
	config *ServerConfig
	schema graphql.ExecutableSchema
}

func NewServer(config *ServerConfig, deps *ServerDependencies) *Server {
	return &Server{
		config: config,
		schema: schema.NewExecutableSchema(schema.Config{
			Resolvers: &resolvers.Root{
				BusResolver:      &resolvers.BusResolvers{},
				LocationResolver: &resolvers.LocationResolvers{},
				QueryResolver: &resolvers.QueryResolver{
					Stops:              deps.Stops,
					StopRoutes:         deps.StopRoutes,
					StopLocationSearch: deps.StopLocationSearch,
					StopTextSearch:     deps.StopTextSearch,
					QueryTravelPlanner: &resolvers.QueryTravelPlanner{
						Planner:   deps.TravelPlanner,
						Scheduler: deps.TravelScheduler,
					},
				},
				RouteResolver:          &resolvers.RouteResolvers{},
				ScheduleResolver:       &resolvers.ScheduleResolvers{},
				ScheduleResultResolver: &resolvers.ScheduleResultResolvers{},
				ServiceResolver:        &resolvers.ServiceResolvers{},
				StopResolver: &resolvers.StopResolvers{
					StopRoutes: deps.StopRoutes,
				},
				StopRouteResolver: &resolvers.StopRouteResolvers{
					Stops:            deps.Stops,
					Routes:           deps.Routes,
					Schedules:        deps.Schedules,
					Reach:            deps.Reach,
					OCTranspo:        deps.OCTranspo,
					StaticMapEncoder: deps.StaticMapEncoder,
				},
				StopTimeResolver: &resolvers.StopTimeResolvers{
					Trips: deps.Trips,
					Stops: deps.Stops,
				},
				TransitResolver: &resolvers.TransitResolvers{
					Routes: deps.Routes,
					Trips:  deps.Trips,
					Shapes: deps.Shapes,
				},
				TravelScheduleResolver:    &resolvers.TravelScheduleResolvers{},
				TravelScheduleLegResolver: &resolvers.TravelScheduleLegResolvers{},
				TravelScheduleNodeResolver: &resolvers.TravelScheduleNodeResolvers{
					Stops: deps.Stops,
				},
				TripResolver: &resolvers.TripResolvers{
					Routes:          deps.Routes,
					Shapes:          deps.Shapes,
					Services:        deps.Services,
					StopTimesByTrip: deps.StopTimesByTrip,
				},
			},
		}),
	}
}

func (s *Server) Run(port string) {
	// https://github.com/danielholmes839/stop-checker.com-2/commit/55a7a8a3d7f46111111c1a482e1b504c8e216520#diff-2395afde28d649c46a74a24d548789a3384fbca13e5ef45e71040ac7f8e0c7cd
	r := chi.NewRouter()
	graphqlHandler := handler.NewDefaultServer(s.schema)

	r.Get("/graphql-playground", playground.Handler("stop-checker.com", "/graphql"))

	r.Post("/graphql", func(w http.ResponseWriter, r *http.Request) {
		graphqlHandler.ServeHTTP(w, r)
	})

	http.ListenAndServe(port, r)
}
