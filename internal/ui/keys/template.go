package keys

import (
	"github.com/gdamore/tcell/v2"
	"leveling/internal/constract"
)

type Func interface {
	Execute(server *constract.Server, event *tcell.EventKey) *tcell.EventKey
	handleEvent(server *constract.Server, event *tcell.EventKey) *tcell.EventKey
}

type T struct {
	Func
	next Func
}

func (t *T) Execute(server *constract.Server, event *tcell.EventKey) *tcell.EventKey {
	if t.handleEvent(server, event) == nil {
		return nil
	}
	if t.next != nil {
		return t.next.Execute(server, event)
	}
	return event
}
