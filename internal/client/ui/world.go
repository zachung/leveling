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
	textView.SetChangedFunc(func(i int, s string, s2 string, r rune) {
		// send select target event to server
		selectTarget(s2)
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
	if len(heroes) == 0 {
		s.textView.Clear()
		return
	}
	// sort
	m := make(map[string]contract.Hero)
	keys := make([]string, 0, len(heroes))
	curName := service.Connector().GetCurName()
	for _, hero := range heroes {
		// ignore self
		if hero.Name == curName {
			continue
		}
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
		if m[k].Target != nil {
			if m[k].Target.Name == curName {
				// you are the target
				mainText += "⚔️"
			}
		}
		s.textView.AddItem(mainText, name, rune('1'+i), nil)
	}
	s.textView.SetCurrentItem(curInx)
	s.app.Draw()
}

func (s *World) Focus() {
	s.app.SetFocus(s.textView)
}

func (s *World) SelectNext() {
	index := s.textView.GetCurrentItem() + 1
	if index >= s.textView.GetItemCount() {
		index = 0
	}
	s.textView.SetCurrentItem(index)
}
