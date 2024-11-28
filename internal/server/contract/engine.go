package contract

import contract2 "leveling/internal/contract"

type Server interface {
	Start()
	NewClientConnect(client *Client) *IHero
	LeaveClientConnect(client *Client)
}

type Console interface {
	Info(msg string, args ...any)
	Debug(msg string, args ...any)
}

type Client interface {
	Send(msg []byte) bool
	GetName() string
	Close()
}

type Hub interface {
	Run()
	SendAction(client *Client, action *contract2.Action)
}
