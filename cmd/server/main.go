package main

import (
	"leveling/internal/server"
)

func main() {
	s := *server.NewServer()
	s.Start()
}
