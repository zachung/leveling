package contract

import (
	"encoding/json"
)

type MessageType int

const (
	Unknown MessageType = iota
	StateChange
	HeroDie
	Action
)

func (e Event) GetType() MessageType {
	return e.Type
}

func Serialize(m Message) []byte {
	serializedSpell, _ := json.Marshal(m)
	return serializedSpell
}

func UnSerialize(bytes []byte) Message {
	var t struct {
		Event `json:"event,omitempty"`
	}
	json.Unmarshal(bytes, &t)

	switch t.Type {
	case Action:
		var message ActionEvent
		json.Unmarshal(bytes, &message)
		return message
	case StateChange:
		var message StateChangeEvent
		json.Unmarshal(bytes, &message)
		return message
	case HeroDie:
		var message HeroDieEvent
		json.Unmarshal(bytes, &message)
		return message
	default:
		return nil
	}
}
