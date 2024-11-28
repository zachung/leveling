package keys

import (
	"github.com/gdamore/tcell/v2"
	"leveling/internal/client/contract"
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
		if r == "1" {
			(*controller).Send(spell.Serialize())
		}
		return nil
	}
	return event
}
