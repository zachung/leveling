package message

import (
	"github.com/gorilla/websocket"
	"leveling/internal/client/contract"
	"leveling/internal/client/service"
	"time"
)

type Connector struct {
	conn *websocket.Conn
}

var done chan interface{}

func NewConnection() *contract.Connector {
	c := &Connector{}
	connector := contract.Connector(c)

	return &connector
}

func (c *Connector) Connect() bool {
	socketUrl := "ws://localhost:8080" + "/socket"
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		service.Logger().Info("Error connecting to Websocket Server:%v\n", err)
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
			service.Logger().Info("Error in receive:%v\n", err)
			return
		}
		service.Logger().Info("Received: %s\n", msg)
	}
}

func (c *Connector) Close() {
	if c.conn == nil {
		return
	}
	defer c.conn.Close()

	// Terminate gracefully...
	service.Logger().Info("Closing all pending connections\n")

	// Close our websocket connection
	err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		service.Logger().Info("Error during closing websocket:%v\n", err)
		return
	}

	select {
	case <-done:
		service.Logger().Info("Receiver Channel Closed! Exiting....\n")
	case <-time.After(time.Duration(1) * time.Second):
		service.Logger().Info("Timeout in closing receiving channel. Exiting....\n")
	}
	return
}

func (c *Connector) SendMessage(message string) {
	c.conn.WriteMessage(websocket.TextMessage, []byte(message))
}

func (c *Connector) StartTest() {
	// Our main loop for the client
	// We send our relevant packets here
	for {
		select {
		case <-time.After(time.Duration(1) * time.Millisecond * 1000):
			// Send an echo packet every second
			err := c.conn.WriteMessage(websocket.TextMessage, []byte("Hello from GolangDocs!"))
			if err != nil {
				service.Logger().Info("Error during writing to websocket:%v\n", err)
				return
			}
		}
	}
}
