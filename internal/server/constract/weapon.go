package constract

type IWeapon interface {
	Attack(hero *IHero)
	SetHolder(hero *IHero)
	GetSpeed() float64
}

const (
	Sword = iota
	Dagger
	Axe
)
