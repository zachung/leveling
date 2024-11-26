package constract

import "io"

type Game interface {
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
}
