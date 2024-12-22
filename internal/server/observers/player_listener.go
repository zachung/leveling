package observers

import (
	contract2 "leveling/internal/contract"
	"leveling/internal/server/contract"
)

type PlayerListener struct {
	client contract.Client
}

func NewPlayerListener(client contract.Client) *PlayerListener {
	return &PlayerListener{client: client}
}

func (p PlayerListener) OnNotify(event contract2.Message) {
	switch event.(type) {
	case contract2.StateChangeEvent:
		p.client.Send(event)
	case contract2.HeroDieEvent:
		p.client.Send(event)
	}
}
