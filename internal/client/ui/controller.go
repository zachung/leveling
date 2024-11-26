package ui

import (
	"leveling/internal/client/message"
	"leveling/internal/constract"
)

type Controller struct {
	connection *message.Connection
}

func NewController() *constract.Controller {
	var controller constract.Controller
	connection := message.NewConnection()
	c := &Controller{connection}
	controller = constract.Controller(c)

	return &controller
}

func (c *Controller) Escape() {
	c.connection.Close()
}

func (c *Controller) Send(message string) {
	c.connection.SendMessage(message)
}
