package connection

import (
	"github.com/gorilla/websocket"
	"leveling/internal/server/service"
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
	http.HandleFunc("/socket", func(w http.ResponseWriter, r *http.Request) {
		serveWs((service.Hub()).(*Hub), w, r)
	})
	http.ListenAndServe("localhost:8080", nil)
}
