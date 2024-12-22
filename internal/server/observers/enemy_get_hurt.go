package observers

import (
	contract2 "leveling/internal/contract"
	"leveling/internal/server/contract"
)

type EnemyGetHurt struct{}

func (EnemyGetHurt) OnNotify(hero contract.IHero, event contract2.Message) {
	switch event.(type) {
	case contract2.StateChangeEvent:
		changeEvent := event.(contract2.StateChangeEvent)
		hero.SetTarget(changeEvent.Attacker.Name)
		hero.SetAutoAttack(true)
	}
}
