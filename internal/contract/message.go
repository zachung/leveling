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
	Name         string      `json:"name"`
	Health       int         `json:"health"`
	IsAutoAttack bool        `json:"isAutoAttack"`
	Action       ActionEvent `json:"action,omitempty"`
	Damage       int         `json:"damage,omitempty"`
	Target       Hero        `json:"target,omitempty"`
	Attacker     Hero        `json:"attacker,omitempty"`
}

type HeroDieEvent struct {
	Event `json:"event,omitempty"`
	Name  string `json:"name"`
}

type Hero struct {
	Name   string `json:"name"`
	Health int    `json:"health"`
	Target *Hero  `json:"target,omitempty"`
}

type WorldEvent struct {
	Event  `json:"event,omitempty"`
	Heroes []Hero `json:"heroes,omitempty"`
}

type SelectTargetEvent struct {
	Event `json:"event,omitempty"`
	Name  string `json:"name,omitempty"`
}
