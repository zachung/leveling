package message

import (
	"github.com/gorilla/websocket"
	"leveling/internal/client/ui"
	"leveling/internal/constract"
	"time"
)

type Connection struct {
	console *constract.Console
	conn    *websocket.Conn
}

var done chan interface{}

func NewConnection(console *constract.Console) *constract.Connection {
	c := &Connection{console: console}
	connection := constract.Connection(c)

	return &connection
}

func (c *Connection) Connect() bool {
	socketUrl := "ws://localhost:8080" + "/socket"
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		ui.Logger().Info("Error connecting to Websocket Server:%v\n", err)
		return false
	}
	c.conn = conn
	done = make(chan interface{})
	go receiveHandler(conn)

	return true
}

func receiveHandler(connection *websocket.Conn) {
	defer close(done)
	for {
		_, msg, err := connection.ReadMessage()
		if err != nil {
			ui.Logger().Info("Error in receive:%v\n", err)
			return
		}
		ui.Logger().Info("Received: %s\n", msg)
	}
}

func (c *Connection) Close() {
	if c.conn == nil {
		return
	}
	defer c.conn.Close()

	// Terminate gracefully...
	ui.Logger().Info("Closing all pending connections\n")

	// Close our websocket connection
	err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		ui.Logger().Info("Error during closing websocket:%v\n", err)
		return
	}

	select {
	case <-done:
		ui.Logger().Info("Receiver Channel Closed! Exiting....\n")
	case <-time.After(time.Duration(1) * time.Second):
		ui.Logger().Info("Timeout in closing receiving channel. Exiting....\n")
	}
	return
}

func (c *Connection) SendMessage(message string) {
	c.conn.WriteMessage(websocket.TextMessage, []byte(message))
}

func (c *Connection) StartTest() {
	// Our main loop for the client
	// We send our relevant packets here
	for {
		select {
		case <-time.After(time.Duration(1) * time.Millisecond * 1000):
			// Send an echo packet every second
			err := c.conn.WriteMessage(websocket.TextMessage, []byte("Hello from GolangDocs!"))
			if err != nil {
				ui.Logger().Info("Error during writing to websocket:%v\n", err)
				return
			}
		}
	}
}
