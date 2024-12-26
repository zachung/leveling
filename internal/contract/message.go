package contract

import "golang.org/x/image/math/f64"

type Message interface {
	Serialize() []byte
}

type Event struct {
	Message `json:"message,omitempty"`
	Type    MessageType `json:"type"`
}

type ActionEvent struct {
	Event    `json:"event,omitempty"`
	Id       KeyFunc `json:"id"`
	IsCancel bool    `json:"is_cancel"`
}

type MoveEvent struct {
	Event  `json:"event,omitempty"`
	Vector f64.Vec2 `json:"vector"`
}

type StateChangeEvent struct {
	Event        `json:"event,omitempty"`
	Name         string      `json:"name"`
	Health       int         `json:"health"`
	IsAutoAttack bool        `json:"isAutoAttack"`
	Position     f64.Vec2    `json:"position"`
	Action       ActionEvent `json:"action,omitempty"`
	Damage       int         `json:"damage,omitempty"`
	Target       Hero        `json:"target,omitempty"`
	Attacker     Hero        `json:"attacker,omitempty"`
}

type GetHurtEvent struct {
	Event  `json:"event,omitempty"`
	Damage int  `json:"damage,omitempty"`
	From   Hero `json:"from,omitempty"`
}

type MakeDamageEvent struct {
	Event  `json:"event,omitempty"`
	Damage int  `json:"damage,omitempty"`
	To     Hero `json:"to,omitempty"`
}

type HeroDieEvent struct {
	Event `json:"event,omitempty"`
	Name  string `json:"name"`
}

type Hero struct {
	Name     string   `json:"name"`
	Health   int      `json:"health"`
	Position f64.Vec2 `json:"position,omitempty"`
	Target   *Hero    `json:"target,omitempty"`
}

type WorldEvent struct {
	Event  `json:"event,omitempty"`
	Heroes []Hero `json:"heroes,omitempty"`
}

type SelectTargetEvent struct {
	Event `json:"event,omitempty"`
	Name  string `json:"name,omitempty"`
}
