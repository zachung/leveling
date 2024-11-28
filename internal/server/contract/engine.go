package contract

type Server interface {
	Start()
	NewClientConnect(client *Client)
	LeaveClientConnect(client *Client)
}

type Console interface {
	Info(msg string, args ...any)
}

type Client interface {
	Send(msg string)
	GetName() string
}
