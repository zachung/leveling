package ui

import (
	"github.com/rivo/tview"
	"leveling/internal/client/contract"
)

var console *contract.Console

type Report struct {
	textView *tview.TextView
	app      *tview.Application
}

func battleReport(app *tview.Application) *Report {
	textView := tview.NewTextView()
	textView.SetTitle("Report(F1)").
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true)
	textView.SetDynamicColors(true)
	textView.SetChangedFunc(func() {
		app.Draw()
		textView.ScrollToEnd()
	})
	textView.SetInputCapture(handleReportKeys)
	console = NewConsole(textView)

	return &Report{textView: textView, app: app}
}

func (r *Report) Focus() {
	r.app.SetFocus(r.textView)
}
