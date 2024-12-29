package entity

type AbilityType int

const (
	AutoAttack AbilityType = iota
	Movement
	Action
)

type Ability interface {
	Update(dt float64)
}
