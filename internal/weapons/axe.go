package weapons

import (
	"leveling/internal/constract"
)

type axe struct {
	Weapon
}

func newAxe() constract.IWeapon {
	return &axe{
		Weapon: Weapon{
			power: 6,
		},
	}
}

func (weapon axe) Attack(hero *constract.IHero) {
	(*hero).ApplyDamage(weapon.holder, weapon.power)
}
