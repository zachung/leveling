package service

import (
	"leveling/internal/client/contract"
)

type Locator struct {
	console     *contract.Console
	controller  *contract.Controller
	ui          *contract.UI
	connector   *contract.Connector
	sideConsole *contract.Console
	bus         contract.Bus
}

var locator *Locator

func GetLocator() *Locator {
	if locator == nil {
		locator = new(Locator)
	}
	return locator
}

func Logger() contract.Console {
	return *locator.console
}

func Connector() contract.Connector {
	return *locator.connector
}

func UI() contract.UI {
	return *locator.ui
}

func Controller() contract.Controller {
	return *locator.controller
}

func SideLogger() contract.Console {
	return *locator.sideConsole
}

func EventBus() contract.Bus {
	return locator.bus
}

func (locator *Locator) SetLogger(console *contract.Console) *Locator {
	locator.console = console

	return locator
}

func (locator *Locator) SetController(controller *contract.Controller) *Locator {
	locator.controller = controller

	return locator
}

func (locator *Locator) SetConnector(connector *contract.Connector) *Locator {
	locator.connector = connector

	return locator
}

func (locator *Locator) SetUI(ui *contract.UI) *Locator {
	locator.ui = ui

	return locator
}

func (locator *Locator) SetKeyLogger(console *contract.Console) *Locator {
	locator.sideConsole = console

	return locator
}

func (locator *Locator) SetBus(bus contract.Bus) *Locator {
	locator.bus = bus

	return locator
}
