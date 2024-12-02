package ui

import (
	"github.com/rivo/tview"
	"leveling/internal/client/contract"
	"leveling/internal/client/message"
	"leveling/internal/client/service"
)

type UI struct {
	app      *tview.Application
	stopChan chan bool
	state    *contract.State
	world    *contract.World
}

func NewUi() *contract.UI {
	var ui contract.UI
	app := tview.NewApplication()

	sideView := sidebar()
	reportView := battleReport(app)
	state := newState(app)
	world := newWorld(app)

	grid := tview.NewGrid().
		SetRows(-2, -2, -2).
		SetColumns(-3, 0).
		AddItem(state.textView, 0, 0, 1, 1, 0, 0, false).
		AddItem(world.textView, 0, 1, 2, 1, 0, 0, false).
		AddItem(sideView, 2, 1, 1, 1, 0, 0, false).
		AddItem(reportView, 1, 0, 2, 1, 0, 0, false)

	app.SetRoot(grid, true).SetFocus(reportView)

	go func() {
		if err := app.Run(); err != nil {
			panic(err)
		}
	}()
	s := contract.State(state)
	w := contract.World(world)
	u := &UI{app: app, stopChan: make(chan bool), state: &s, world: &w}
	ui = contract.UI(u)

	return &ui
}

func (u *UI) SetKeyBinding() {
	u.app.SetInputCapture(service.Controller().GetKeyBinding())
}

func (u *UI) Run() {
	locator := service.GetLocator().SetLogger(u.Logger())
	service.Logger().Info("Initializing...\n")
	go func() {
		ui := contract.UI(u)
		locator.
			SetUI(&ui).
			SetKeyLogger(u.SideLogger()).
			SetConnector(message.NewConnection()).
			SetController(NewController())
		u.SetKeyBinding()
		service.Logger().Info("Ready for connect, press T/S/B start.\n")
	}()
	<-u.stopChan
}

func (u *UI) Stop() {
	u.app.Stop()
	u.stopChan <- true
}

func (u *UI) Logger() *contract.Console {
	return console
}

func (u *UI) SideLogger() *contract.Console {
	return keyConsole
}

func (u *UI) State() contract.State {
	return *u.state
}

func (u *UI) World() contract.World {
	return *u.world
}
