package hero

import "fmt"

type AttackType int

const (
	StraightSword AttackType = iota
	CurvedSword
	Dagger
	Claw
	Polearm
	Mace
	Hammer
	Axe
)

type Hero struct {
	Name       string
	Health     int
	AttackType AttackType
	Strength   int
}

func (hero *Hero) Attack(target *Hero) {
	switch hero.AttackType {
	case StraightSword:
		target.applyDamage(hero, 6)
	case CurvedSword:
		target.applyDamage(hero, 8)
	case Dagger:
		target.applyDamage(hero, 2)
		target.applyDamage(hero, 2)
		target.applyDamage(hero, 2)
	case Claw:
		target.applyDamage(hero, 3)
		target.applyDamage(hero, 3)
	case Polearm:
		target.applyDamage(hero, 13)
	case Mace:
		target.applyDamage(hero, 8)
	case Hammer:
		target.applyDamage(hero, 10)
	case Axe:
		target.applyDamage(hero, 10)
	default:
		target.applyDamage(hero, hero.Strength)
	}
}

func (hero *Hero) applyDamage(from *Hero, damage int) {
	hero.Health -= damage
	fmt.Printf("%s(%v) attack %s(%v) make %v damage\n", from.Name, from.Health, hero.Name, hero.Health, damage)
}

func (hero *Hero) IsDie() bool {
	return hero.Health <= 0
}
