package dao

import "golang.org/x/image/math/f64"

type Hero struct {
	ID       int64    `json:"id"`
	Name     string   `json:"name"`
	Health   int      `json:"health"`
	Strength int      `json:"strength"`
	MainHand int      `json:"mainHand"`
	Position f64.Vec2 `json:"position"`
}
