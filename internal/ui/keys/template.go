package keys

import (
	"github.com/gdamore/tcell/v2"
	"leveling/internal/constract"
)

type Func interface {
	Execute(game *constract.Game, event *tcell.EventKey) *tcell.EventKey
	handleEvent(game *constract.Game, event *tcell.EventKey) *tcell.EventKey
}

type T struct {
	Func
	next Func
}

func (t *T) Execute(game *constract.Game, event *tcell.EventKey) *tcell.EventKey {
	if t.handleEvent(game, event) == nil {
		return nil
	}
	if t.next != nil {
		return t.next.Execute(game, event)
	}
	return event
}
