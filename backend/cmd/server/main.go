package main

import (
	"fmt"

	"stop-checker.com/server"
)

func main() {
	server.ReadConfig()
	config := server.GetConfig()

	fmt.Println(config)

	s := server.Server{}
	s.Listen(config)
}
