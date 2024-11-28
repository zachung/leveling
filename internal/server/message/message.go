package message

import (
	"github.com/gorilla/websocket"
	"leveling/internal/server/service"
	"net/http"
)

type Message struct {
	Type int16
}

var upgrader = websocket.Upgrader{} // use default options

func NewMessenger() {
	http.HandleFunc("/socket", socketHandler)
	http.ListenAndServe("localhost:8080", nil)
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		service.Logger().Info("Error during connection upgradation:%v\n", err)
		return
	}
	defer conn.Close()

	// The event loop
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			service.Logger().Info("Error during message reading:%v\n", err)
			break
		}
		service.Logger().Info("Received: %s\n", message)
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			service.Logger().Info("Error during message writing:%v\n", err)
			break
		}
	}
}
