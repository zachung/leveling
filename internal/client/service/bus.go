package service

import (
	contract2 "leveling/internal/client/contract"
	"leveling/internal/contract"
)

type Bus struct {
	stateEvent contract.StateChangeEvent
	worldEvent contract.WorldEvent
	report     string
	observers  map[contract2.BusEvent][]func()
}

func NewBus() contract2.Bus {
	return contract2.Bus(&Bus{
		observers: make(map[contract2.BusEvent][]func()),
	})
}

func (b *Bus) AddObserver(event contract2.BusEvent, observer func()) {
	b.observers[event] = append(b.observers[event], observer)
}

func (b *Bus) SetState(event contract.StateChangeEvent) {
	b.stateEvent = event
	for _, f := range b.observers[contract2.OnStateChanged] {
		f()
	}
}

func (b *Bus) GetState() contract.StateChangeEvent {
	return b.stateEvent
}

func (b *Bus) SetWorldState(event contract.WorldEvent) {
	b.worldEvent = event
	for _, f := range b.observers[contract2.OnWorldChanged] {
		f()
	}
}

func (b *Bus) GetWorldState() contract.WorldEvent {
	return b.worldEvent
}

func (b *Bus) SelectNext() {
	for _, f := range b.observers[contract2.OnSelectTarget] {
		f()
	}
}

func (b *Bus) AppendReport(text string) {
	b.report += text
	for _, f := range b.observers[contract2.OnReportAppend] {
		f()
	}
}

func (b *Bus) GetReport() string {
	return b.report
}
