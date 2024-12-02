package keys

import (
	"github.com/gdamore/tcell/v2"
	"leveling/internal/client/service"
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

func (c CtrlC) handleEvent(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyCtrlC {
		service.Controller().Escape()
		return nil
	}
	return event
}
