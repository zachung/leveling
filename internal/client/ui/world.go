package ui

import (
	"fmt"
	"github.com/rivo/tview"
	"leveling/internal/client/service"
	"leveling/internal/contract"
	"sort"
)

type World struct {
	textView *tview.List
	app      *tview.Application
}

func newWorld(app *tview.Application) *World {
	textView := tview.NewList()
	textView.SetTitle("World(F2)").
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true)
	textView.ShowSecondaryText(false)
	textView.SetSelectedFunc(func(i int, s string, s2 string, r rune) {
		service.Logger().Info("Selected %s\n", s2)
		// send select target event to server
		selectTarget(s2)
		service.UI().Report().Focus()
	})

	return &World{textView, app}
}

func selectTarget(name string) {
	event := contract.SelectTargetEvent{
		Event: contract.Event{
			Type: contract.SelectTarget,
		},
		Name: name,
	}
	service.Controller().Send(event)
}

func (s *World) UpdateWorld(event contract.WorldEvent) {
	heroes := event.Heroes
	// sort
	m := make(map[string]contract.Hero)
	keys := make([]string, 0, len(heroes))
	for _, hero := range heroes {
		// TODO: ignore self
		m[hero.Name] = hero
		keys = append(keys, hero.Name)
	}
	sort.Strings(keys)
	// 紀錄之前的選擇
	var curText string
	if s.textView.GetItemCount() > 0 {
		_, curText = s.textView.GetItemText(s.textView.GetCurrentItem())
		s.textView.Clear()
	}
	// make list
	curInx := 0
	for i, k := range keys {
		name := m[k].Name
		if curText == name {
			curInx = i
		}
		mainText := fmt.Sprintf("%s(%d)", name, m[k].Health)
		s.textView.AddItem(mainText, name, rune('1'+i), nil)
	}
	if curInx != 0 {
		s.textView.SetCurrentItem(curInx)
	}
	s.app.Draw()
}

func (s *World) Focus() {
	s.app.SetFocus(s.textView)
}

func (s *World) SelectTarget(index int) {
	s.textView.SetCurrentItem(index)
}