package main

import (
	"leveling/internal/client/message"
	"leveling/internal/client/ui"
	"leveling/internal/constract"
)

func main() {
	newUi := ui.NewUi()
	connection := message.NewConnection(newUi.Logger())
	u := constract.UI(newUi)
	controller := ui.NewController(&u, connection)
	newUi.SetController(controller)
	newUi.Run()
}
