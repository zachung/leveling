package ui

import (
	"fmt"
	"github.com/ebitenui/ebitenui"
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

var keyConsole *Console
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

func (u *UI) Logger() *contract.Console {
	c := contract.Console(console)
	return &c
}

func (u *UI) SideLogger() *contract.Console {
	c := contract.Console(keyConsole)
	return &c
}

func (u *UI) Run() {
	defer u.Stop()

	locator := service.GetLocator().SetBus(service.NewBus())

	console = &Console{}
	keyConsole = &Console{}

	ui := ebitenui.UI{
		Container: layoutRoot(),
	}
	game := Game{
		ui: &ui,
	}
	game.state = newState()
	game.world = newWorld()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Leveling")

	go func() {
		locator.SetLogger(u.Logger())
		service.Logger().Info("Initializing...\n")

		ui2 := contract.UI(u)
		locator.
			SetUI(&ui2).
			SetKeyLogger(u.SideLogger()).
			SetConnector(message.NewConnection()).
			SetController(NewController())
		service.Logger().Info("Ready for connect, press T/S/B start.\n")
	}()

	// run Ebiten main loop
	err := ebiten.RunGame(&game)
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
