package keys

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"leveling/internal/client/service"
	contract2 "leveling/internal/contract"
)

type Rune struct {
	*T
}

func NewRune(next Func) *Rune {
	t := &T{next: next}
	i := &Rune{T: t}
	t.Func = i

	return i
}

func (c Rune) handleEvent() *ebiten.Key {
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		spell := contract2.ActionEvent{Event: contract2.Event{Type: contract2.Action}}
		spell.Id = 1
		service.Controller().Send(spell)
		key := ebiten.Key1
		return &key
	} else if inpututil.IsKeyJustPressed(ebiten.Key2) {
		spell := contract2.ActionEvent{Event: contract2.Event{Type: contract2.Action}}
		spell.Id = 2
		service.Controller().Send(spell)
		key := ebiten.Key2
		return &key
	} else if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		service.Controller().Connect("Sin")
		key := ebiten.KeyS
		return &key
	} else if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		service.Controller().Connect("Taras")
		key := ebiten.KeyT
		return &key
	} else if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		service.Controller().Connect("Brian")
		key := ebiten.KeyB
		return &key
	}
	return nil
}
