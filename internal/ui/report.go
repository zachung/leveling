package ui

import (
	"github.com/rivo/tview"
	"leveling/internal/constract"
)

var console *constract.Console

func battleReport(app *tview.Application) tview.Primitive {
	textView := tview.NewTextView().
		SetDynamicColors(true)
	textView.SetChangedFunc(func() {
		app.Draw()
		textView.ScrollToEnd()
	})
	console = NewConsole(textView)

	return textView
}
