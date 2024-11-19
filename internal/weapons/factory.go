package weapons

import (
	"leveling/internal/constract"
)

type Weapon struct {
	power  int
	speed  float64
	holder *constract.IHero
}

func (weapon *Weapon) Attack(hero *constract.IHero) {
}

func (weapon *Weapon) SetHolder(hero *constract.IHero) {
	weapon.holder = hero
}

func (weapon *Weapon) GetSpeed() float64 {
	return weapon.speed
}

func NewWeapon(weaponId int) constract.IWeapon {
	if weaponId == constract.Sword {
		return newSword()
	}
	if weaponId == constract.Dagger {
		return newDagger()
	}
	if weaponId == constract.Axe {
		return newAxe()
	}

	return nil
}
