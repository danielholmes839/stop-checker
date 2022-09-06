package main

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"stop-checker.com/server"
)

func main() {
	t0 := time.Now()
	log.Info().Msg("server starting")

	s := server.Server{}
	s.HandleGraphQL()

	log.Info().Dur("total-duration", time.Since(t0)).Msg("server listening")

	err := http.ListenAndServe(":3001", nil)

	log.Error().Err(err).Msg("server shutdown")
}
