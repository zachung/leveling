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

func NewMessenger(game *constract.Game) {
	http.HandleFunc("/socket", newSocketHandler(game))
	http.ListenAndServe("localhost:8080", nil)
}

func newSocketHandler(game *constract.Game) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Upgrade our raw HTTP connection to a websocket based one
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			(*game).Log("Error during connection upgradation:%v\n", err)
			return
		}
		defer conn.Close()

		// The event loop
		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				(*game).Log("Error during message reading:%v\n", err)
				break
			}
			(*game).Log("Received: %s\n", message)
			err = conn.WriteMessage(messageType, message)
			if err != nil {
				(*game).Log("Error during message writing:%v\n", err)
				break
			}
		}
	}
}
