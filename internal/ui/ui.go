package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"leveling/internal/constract"
)

type UI struct {
	app *tview.Application
}

var console *Console

func NewUi(game constract.Game) *UI {
	app := tview.NewApplication()
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			game.Stop()
			return nil
		}
		return event
	})

	textView := tview.NewTextView().
		SetDynamicColors(true)
	textView.SetChangedFunc(func() {
		app.Draw()
		textView.ScrollToEnd()
	})
	textView.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			app.Stop()
		}
	})
	textView.SetBorder(true)

	app.SetRoot(textView, true).SetFocus(textView)
	console = NewConsole(textView)

	return &UI{app}
}

func (ui *UI) Run() *UI {
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

func Logger() *Console {
	return console
}
