package keys

import (
	"github.com/gdamore/tcell/v2"
	"leveling/internal/client/contract"
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

func (c CtrlC) handleEvent(controller *contract.Controller, event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyCtrlC {
		(*controller).Escape()
		return nil
	}
	return event
}
