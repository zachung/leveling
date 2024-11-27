package main

import (
	"io"
	"leveling/internal/server"
	"os"
)

func main() {
	server := (*server.NewServer())
	stdout := io.Writer(os.Stdout)
	server.SetConsole(&stdout)
	server.Start()
}
