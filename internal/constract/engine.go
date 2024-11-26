package constract

import "io"

type Server interface {
	Start()
	Stop()
	SetConsole(writer *io.Writer)
	Log(string, ...any)
}

type Console interface {
	Info(msg string)
}

type UI interface {
	Logger() *Console
	SideLogger() *Console
	SetController(controller *Controller)
}

type Controller interface {
	Escape()
	Send(message string)
}
