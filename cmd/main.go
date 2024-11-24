package main

import "leveling/internal/engine"

func main() {
	game := engine.NewGame()
	(*game).Start()
}
