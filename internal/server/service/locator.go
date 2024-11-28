package service

import (
	"leveling/internal/server/constract"
)

type Locator struct {
	console *constract.Console
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

func (locator *Locator) SetLogger(console *constract.Console) *Locator {
	locator.console = console

	return locator
}
