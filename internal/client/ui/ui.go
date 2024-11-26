package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	keys2 "leveling/internal/client/ui/keys"
	"leveling/internal/constract"
)

type UI struct {
	app        *tview.Application
	controller *constract.Controller
}

func NewUi() *UI {
	app := tview.NewApplication()

	sideView := sidebar()
	reportView := battleReport(app)

	grid := tview.NewGrid().
		SetRows(0).
		SetColumns(-3, 0).
		SetBorders(true).
		AddItem(sideView, 0, 1, 1, 1, 0, 0, false).
		AddItem(reportView, 0, 0, 1, 1, 0, 0, false)

	app.SetRoot(grid, true).SetFocus(reportView)

	return &UI{app: app}
}

func (ui *UI) keyBinding() {
	ui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// chain of responsibility
		keyHandlers := keys2.NewCtrlC(keys2.NewRune(nil))

		if (*keyHandlers).Execute(ui.controller, event) == nil {
			return nil
		}
		return event
	})
}

func (ui *UI) Run() *UI {
	ui.keyBinding()

	go func() {
		if err := ui.app.Run(); err != nil {
			panic(err)
		}
	}()

	return ui
}

func (ui *UI) Stop() {
	ui.app.Stop()
}

func (ui *UI) SetController(controller *constract.Controller) {
	ui.controller = controller
}

func (ui *UI) Logger() *constract.Console {
	return console
}

func (ui *UI) SideLogger() *constract.Console {
	return keyConsole
}
