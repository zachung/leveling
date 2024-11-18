package weapons

import (
	"leveling/internal/constract"
)

type sword struct {
	Weapon
}

func newSword() constract.IWeapon {
	return &sword{
		Weapon: Weapon{
			power: 3,
		},
	}
}

func (weapon *sword) Attack(hero *constract.IHero) {
	(*hero).ApplyDamage(weapon.holder, weapon.power)
}
