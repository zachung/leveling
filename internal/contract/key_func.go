package contract

type RoleEvent int

const (
	SetAutoAttack RoleEvent = iota
	SwitchTarget
	Skill1
	Skill2

	Movement
	Up
	Down
	Left
	Right

	CancelAction
)
