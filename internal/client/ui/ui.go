package ui

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"leveling/internal/client/contract"
	"leveling/internal/client/message"
	"leveling/internal/client/service"
	"log"
)

type UI struct {
	game *Game
}

func NewUi() contract.UI {
	var ui contract.UI

	u := &UI{}
	ui = contract.UI(u)

	return ui
}

func (u *UI) Run() {
	defer u.Stop()

	locator := service.GetLocator().
		SetBus(service.NewBus()).
		SetChat(&Chat{})
	service.Chat().Info("Initializing...\n")

	go func() {
		locator.
			SetUI(u).
			SetConnector(message.NewConnection()).
			SetController(NewController())
		service.Chat().Info("Ready for connect, press T/S/B start.\n")
	}()

	game := NewGame()

	ebiten.SetWindowSize(contract.ScreenWidth, contract.ScreenHeight)
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

type Chat struct {
	text string
}

func (c *Chat) Info(msg string, args ...any) {
	bus := service.EventBus()
	bus.AppendReport(fmt.Sprintf(msg, args...))
}
