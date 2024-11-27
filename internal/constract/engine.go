package constract

import (
	"github.com/gdamore/tcell/v2"
	"io"
)

type Server interface {
	Start()
	Stop()
	SetConsole(writer *io.Writer)
	Log(string, ...any)
}

type Console interface {
	Info(msg string, args ...any)
}

type UI interface {
	Logger() *Console
	SideLogger() *Console
	SetController(controller *Controller)
	SetKeyBinding(keyBinding func(event *tcell.EventKey) *tcell.EventKey)
	Stop()
}

type Controller interface {
	Connect()
	GetKeyBinding() func(event *tcell.EventKey) *tcell.EventKey
	Escape()
	Send(message string)
}

type Connection interface {
	Connect() bool
	Close()
	SendMessage(message string)
}
