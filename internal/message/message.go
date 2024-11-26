package message

import (
	"github.com/gorilla/websocket"
	"leveling/internal/constract"
	"net/http"
)

type Message struct {
	Type int16
}

var upgrader = websocket.Upgrader{} // use default options

func NewMessenger(server *constract.Server) {
	http.HandleFunc("/socket", newSocketHandler(server))
	http.ListenAndServe("localhost:8080", nil)
}

func newSocketHandler(server *constract.Server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Upgrade our raw HTTP connection to a websocket based one
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			(*server).Log("Error during connection upgradation:%v\n", err)
			return
		}
		defer conn.Close()

		// The event loop
		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				(*server).Log("Error during message reading:%v\n", err)
				break
			}
			(*server).Log("Received: %s\n", message)
			err = conn.WriteMessage(messageType, message)
			if err != nil {
				(*server).Log("Error during message writing:%v\n", err)
				break
			}
		}
	}
}
