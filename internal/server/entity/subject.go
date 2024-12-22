package entity

import (
	contract2 "leveling/internal/contract"
	"leveling/internal/server/contract"
)

type RoleSubject struct {
	observers []contract.Observer
}

func NewRoleSubject() contract.Subject {
	return new(RoleSubject)
}

func (s *RoleSubject) AddObserver(observer contract.Observer) {
	s.observers = append(s.observers, observer)
}

func (s *RoleSubject) Notify(event contract2.Message) {
	for _, observer := range s.observers {
		observer.OnNotify(event)
	}
}
