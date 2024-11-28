package contract

import "encoding/json"

type Message interface {
	Serialize() []byte
}

type Action struct {
	Id int `json:"id"`
}

func (s Action) Serialize() []byte {
	serializedSpell, _ := json.Marshal(s)
	return serializedSpell
}

func UnSerialize(bytes []byte) Action {
	var action Action
	json.Unmarshal(bytes, &action)

	return action
}
