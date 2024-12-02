package keys

import (
	"github.com/gdamore/tcell/v2"
	"leveling/internal/client/service"
)

type WorldPanel struct {
	*T
}

func NewWorldPanel(next Func) *WorldPanel {
	t := &T{next: next}
	i := &WorldPanel{T: t}
	t.Func = i

	return i
}

func (c WorldPanel) handleEvent(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyF2 {
		service.UI().World().Focus()
		return nil
	}
	return event
}
