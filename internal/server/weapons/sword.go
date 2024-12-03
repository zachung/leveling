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
