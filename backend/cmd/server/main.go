package main

import (
	"stop-checker.com/server"
)

func main() {
	server.ReadConfig()
	config := server.GetConfig()

	s := server.Server{}
	s.Listen(config)
}
