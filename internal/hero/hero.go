package hero

import (
	"fmt"
	"leveling/internal/constract"
	"leveling/internal/entity"
	"leveling/internal/weapons"
)

type Hero struct {
	Name     string
	Health   int
	MainHand *constract.IWeapon
	Strength int
}

func New(data entity.Hero) *Hero {
	weapon := weapons.NewWeapon(data.MainHand)
	hero := &Hero{
		Name:     data.Name,
		Health:   data.Health,
		Strength: data.Strength,
		MainHand: &weapon,
	}
	iHero := constract.IHero(hero)
	weapon.SetHolder(&iHero)

	return hero
}

func (hero *Hero) Attack(targets []*constract.IHero) {
	(*hero.MainHand).Attack(targets[0])
}

func (hero *Hero) ApplyDamage(from *constract.IHero, damage int) {
	hero.Health -= damage
	attacker := (*from).(*Hero)
	fmt.Printf("%s(%v) attacked by %s take %v damage\n", hero.Name, hero.Health, attacker.Name, damage)
}

func (hero *Hero) IsDie() bool {
	return hero.Health <= 0
}
