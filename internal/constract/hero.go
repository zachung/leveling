package constract

type IHero interface {
	Attack(target *IHero)
	ApplyDamage(from *IHero, damage int)
	IsDie() bool
}
