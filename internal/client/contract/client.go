package contract

import (
	"github.com/gdamore/tcell/v2"
)

type Console interface {
	Info(msg string, args ...any)
}

type UI interface {
	Logger() *Console
	SideLogger() *Console
	Run()
	Stop()
}

type Controller interface {
	Connect()
	GetKeyBinding() func(event *tcell.EventKey) *tcell.EventKey
	Escape()
	Send(message string)
}

type Connector interface {
	Connect() bool
	Close()
	SendMessage(message string)
}
