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
	Run(name string)
	Stop()
}

type Controller interface {
	Connect()
	GetKeyBinding() func(event *tcell.EventKey) *tcell.EventKey
	Escape()
	Send(message []byte)
}

type Connector interface {
	Connect() bool
	Close()
	SendMessage(message []byte)
}
