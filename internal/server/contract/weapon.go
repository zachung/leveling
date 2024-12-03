package contract

type IWeapon interface {
	SetHolder(hero *IHero)
	GetSpeed() float64
	GetPower() int
}

const (
	Sword = iota
	Dagger
	Axe
)
