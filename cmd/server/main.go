package main

import (
	"io"
	"leveling/internal/engine"
	"os"
)

func main() {
	server := (*engine.NewServer())
	stdout := io.Writer(os.Stdout)
	server.SetConsole(&stdout)
	server.Start()
}
