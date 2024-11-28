package message

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type Message struct {
	Type int16
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewMessenger() {
	hub := newHub()
	go hub.run()
	http.HandleFunc("/socket", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	http.ListenAndServe("localhost:8080", nil)
}
