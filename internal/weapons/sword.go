package weapons

import (
	"leveling/internal/server/constract"
)

type sword struct {
	Weapon
}

func newSword() constract.IWeapon {
	return &sword{
		Weapon: Weapon{
			power: 3,
			speed: 2,
		},
	}
}

func (weapon sword) Attack(hero *constract.IHero) {
	(*hero).ApplyDamage(weapon.holder, weapon.power)
}
