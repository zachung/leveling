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
	Id       RoleEvent `json:"id"`
	IsEnable bool      `json:"is_cancel"`
}

type MoveEvent struct {
	Event  `json:"event,omitempty"`
	Vector f64.Vec2 `json:"vector"`
}

type StateChangeEvent struct {
	Event        `json:"event,omitempty"`
	Hero         Hero        `json:"hero,omitempty"`
	IsAutoAttack bool        `json:"isAutoAttack"`
	Action       ActionEvent `json:"action,omitempty"`
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
	Vector   f64.Vec2 `json:"vector,omitempty"`
	Target   *Hero    `json:"target,omitempty"`
}

type WorldEvent struct {
	Event  `json:"event,omitempty"`
	Heroes map[string]Hero `json:"heroes,omitempty"`
}

type SelectTargetEvent struct {
	Event `json:"event,omitempty"`
	Name  string `json:"name,omitempty"`
}
