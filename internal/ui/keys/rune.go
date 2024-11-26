package keys

import (
	"github.com/gdamore/tcell/v2"
	"leveling/internal/constract"
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

func (c Rune) handleEvent(server *constract.Server, event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyRune {
		// TODO: the control message for client
		//(*(*(*server).UI()).SideLogger()).Info(fmt.Sprintf("type in %v", string(event.Rune())))
		return nil
	}
	return event
}
