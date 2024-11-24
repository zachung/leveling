package constract

type Game interface {
	Start()
	Stop()
	UI() *UI
}

type Console interface {
	Info(msg string)
}

type UI interface {
	Logger() *Console
	SideLogger() *Console
}
