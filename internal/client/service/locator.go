package service

import (
	"leveling/internal/client/contract"
)

type Locator struct {
	chat       contract.Chat
	controller contract.Controller
	ui         contract.UI
	connector  contract.Connector
	bus        contract.Bus
}

var locator *Locator

func GetLocator() *Locator {
	if locator == nil {
		locator = new(Locator)
	}
	return locator
}

func Chat() contract.Chat {
	return locator.chat
}

func Connector() contract.Connector {
	return locator.connector
}

func UI() contract.UI {
	return locator.ui
}

func Controller() contract.Controller {
	return locator.controller
}

func EventBus() contract.Bus {
	return locator.bus
}

func (locator *Locator) SetChat(chat contract.Chat) *Locator {
	locator.chat = chat

	return locator
}

func (locator *Locator) SetController(controller contract.Controller) *Locator {
	locator.controller = controller

	return locator
}

func (locator *Locator) SetConnector(connector contract.Connector) *Locator {
	locator.connector = connector

	return locator
}

func (locator *Locator) SetUI(ui contract.UI) *Locator {
	locator.ui = ui

	return locator
}

func (locator *Locator) SetBus(bus contract.Bus) *Locator {
	locator.bus = bus

	return locator
}
