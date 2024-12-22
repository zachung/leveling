package service

import (
	"leveling/internal/server/contract"
)

type Locator struct {
	console contract.Console
	server  contract.Server
	hub     contract.Hub
}

var locator *Locator

func GetLocator() *Locator {
	if locator == nil {
		locator = new(Locator)
	}
	return locator
}

func Logger() contract.Console {
	return locator.console
}

func Server() contract.Server {
	return locator.server
}

func Hub() contract.Hub {
	return locator.hub
}

func (locator *Locator) SetLogger(console contract.Console) *Locator {
	locator.console = console

	return locator
}

func (locator *Locator) SetServer(server contract.Server) *Locator {
	locator.server = server

	return locator
}

func (locator *Locator) SetHub(broadcaster contract.Hub) *Locator {
	locator.hub = broadcaster

	return locator
}
