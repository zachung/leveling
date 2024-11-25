package main

import (
	"io"
	"leveling/internal/engine"
	"os"
)

func main() {
	game := (*engine.NewGame())
	stdout := io.Writer(os.Stdout)
	game.SetConsole(&stdout)
	game.Start()
}
