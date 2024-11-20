package weapons

import (
	"leveling/internal/constract"
)

type dagger struct {
	Weapon
}

func newDagger() constract.IWeapon {
	return &dagger{
		Weapon: Weapon{
			power: 1,
			speed: 1,
		},
	}
}

func (weapon dagger) Attack(hero *constract.IHero) {
	(*hero).ApplyDamage(weapon.holder, weapon.power)
}
