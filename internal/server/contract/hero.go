package contract

import contract2 "leveling/internal/contract"

type IHero interface {
	Attack(dt float64, targets []*IHero)
	ApplyDamage(from *IHero, damage int)
	IsDie() bool
	SetNextAction(action *contract2.ActionEvent)
	GetName() string
	GetHealth() int
}
