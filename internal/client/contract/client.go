package contract

import (
	"leveling/internal/contract"
)

type Chat interface {
	Info(msg string, args ...any)
}

type UI interface {
	Run()
	Stop()
	State() State
}

type State interface {
	UpdateState(event contract.StateChangeEvent)
}

type Controller interface {
	Connect(name string)
	Escape()
	Send(message contract.Message)
}

type Connector interface {
	Connect(name string) bool
	Close()
	SendMessage(message contract.Message)
	GetCurName() string
}

type BusEvent int

const (
	OnStateChanged BusEvent = iota
	OnWorldChanged
	OnReportAppend
)

type Bus interface {
	AddObserver(event BusEvent, observer func())
	SetState(event contract.StateChangeEvent)
	GetState() contract.StateChangeEvent
	SetWorldState(event contract.WorldEvent)
	GetWorldState() contract.WorldEvent
	SelectNext()
	AppendReport(text string)
	GetReport() string
}
