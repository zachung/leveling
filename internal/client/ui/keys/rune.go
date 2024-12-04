package keys

import (
	"github.com/gdamore/tcell/v2"
	"leveling/internal/client/service"
	contract2 "leveling/internal/contract"
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

func (c Rune) handleEvent(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyRune {
		spell := contract2.ActionEvent{Event: contract2.Event{Type: contract2.Action}}

		switch event.Rune() {
		case '1':
			spell.Id = 1
			service.Controller().Send(spell)
		case '2':
			spell.Id = 2
			service.Controller().Send(spell)
		case 's':
			service.Controller().Connect("Sin")
		case 't':
			service.Controller().Connect("Taras")
		case 'b':
			service.Controller().Connect("Brian")
		}
		return nil
	}
	return event
}
