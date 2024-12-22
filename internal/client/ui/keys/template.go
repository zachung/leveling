package keys

import (
	"github.com/hajimehoshi/ebiten/v2"
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
		// key handled
		return nil
	}
	if t.next != nil {
		return t.next.Execute()
	}
	return key
}
