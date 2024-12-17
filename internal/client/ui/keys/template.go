package keys

import (
	"github.com/hajimehoshi/ebiten/v2"
	"leveling/internal/client/service"
)

type Func interface {
	Execute() *ebiten.Key
	handleEvent() *ebiten.Key
}

type T struct {
	Func
	next Func
}

func (t *T) Execute() *ebiten.Key {
	key := t.handleEvent()
	if key != nil {
		service.SideLogger().Info("%v\n", key)
		return nil
	}
	if t.next != nil {
		return t.next.Execute()
	}
	return key
}
