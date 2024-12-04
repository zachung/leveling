package observers

import (
	contract2 "leveling/internal/contract"
	"leveling/internal/server/contract"
)

type GetHurt struct{}

func (GetHurt) OnNotify(hero contract.IHero, event contract2.Message) {
	switch event.(type) {
	case contract2.StateChangeEvent:
		changeEvent := event.(contract2.StateChangeEvent)
		hero.SetTarget(changeEvent.Attacker.Name)
		// enable auto attack
		action := contract2.ActionEvent{Id: 1}
		hero.SetNextAction(&action)
	}
}
