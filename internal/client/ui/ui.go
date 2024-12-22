package ui

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"leveling/internal/client/contract"
	"leveling/internal/client/message"
	"leveling/internal/client/service"
	"log"
)

const (
	screenWidth  = 1024
	screenHeight = 768
)

var console *Console

type UI struct {
	game *Game
}

func NewUi() contract.UI {
	var ui contract.UI

	u := &UI{}
	ui = contract.UI(u)

	return ui
}

func (u *UI) Chat() *contract.Chat {
	c := contract.Chat(console)
	return &c
}

func (u *UI) Run() {
	defer u.Stop()

	locator := service.GetLocator().SetBus(service.NewBus())

	go func() {
		locator.SetChat(u.Chat())
		service.Chat().Info("Initializing...\n")

		ui := contract.UI(u)
		locator.
			SetUI(&ui).
			SetConnector(message.NewConnection()).
			SetController(NewController())
		service.Chat().Info("Ready for connect, press T/S/B start.\n")
	}()

	game := NewGame()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Leveling")

	// run Ebiten main loop
	err := ebiten.RunGame(game)
	if err != nil {
		log.Println(err)
	}
}

func (u *UI) Stop() {
	log.Println("Stopping...")
	service.Connector().Close()
	log.Println("bye bye")
}

func (u *UI) State() contract.State {
	return contract.State(u.game.state)
}

func (u *UI) World() contract.World {
	return contract.World(u.game.world)
}

func (u *UI) Report() contract.Panel {
	return nil
}

type Console struct {
	text string
}

func (c *Console) Info(msg string, args ...any) {
	bus := service.EventBus()
	bus.AppendReport(fmt.Sprintf(msg, args...))
}
