package constract

type IHero interface {
	Attack(targets []*IHero)
	ApplyDamage(from *IHero, damage int)
	IsDie() bool
}
