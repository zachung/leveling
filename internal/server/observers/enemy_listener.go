package observers

import (
	log "github.com/sirupsen/logrus"
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
	case contract2.GetHurtEvent:
		hurtEvent := event.(contract2.GetHurtEvent)
		e.hero.SetTarget(hurtEvent.From.Name)
		e.hero.SetAutoAttack(true)
		log.Infof("%v target %s", e.hero.GetName(), hurtEvent.From.Name)
	}
}
