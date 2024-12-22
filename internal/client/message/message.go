package message

import (
	"fmt"
	"github.com/gorilla/websocket"
	"leveling/internal/client/contract"
	"leveling/internal/client/service"
	contract2 "leveling/internal/contract"
	"net/http"
	"time"
)

type Connector struct {
	conn    *websocket.Conn
	curName string
}

var done chan interface{}

func NewConnection() contract.Connector {
	return new(Connector)
}

func (c *Connector) Connect(name string) bool {
	service.Chat().Info("%s connecting...\n", name)
	socketUrl := "ws://localhost:8080" + "/socket"
	header := http.Header{}
	header.Add("Authorization", name)
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, header)
	if err != nil {
		service.Chat().Info("Error connecting to Websocket Server:%v\n", err)
		return false
	}
	c.conn = conn
	c.curName = name
	go receiveHandler(conn, name)
	service.Chat().Info("Connected!\n")

	return true
}

func receiveHandler(connection *websocket.Conn, name string) {
	for {
		_, msg, err := connection.ReadMessage()
		if err != nil {
			service.Chat().Info("Error in receive:%v\n", err)
			return
		}
		unSerialize := contract2.UnSerialize(msg)
		switch unSerialize.(type) {
		case contract2.StateChangeEvent:
			event := unSerialize.(contract2.StateChangeEvent)
			var message string
			if event.Name == name {
				if event.Attacker.Name != "" {
					message = fmt.Sprintf("[color=ff0000]-%v health[/color] from %s, remain %v\n",
						event.Damage,
						event.Attacker.Name,
						event.Health,
					)
				}
				service.EventBus().SetState(event)
			} else {
				message = fmt.Sprintf("attack [color=ff0000]%s(%v)[/color] make [color=ff0000]%v[/color] damage\n",
					event.Name,
					event.Health,
					event.Damage,
				)
			}
			if message != "" {
				service.Chat().Info(message)
			}
		case contract2.HeroDieEvent:
			event := unSerialize.(contract2.HeroDieEvent)
			var message string
			if event.Name == name {
				message = "[color=ff0000]You Died[/color].\n"
				service.EventBus().SetState(contract2.StateChangeEvent{Name: name, Health: 0})
			} else {
				message = fmt.Sprintf("%v is Died.\n", event.Name)
			}
			service.Chat().Info(message)
		case contract2.WorldEvent:
			event := unSerialize.(contract2.WorldEvent)
			service.EventBus().SetWorldState(event)
		default:
			service.Chat().Info("Received unknown message: %+v %T\n", unSerialize, unSerialize)
		}
	}
}

func (c *Connector) Close() {
	if c.conn == nil {
		return
	}
	defer c.conn.Close()

	// Terminate gracefully...
	service.Chat().Info("Closing all pending connections\n")

	// Close our websocket connection
	err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		service.Chat().Info("Error during closing websocket:%v\n", err)
		return
	}

	select {
	case <-done:
		service.Chat().Info("Receiver Channel Closed! Exiting....\n")
	case <-time.After(time.Duration(1) * time.Second):
		service.Chat().Info("Timeout in closing receiving channel. Exiting....\n")
	}
	return
}

func (c *Connector) SendMessage(message contract2.Message) {
	if c.conn == nil {
		return
	}
	c.conn.WriteMessage(websocket.TextMessage, contract2.Serialize(message))
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
				service.Chat().Info("Error during writing to websocket:%v\n", err)
				return
			}
		}
	}
}

func (c *Connector) GetCurName() string {
	return c.curName
}
