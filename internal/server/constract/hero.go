package constract

type IHero interface {
	Attack(dt float64, targets []*IHero)
	ApplyDamage(from *IHero, damage int)
	IsDie() bool
}
