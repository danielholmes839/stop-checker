package main

import (
	"net/http"

	"stop-checker.com/backend"
)

func main() {
	s := backend.Server{}

	s.HandleGraphQL()

	http.ListenAndServe(":5000", nil)
}
