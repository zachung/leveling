package weapons

import (
	"leveling/internal/server/contract"
)

type axe struct {
	Weapon
}

func newAxe() contract.IWeapon {
	return &axe{
		Weapon: Weapon{
			power: 600,
			speed: 4,
		},
	}
}

func (weapon axe) Attack(hero *contract.IHero) {
	(*hero).ApplyDamage(weapon.holder, weapon.power)
}
