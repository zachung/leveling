package constract

type Server interface {
	Start()
	Stop()
}

type Console interface {
	Info(msg string, args ...any)
}
