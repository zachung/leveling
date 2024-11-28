package service

import (
	"leveling/internal/server/contract"
)

type Locator struct {
	console *contract.Console
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

func (locator *Locator) SetLogger(console *contract.Console) *Locator {
	locator.console = console

	return locator
}
