package hero

import (
	"fmt"
	"leveling/internal/constract"
)

type Hero struct {
	Name     string
	Health   int
	MainHand *constract.IWeapon
	Strength int
}

func (hero *Hero) Hold(weapon *constract.IWeapon) {
	hero.MainHand = weapon
	iHero := constract.IHero(hero)
	(*weapon).SetHolder(&iHero)
}

func (hero *Hero) Attack(target *constract.IHero) {
	(*hero.MainHand).Attack(target)
}

func (hero *Hero) ApplyDamage(from *constract.IHero, damage int) {
	hero.Health -= damage
	attacker := (*from).(*Hero)
	fmt.Printf("%s(%v) attacked by %s take %v damage\n", hero.Name, hero.Health, attacker.Name, damage)
}

func (hero *Hero) IsDie() bool {
	return hero.Health <= 0
}
