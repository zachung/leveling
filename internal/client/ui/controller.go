package ui

import (
	"leveling/internal/client/contract"
	"leveling/internal/client/service"
	contract2 "leveling/internal/contract"
)

type Controller struct {
}

func NewController() contract.Controller {
	return new(Controller)
}

func (c *Controller) Connect(name string) {
	go func() {
		if service.Connector().Connect(name) {
			// TODO: another key binding
		}
	}()
}

func (c *Controller) Escape() {
	go func() {
		service.Connector().Close()
		service.UI().Stop()
	}()
}

func (c *Controller) Send(message contract2.Message) {
	service.Connector().SendMessage(message)
}
