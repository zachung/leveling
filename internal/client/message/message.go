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
	conn *websocket.Conn
}

var done chan interface{}

func NewConnection() *contract.Connector {
	c := &Connector{}
	connector := contract.Connector(c)

	return &connector
}

func (c *Connector) Connect(name string) bool {
	service.Logger().Info("%s connecting...\n", name)
	socketUrl := "ws://localhost:8080" + "/socket"
	header := http.Header{}
	header.Add("Authorization", name)
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, header)
	if err != nil {
		service.Logger().Info("Error connecting to Websocket Server:%v\n", err)
		return false
	}
	c.conn = conn
	go receiveHandler(conn, name)
	service.Logger().Info("Connected!\n")

	return true
}

func receiveHandler(connection *websocket.Conn, name string) {
	for {
		_, msg, err := connection.ReadMessage()
		if err != nil {
			service.Logger().Info("Error in receive:%v\n", err)
			return
		}
		unSerialize := contract2.UnSerialize(msg)
		switch unSerialize.(type) {
		case contract2.GetHurtEvent:
			event := unSerialize.(contract2.GetHurtEvent)
			var message string
			if event.Name == name {
				message = fmt.Sprintf("[red]-%v health[white] from [::u]%s[::U], remain %v\n",
					event.Damage,
					event.AttackerName,
					event.Health,
				)
			} else {
				message = fmt.Sprintf("attack [red]%s(%v)[white] make [red]%v[white] damage\n",
					event.Name,
					event.Health,
					event.Damage,
				)
			}
			service.Logger().Info(message)
		case contract2.HeroDieEvent:
			event := unSerialize.(contract2.HeroDieEvent)
			var message string
			if event.Name == name {
				message = "[red]You Died[white].\n"
			} else {
				message = fmt.Sprintf("%v is Died.\n", event.Name)
			}
			service.Logger().Info(message)
		default:
			service.Logger().Info("Received unknown message: %+v %T\n", unSerialize, unSerialize)
		}
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
				service.Logger().Info("Error during writing to websocket:%v\n", err)
				return
			}
		}
	}
}
