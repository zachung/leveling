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
	p.client.Send(event)
}
