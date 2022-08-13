package main

import (
	"fmt"
	"net/http"

	"stop-checker.com/backend"
)

func main() {
	s := backend.Server{}

	s.HandleGraphQL()

	fmt.Println("Listening...")
	err := http.ListenAndServe(":3001", nil)
	fmt.Println(err)
}
