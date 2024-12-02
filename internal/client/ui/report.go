package ui

import (
	"github.com/rivo/tview"
	"leveling/internal/client/contract"
)

var console *contract.Console

func battleReport(app *tview.Application) tview.Primitive {
	textView := tview.NewTextView()
	textView.SetTitle("Report").
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true)
	textView.SetDynamicColors(true)
	textView.SetChangedFunc(func() {
		app.Draw()
		textView.ScrollToEnd()
	})
	console = NewConsole(textView)

	return textView
}
