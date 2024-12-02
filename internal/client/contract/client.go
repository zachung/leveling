package contract

import (
	"github.com/gdamore/tcell/v2"
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
}

type State interface {
	UpdateState(event contract.StateChangeEvent)
}

type World interface {
	UpdateWorld(event contract.WorldEvent)
}

type Controller interface {
	Connect(name string)
	GetKeyBinding() func(event *tcell.EventKey) *tcell.EventKey
	Escape()
	Send(message contract.Message)
}

type Connector interface {
	Connect(name string) bool
	Close()
	SendMessage(message contract.Message)
}
