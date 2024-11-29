package keys

import (
	"github.com/gdamore/tcell/v2"
	"leveling/internal/client/contract"
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

func (c Rune) handleEvent(controller *contract.Controller, event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyRune {
		spell := contract2.Action{Id: 1}

		r := string(event.Rune())
		switch r {
		case "1":
			(*controller).Send(spell.Serialize())
		case "s":
			service.Controller().Connect("Sin")
		case "t":
			service.Controller().Connect("Taras")
		case "b":
			service.Controller().Connect("Brian")
		}
		return nil
	}
	return event
}
