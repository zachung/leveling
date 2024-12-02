package ui

import (
	"fmt"
	"github.com/rivo/tview"
	"leveling/internal/contract"
	"sort"
)

type World struct {
	textView *tview.TextView
	app      *tview.Application
}

func newWorld(app *tview.Application) *World {
	textView := tview.NewTextView()
	textView.SetTitle("World").
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true)
	textView.SetChangedFunc(func() {
		app.Draw()
	})

	return &World{textView, app}
}

func (s *World) UpdateWorld(event contract.WorldEvent) {
	heroes := event.Heroes
	// sort
	m := make(map[string]contract.Hero)
	keys := make([]string, 0, len(heroes))
	for _, hero := range heroes {
		m[hero.Name] = hero
		keys = append(keys, hero.Name)
	}
	sort.Strings(keys)
	// make string
	var h string
	for i, k := range keys {
		h = fmt.Sprintf("%s%d) %s(%d)\n", h, i, m[k].Name, m[k].Health)
	}
	s.textView.SetText(h)
}
