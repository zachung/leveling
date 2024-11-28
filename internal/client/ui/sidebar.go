package ui

import (
	"github.com/rivo/tview"
	"leveling/internal/client/contract"
)

var keyConsole *contract.Console

func sidebar() tview.Primitive {
	textView := tview.NewTextView()
	textView.SetTitle("Client events").
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true)
	textView.SetMaxLines(10)
	textView.SetChangedFunc(func() {
		textView.ScrollToEnd()
	})
	keyConsole = NewConsole(textView)

	return textView
}
