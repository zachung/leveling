package observers

import (
	contract2 "leveling/internal/contract"
	"leveling/internal/server/contract"
)

type EnemyListener struct {
	hero contract.IHero
}

func NewEnemyListener(hero contract.IHero) *EnemyListener {
	return &EnemyListener{hero: hero}
}

func (e EnemyListener) OnNotify(event contract2.Message) {
	switch event.(type) {
	case contract2.StateChangeEvent:
		changeEvent := event.(contract2.StateChangeEvent)
		e.hero.SetTarget(changeEvent.Attacker.Name)
		e.hero.SetAutoAttack(true)
	}
}
