package weapons

import (
	"leveling/internal/server/contract"
)

type Weapon struct {
	power int
	speed float64
}

func (weapon *Weapon) GetSpeed() float64 {
	return weapon.speed
}

func (weapon *Weapon) GetPower() int {
	return weapon.power
}

func NewWeapon(weaponId int) contract.IWeapon {
	if weaponId == contract.Sword {
		return newSword()
	}
	if weaponId == contract.Dagger {
		return newDagger()
	}
	if weaponId == contract.Axe {
		return newAxe()
	}

	return nil
}
