package ui

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"leveling/internal/client/service"
	"leveling/internal/client/ui/component"
	"leveling/internal/contract"
	"sort"
)

type World struct {
	event         contract.WorldEvent
	currentTarget string
}

func newWorld() *World {
	return &World{}
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

func (w *World) UpdateWorld(event contract.WorldEvent) {
	w.event = event
}

// Focus deprecated
func (w *World) Focus() {
}

func (w *World) SelectNext() {
	if len(w.event.Heroes) == 0 {
		return
	}
	curIndex := 0
	for i, hero := range w.event.Heroes {
		if hero.Name == w.currentTarget {
			curIndex = i
		}
	}
	index := curIndex + 1
	if index >= len(w.event.Heroes) {
		index = 0
	}
	selectTarget(w.event.Heroes[index].Name)
}

func (w *World) Draw(dst *ebiten.Image) {
	event := service.EventBus().GetWorldState()
	heroes := event.Heroes
	if len(heroes) == 0 {
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
	// make list
	listBox := component.NewListBox(400, 200)
	for _, k := range keys {
		name := m[k].Name
		mainText := fmt.Sprintf("%s(%d)", name, m[k].Health)
		if m[k].Target != nil {
			if m[k].Target.Name == curName {
				// you are the target
				mainText += "⚔️"
			}
		}
		listBox.AppendItem(mainText)
	}
	listBox.Draw(dst)
}
