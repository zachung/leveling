package entity

import (
	contract2 "leveling/internal/contract"
	"leveling/internal/server/contract"
)

type ActionAbility struct {
	hero       *Hero
	nextAction *contract2.ActionEvent
}

func NewActionAbility(hero *Hero) Ability {
	return &ActionAbility{hero: hero}
}

func (a *ActionAbility) Update(dt float64) {
	hero := a.hero
	hero.roundCooldown += dt
	if hero.roundCooldown < globalCoolDown {
		return
	}
	if a.nextAction == nil {
		return
	}
	hero.roundCooldown = 0
	a.doAction(hero)
}

func (a *ActionAbility) doAction(hero *Hero) {
	if hero.target == nil {
		return
	}
	target := hero.target.(*Hero)
	if target.IsDie() {
		hero.target = nil
		return
	}
	hero.isActive = true
	damage := contract.Damage(hero.mainHand.GetPower() + hero.strength)
	target.ApplyDamage(damage)
	if target.IsDie() {
		hero.SetAction(contract2.ActionEvent{Event: contract2.Event{Type: contract2.Action}, Id: contract2.CancelAction})
	}
	a.nextAction = nil

	// send message to client
	messageEvent(hero, damage, target)
}
