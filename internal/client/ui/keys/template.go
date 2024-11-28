package keys

import (
	"github.com/gdamore/tcell/v2"
	"leveling/internal/client/contract"
)

type Func interface {
	Execute(controller *contract.Controller, event *tcell.EventKey) *tcell.EventKey
	handleEvent(controller *contract.Controller, event *tcell.EventKey) *tcell.EventKey
}

type T struct {
	Func
	next Func
}

func (t *T) Execute(controller *contract.Controller, event *tcell.EventKey) *tcell.EventKey {
	if t.handleEvent(controller, event) == nil {
		return nil
	}
	if t.next != nil {
		return t.next.Execute(controller, event)
	}
	return event
}
