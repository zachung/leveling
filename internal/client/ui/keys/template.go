package keys

import (
	"github.com/gdamore/tcell/v2"
	"leveling/internal/client/service"
)

type Func interface {
	Execute(event *tcell.EventKey) *tcell.EventKey
	handleEvent(event *tcell.EventKey) *tcell.EventKey
}

type T struct {
	Func
	next Func
}

func (t *T) Execute(event *tcell.EventKey) *tcell.EventKey {
	if t.handleEvent(event) == nil {
		service.SideLogger().Info("%v\n", event.Name())
		return nil
	}
	if t.next != nil {
		return t.next.Execute(event)
	}
	return event
}
