package weapons

import (
	"leveling/internal/server/contract"
)

type sword struct {
	Weapon
}

func newSword() contract.IWeapon {
	return &sword{
		Weapon: Weapon{
			power: 3,
			speed: 2,
		},
	}
}

func (weapon sword) Attack(hero *contract.IHero) {
	(*hero).ApplyDamage(weapon.holder, weapon.power)
}
