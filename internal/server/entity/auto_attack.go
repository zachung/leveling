package entity

import (
	contract2 "leveling/internal/contract"
	"leveling/internal/server/contract"
)

type AutoAttackAbility struct {
	hero           *Hero
	isAutoAttack   bool
	attackCooldown float64 // weapon auto attack cooldown
}

func NewAutoAttackAbility(hero *Hero) Ability {
	a := &AutoAttackAbility{hero: hero}
	hero.AddOperationListener(contract2.SetAutoAttack, func(args ...any) {
		if hero.target == nil {
			return
		}
		a.isAutoAttack = args[0].(bool)
	})
	hero.AddOperationListener(contract2.CancelAction, func(...any) {
		a.isAutoAttack = false
	})

	return a
}

func (a *AutoAttackAbility) Update(dt float64) {
	hero := a.hero
	a.attackCooldown += dt / hero.mainHand.GetSpeed()
	for {
		if a.attackCooldown < roundTimeSecond {
			return
		}
		a.attackCooldown -= roundTimeSecond
		a.doAutoAttack(hero)
	}
}

func (a *AutoAttackAbility) doAutoAttack(hero *Hero) {
	if !a.isAutoAttack {
		return
	}
	if hero.target == nil {
		a.isAutoAttack = false
		return
	}
	target := hero.target.(*Hero)
	if target.IsDie() {
		a.isAutoAttack = false
		hero.target = nil
		return
	}
	hero.isActive = true
	damage := contract.Damage(hero.mainHand.GetPower())
	target.ApplyDamage(damage)
	if target.IsDie() {
		a.isAutoAttack = false
		hero.SetAction(contract2.ActionEvent{Event: contract2.Event{Type: contract2.Action}, Id: contract2.CancelAction})
	}
	// send message to client
	messageEvent(hero, damage, target)
}
