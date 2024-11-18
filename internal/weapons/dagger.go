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
			power: 2,
		},
	}
}

func (weapon dagger) Attack(hero *constract.IHero) {
	for i := 0; i < 4; i++ {
		(*hero).ApplyDamage(weapon.holder, weapon.power)
	}
}
