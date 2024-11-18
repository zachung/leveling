package constract

type IWeapon interface {
	Attack(hero *IHero)
	SetHolder(hero *IHero)
}

const (
	Sword = iota
	Dagger
	Axe
)
