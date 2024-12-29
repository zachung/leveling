package entity

import (
	"golang.org/x/image/math/f64"
	contract2 "leveling/internal/contract"
)

type MoveAbility struct {
	hero          *Hero
	moveThreshold float64
	position      f64.Vec2
	vector        f64.Vec2
}

func NewMoveAbility(hero *Hero) Ability {
	a := &MoveAbility{hero: hero}
	hero.AddOperationListener(contract2.Right, func(args ...interface{}) {
		event := args[0].(contract2.MoveEvent)
		a.vector = event.Vector
	})

	return a
}

func (a *MoveAbility) Update(dt float64) {
	hero := a.hero
	a.moveThreshold += dt
	if a.moveThreshold <= 0.03 {
		return
	}
	var dv float64
	a.moveThreshold, dv = 0, a.moveThreshold
	if a.vector[0] == 0 && a.vector[1] == 0 {
		return
	}
	hero.position[0] += a.vector[0] * dv * 160
	hero.position[1] += a.vector[1] * dv * 160
}
