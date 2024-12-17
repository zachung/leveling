package keys

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"leveling/internal/client/service"
)

type SwitchTarget struct {
	*T
}

func NewSwitchTarget(next Func) *SwitchTarget {
	t := &T{next: next}
	i := &SwitchTarget{T: t}
	t.Func = i

	return i
}

func (c SwitchTarget) handleEvent() *ebiten.Key {
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		service.EventBus().SelectNext()
		key := ebiten.KeyTab
		return &key
	}
	return nil
}
