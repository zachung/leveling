package contract

type Message interface {
	Serialize() []byte
}

type Event struct {
	Message `json:"message,omitempty"`
	Type    MessageType `json:"type"`
}

type ActionEvent struct {
	Event `json:"event,omitempty"`
	Id    int `json:"id"`
}

type StateChangeEvent struct {
	Event        `json:"event,omitempty"`
	Name         string `json:"name"`
	Health       int    `json:"health"`
	Damage       int    `json:"damage,omitempty"`
	AttackerName string `json:"attacker_name,omitempty"`
}

type HeroDieEvent struct {
	Event `json:"event,omitempty"`
	Name  string `json:"name"`
}
