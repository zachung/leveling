package main

import (
	"leveling/internal/client/ui"
	"time"
)

func main() {
	c := make(chan bool)
	newUi := ui.NewUi()
	newUi.Run()
	newUi.SetController(ui.NewController())
	go func() {
		for {
			time.Sleep(1 * time.Second)
		}
	}()
	<-c
	//message.Connect()
}
