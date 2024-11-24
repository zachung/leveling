package keys

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"leveling/internal/constract"
)

type Rune struct {
	next Func
}

func NewRune(next Func) *Rune {
	return &Rune{next}
}

func (c Rune) Execute(game *constract.Game, event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyRune {
		(*(*(*game).UI()).SideLogger()).Info(fmt.Sprintf("type in %v", string(event.Rune())))
		return nil
	}
	if c.next != nil {
		return c.next.Execute(game, event)
	}
	return event
}
