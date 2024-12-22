package connection

import (
	contract2 "leveling/internal/contract"
	"leveling/internal/server/contract"
	"leveling/internal/server/service"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]contract.IHero

	// Inbound messages from the clients.
	broadcast chan contract2.Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() contract.Hub {
	h := &Hub{
		broadcast:  make(chan contract2.Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]contract.IHero),
	}

	return h
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			hero := service.Server().NewClientConnect(client)
			h.clients[client] = hero
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
			}
			service.Server().LeaveClientConnect(client)
		case message := <-h.broadcast:
			for client := range h.clients {
				if !(*client).Send(message) {
					(*client).Close()
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *Hub) SendAction(client contract.Client, action contract2.Message) {
	c := client.(*Client)
	iHero := h.clients[c]
	switch action.(type) {
	case contract2.ActionEvent:
		event := action.(contract2.ActionEvent)
		iHero.SetNextAction(&event)
	case contract2.SelectTargetEvent:
		event := action.(contract2.SelectTargetEvent)
		iHero.SetTarget(event.Name)
	}
}

func (h *Hub) Broadcast(m contract2.Message) {
	h.broadcast <- m
}
