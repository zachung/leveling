package message

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type Connection struct {
	conn *websocket.Conn
}

var done chan interface{}

func NewConnection() *Connection {
	conn := connect()
	go receiveHandler(conn)

	return &Connection{conn}
}

func connect() *websocket.Conn {
	socketUrl := "ws://localhost:8080" + "/socket"
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err)
	}

	return conn
}

func receiveHandler(connection *websocket.Conn) {
	defer close(done)
	for {
		_, msg, err := connection.ReadMessage()
		if err != nil {
			log.Println("Error in receive:", err)
			return
		}
		log.Printf("Received: %s\n", msg)
	}
}

func (c *Connection) Close() {
	defer c.conn.Close()

	// Terminate gracefully...
	log.Println("Received SIGINT interrupt signal. Closing all pending connections")

	// Close our websocket connection
	err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("Error during closing websocket:", err)
		return
	}

	select {
	case <-done:
		log.Println("Receiver Channel Closed! Exiting....")
	case <-time.After(time.Duration(1) * time.Second):
		log.Println("Timeout in closing receiving channel. Exiting....")
	}
	return
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
				log.Println("Error during writing to websocket:", err)
				return
			}
		}
	}
}

func (c Connection) SendMessage(message string) {
	c.conn.WriteMessage(websocket.TextMessage, []byte(message))
}
