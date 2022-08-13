package main

import (
	"fmt"
	"net/http"

	"stop-checker.com/server"
)

func main() {
	s := server.Server{}

	s.HandleGraphQL()

	fmt.Println("Listening...")
	err := http.ListenAndServe(":3001", nil)
	fmt.Println(err)
}
