package contract

import contract2 "leveling/internal/contract"

type Damage int

type IHero interface {
	Update(dt float64) bool
	IsDie() bool
	SetNextAction(action *contract2.ActionEvent)
	GetName() string
	GetHealth() int
	SetTarget(name string)
	SetRound(round *Round)
	ApplyDamage(damage Damage)
}
