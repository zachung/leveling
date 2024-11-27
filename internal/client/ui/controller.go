package ui

import (
	"github.com/gdamore/tcell/v2"
	"leveling/internal/client/ui/keys"
	"leveling/internal/constract"
	"time"
)

type Controller struct {
	ui         *constract.UI
	connection *constract.Connection
}

func NewController(ui *constract.UI, connection *constract.Connection) *constract.Controller {
	var controller constract.Controller
	c := &Controller{ui, connection}
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
		if (*c.connection).Connect() {
			KeyLogger().Info("Connected!\n")
			// TODO: another key binding
		}
	}()
}

func (c *Controller) Escape() {
	KeyLogger().Info("Stopping...\n")
	go func() {
		(*c.connection).Close()
		time.Sleep(1 * time.Second)
		(*c.ui).Stop()
	}()
}

func (c *Controller) Send(message string) {
	KeyLogger().Info("%v\n", message)
	(*c.connection).SendMessage(message)
}
