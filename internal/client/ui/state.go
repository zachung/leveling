package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"leveling/internal/contract"
	"time"
)

type State struct {
	textView *tview.TextView
	app      *tview.Application
}

func newState(app *tview.Application) *State {
	textView := tview.NewTextView()
	textView.SetTitle("State").
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true)
	textView.SetChangedFunc(func() {
		app.Draw()
	})

	return &State{textView, app}
}

func (s *State) UpdateState(event contract.StateChangeEvent) {
	s.textView.SetText(fmt.Sprintf("%v: %d", event.Name, event.Health))
	if event.Damage > 0 {
		s.textView.SetBorderColor(tcell.ColorRed)
	}
	go func() {
		time.Sleep(100 * time.Millisecond)
		s.textView.SetBorderColor(tcell.ColorWhite)
		s.app.Draw()
	}()
}
