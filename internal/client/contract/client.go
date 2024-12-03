package contract

import (
	"leveling/internal/contract"
)

type Console interface {
	Info(msg string, args ...any)
}

type UI interface {
	Logger() *Console
	SideLogger() *Console
	Run()
	Stop()
	State() State
	World() World
	Report() Panel
}

type State interface {
	UpdateState(event contract.StateChangeEvent)
}

type Panel interface {
	Focus()
}

type World interface {
	Panel
	UpdateWorld(event contract.WorldEvent)
	SelectTarget(index int)
}

type Controller interface {
	Connect(name string)
	Escape()
	Send(message contract.Message)
}

type Connector interface {
	Connect(name string) bool
	Close()
	SendMessage(message contract.Message)
	GetCurName() string
}
