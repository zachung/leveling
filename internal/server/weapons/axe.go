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
