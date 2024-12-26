package contract

type KeyFunc int

const (
	AutoAttack KeyFunc = iota
	SwitchTarget
	Skill1
	Skill2

	Up
	Down
	Left
	Right

	CancelAction
)
