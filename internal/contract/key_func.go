package contract

type RoleEvent int

const (
	SetAutoAttack RoleEvent = iota
	SwitchTarget
	Skill1
	Skill2

	Up
	Down
	Left
	Right

	CancelAction
)
