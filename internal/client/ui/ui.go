package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"leveling/internal/constract"
)

type UI struct {
	app        *tview.Application
	controller *constract.Controller
	stopChan   chan bool
}

func NewUi() *UI {
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

	return &UI{
		app:      app,
		stopChan: make(chan bool),
	}
}

func (ui *UI) SetKeyBinding(keyBinding func(event *tcell.EventKey) *tcell.EventKey) {
	ui.app.SetInputCapture(keyBinding)
}

func (ui *UI) Run() {
	(*keyConsole).Info("Initializing...\n")
	(*ui.controller).Connect()
	<-ui.stopChan
}

func (ui *UI) Stop() {
	ui.app.Stop()
	ui.stopChan <- true
}

func (ui *UI) SetController(controller *constract.Controller) {
	ui.controller = controller
	ui.SetKeyBinding((*controller).GetKeyBinding())
}

func (ui *UI) Logger() *constract.Console {
	return console
}

func (ui *UI) SideLogger() *constract.Console {
	return keyConsole
}

func Logger() *Console {
	return (*console).(*Console)
}

func KeyLogger() *Console {
	return (*keyConsole).(*Console)
}
