package weapons

import (
	"leveling/internal/server/contract"
)

type dagger struct {
	Weapon
}

func newDagger() contract.IWeapon {
	return &dagger{
		Weapon: Weapon{
			power: 1,
			speed: 1,
		},
	}
}

func (weapon dagger) Attack(hero *contract.IHero) {
	(*hero).ApplyDamage(weapon.holder, weapon.power)
}
