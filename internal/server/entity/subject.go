package entity

import (
	contract2 "leveling/internal/contract"
	"leveling/internal/server/contract"
)

type Subject struct {
	observers []*contract.Observer
}

func (s *Subject) AddObserver(observer *contract.Observer) {
	s.observers = append(s.observers, observer)
}

func (s *Subject) Notify(hero contract.IHero, event contract2.Message) {
	for _, observer := range s.observers {
		(*observer).OnNotify(hero, event)
	}
}
