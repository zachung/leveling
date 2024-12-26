package keys

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"leveling/internal/client/service"
	contract2 "leveling/internal/contract"
)

type Action struct {
	*T
}

func NewAction(next Func) *Action {
	t := &T{next: next}
	i := &Action{T: t}
	t.Func = i

	return i
}

func (c Action) handleEvent() *ebiten.Key {
	var listenKeys = []ebiten.Key{
		ebiten.Key1,
		ebiten.Key2,
		ebiten.Key3,
		ebiten.KeyEscape,
	}
	for _, key := range listenKeys {
		if inpututil.IsKeyJustPressed(key) {
			spell := contract2.ActionEvent{Event: contract2.Event{Type: contract2.Action}}
			spell.Id = KeyMap[key]
			service.Controller().Send(spell)
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		service.Controller().Connect("Taras")
		key := ebiten.KeyT
		return &key
	} else if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		service.Controller().Connect("Brian")
		key := ebiten.KeyB
		return &key
	} else if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		service.Controller().Connect("Sin")
		key := ebiten.KeyN
		return &key
	}
	return nil
}
