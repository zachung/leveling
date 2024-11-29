package message

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
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *contract.Hub {
	h := &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]contract.IHero),
	}
	hub := contract.Hub(h)

	return &hub
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			c := contract.Client(client)
			hero := service.Server().NewClientConnect(&c)
			h.clients[client] = *hero
		case client := <-h.unregister:
			c := contract.Client(client)
			service.Server().LeaveClientConnect(&c)
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				(*client).Close()
			}
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

func (h *Hub) SendAction(client *contract.Client, action *contract2.Action) {
	c := (*client).(*Client)
	h.clients[c].SetNextAction(action)
}

func (h *Hub) Broadcast(m []byte) {
	h.broadcast <- m
}
