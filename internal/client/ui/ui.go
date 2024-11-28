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
}

func NewUi() *contract.UI {
	var ui contract.UI
	app := tview.NewApplication()

	sideView := sidebar()
	reportView := battleReport(app)

	grid := tview.NewGrid().
		SetRows(0).
		SetColumns(-3, 0).
		AddItem(sideView, 0, 1, 1, 1, 0, 0, false).
		AddItem(reportView, 0, 0, 1, 1, 0, 0, false)

	app.SetRoot(grid, true).SetFocus(reportView)

	go func() {
		if err := app.Run(); err != nil {
			panic(err)
		}
	}()
	u := &UI{app: app, stopChan: make(chan bool)}
	ui = contract.UI(u)

	return &ui
}

func (u *UI) SetKeyBinding() {
	u.app.SetInputCapture(service.Controller().GetKeyBinding())
}

func (u *UI) Run(name string) {
	(*keyConsole).Info("Initializing...\n")
	go func() {
		ui := contract.UI(u)
		service.GetLocator().
			SetUI(&ui).
			SetLogger(u.Logger()).
			SetConnector(message.NewConnection(name)).
			SetController(NewController())
		service.Controller().Connect()
		u.SetKeyBinding()
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

func KeyLogger() *Console {
	return (*keyConsole).(*Console)
}
