package keys

import (
	"github.com/gdamore/tcell/v2"
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

func (c SwitchTarget) handleEvent(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyTAB {
		service.UI().World().SelectNext()
		return nil
	}
	return event
}
