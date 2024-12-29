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
		spell := contract2.ActionEvent{Event: contract2.Event{Type: contract2.Action}}
		spell.Id = contract2.SetAutoAttack
		spell.IsEnable = true
		e.hero.SetAction(spell)
		log.Infof("%v target %s", e.hero.GetName(), hurtEvent.From.Name)
	}
}
