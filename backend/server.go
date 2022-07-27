package backend

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"stop-checker.com/backend/graph"
	"stop-checker.com/backend/graph/generated"
	"stop-checker.com/db"
	"stop-checker.com/db/gtfs"
	"stop-checker.com/db/model"
)

type Server struct {
	Database *db.Database
}

func (s *Server) HandleGraphQL() {
	dataset, err := gtfs.NewDatasetFromFilesystem("./db/data")
	if err != nil {
		panic(err)
	}

	octranspo := &model.BaseParser{
		TimeZone:   dataset.TimeZone,
		TimeLayout: "15:04:05",
		DateLayout: "20060102",
	}

	base := model.NewBaseFromGTFS(dataset, octranspo)
	database := db.NewDatabase(base)

	resolvers := handler.NewDefaultServer(generated.NewExecutableSchema(
		generated.Config{
			Resolvers: &graph.Resolver{
				Database: database,
				Timezone: dataset.TimeZone,
			},
		},
	))

	http.Handle("/graphql", resolvers)

	http.Handle("/graphql-playground", playground.Handler("stop-checker", "/graphql"))
}
