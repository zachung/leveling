package keys

import (
	"github.com/gdamore/tcell/v2"
	"leveling/internal/constract"
)

type Func interface {
	Execute(game *constract.Game, event *tcell.EventKey) *tcell.EventKey
}

type CtrlC struct {
	next Func
}

func NewCtrlC(next Func) *CtrlC {
	return &CtrlC{next}
}

func (c *CtrlC) Execute(game *constract.Game, event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyCtrlC {
		(*game).Stop()
		return nil
	}
	if c.next != nil {
		return c.next.Execute(game, event)
	}
	return event
}
