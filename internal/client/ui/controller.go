package ui

import (
	"github.com/gdamore/tcell/v2"
	"leveling/internal/client/contract"
	"leveling/internal/client/service"
	"leveling/internal/client/ui/keys"
	contract2 "leveling/internal/contract"
)

type Controller struct {
}

func NewController() *contract.Controller {
	var controller contract.Controller
	c := &Controller{}
	controller = contract.Controller(c)

	return &controller
}

func (c *Controller) GetKeyBinding() func(event *tcell.EventKey) *tcell.EventKey {
	controller := contract.Controller(c)

	return func(event *tcell.EventKey) *tcell.EventKey {
		// chain of responsibility
		keyHandlers := keys.NewCtrlC(keys.NewRune(nil))

		if (*keyHandlers).Execute(&controller, event) == nil {
			return nil
		}
		return event
	}
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
