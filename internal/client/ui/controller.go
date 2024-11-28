package ui

import (
	"github.com/gdamore/tcell/v2"
	"leveling/internal/client/constract"
	"leveling/internal/client/service"
	"leveling/internal/client/ui/keys"
	"time"
)

type Controller struct {
}

func NewController() *constract.Controller {
	var controller constract.Controller
	c := &Controller{}
	controller = constract.Controller(c)

	return &controller
}

func (c *Controller) GetKeyBinding() func(event *tcell.EventKey) *tcell.EventKey {
	controller := constract.Controller(c)

	return func(event *tcell.EventKey) *tcell.EventKey {
		// chain of responsibility
		keyHandlers := keys.NewCtrlC(nil)

		if (*keyHandlers).Execute(&controller, event) == nil {
			return nil
		}
		return event
	}
}

func (c *Controller) Connect() {
	KeyLogger().Info("Connect to server...\n")
	go func() {
		if service.Connector().Connect() {
			KeyLogger().Info("Connected!\n")
			// TODO: another key binding
		}
	}()
}

func (c *Controller) Escape() {
	KeyLogger().Info("Stopping...\n")
	go func() {
		service.Connector().Close()
		time.Sleep(1 * time.Second)
		service.UI().Stop()
	}()
}

func (c *Controller) Send(message string) {
	KeyLogger().Info("%v\n", message)
	service.Connector().SendMessage(message)
}
