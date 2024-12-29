package keys

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/math/f64"
	"leveling/internal/client/service"
	contract2 "leveling/internal/contract"
)

type Move struct {
	*T
}

var curVec f64.Vec2

func NewMove(next Func) *Move {
	t := &T{next: next}
	i := &Move{T: t}
	t.Func = i

	return i
}

func (m Move) handleEvent() *ebiten.Key {
	var newKeys []ebiten.Key
	newKeys = inpututil.AppendPressedKeys(newKeys[:0])
	vec := f64.Vec2{}
	for _, key := range newKeys {
		switch key {
		case ebiten.KeyW:
			vec[1] -= 1
		case ebiten.KeyS:
			vec[1] += 1
		case ebiten.KeyA:
			vec[0] -= 1
		case ebiten.KeyD:
			vec[0] += 1
		}
	}
	state := service.EventBus().GetState()
	state.Vector = vec
	service.EventBus().SetState(state)
	if vec != curVec {
		curVec = vec
		event := contract2.MoveEvent{Event: contract2.Event{Type: contract2.Move}, Vector: vec}
		service.Controller().Send(event)
	}

	return nil
}
