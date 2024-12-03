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
