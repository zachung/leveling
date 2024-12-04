package contract

type IWeapon interface {
	GetSpeed() float64
	GetPower() int
}

const (
	Sword = iota
	Dagger
	Axe
)
