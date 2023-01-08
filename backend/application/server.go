package application

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
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
	EnableCORS       bool
	EnablePlayground bool
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
	r := chi.NewRouter()

	r.Use(s.CORSMiddleware())

	graphqlHandler := handler.NewDefaultServer(s.schema)

	if s.config.EnablePlayground {
		r.Get("/graphql", playground.Handler("stop-checker.com", "/graphql"))
	}

	r.Post("/graphql", func(w http.ResponseWriter, r *http.Request) {
		graphqlHandler.ServeHTTP(w, r)
	})

	http.ListenAndServe(port, r)
}

func (s *Server) CORSMiddleware() func(h http.Handler) http.Handler {
	if s.config.EnableCORS {
		// allow all origins
		return cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
		}).Handler
	}

	// allow only stop-checker.com
	return cors.New(cors.Options{
		AllowedOrigins: []string{"stop-checker.com", "www.stop-checker.com"},
	}).Handler
}
