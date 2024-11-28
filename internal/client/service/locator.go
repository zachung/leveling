package service

import "leveling/internal/constract"

type Locator struct {
	console    *constract.Console
	controller *constract.Controller
	ui         *constract.UI
	connector  *constract.Connector
}

var locator *Locator

func GetLocator() *Locator {
	if locator == nil {
		locator = new(Locator)
	}
	return locator
}

func Logger() constract.Console {
	return *locator.console
}

func Connector() constract.Connector {
	return *locator.connector
}

func UI() constract.UI {
	return *locator.ui
}

func Controller() constract.Controller {
	return *locator.controller
}

func (locator *Locator) SetLogger(console *constract.Console) *Locator {
	locator.console = console

	return locator
}

func (locator *Locator) SetController(controller *constract.Controller) *Locator {
	locator.controller = controller

	return locator
}

func (locator *Locator) SetConnector(connector *constract.Connector) *Locator {
	locator.connector = connector

	return locator
}

func (locator *Locator) SetUI(ui *constract.UI) *Locator {
	locator.ui = ui

	return locator
}
