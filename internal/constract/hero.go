package constract

type IHero interface {
	Hold(weapon *IWeapon)
	Attack(target *IHero)
	ApplyDamage(from *IHero, damage int)
	IsDie() bool
}
