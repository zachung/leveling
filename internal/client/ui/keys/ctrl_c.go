package keys

import (
	"github.com/gdamore/tcell/v2"
	"leveling/internal/constract"
)

type CtrlC struct {
	*T
}

func NewCtrlC(next Func) *CtrlC {
	t := &T{next: next}
	i := &CtrlC{T: t}
	t.Func = i

	return i
}

func (c CtrlC) handleEvent(controller *constract.Controller, event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyCtrlC {
		(*controller).Escape()
		return nil
	}
	return event
}