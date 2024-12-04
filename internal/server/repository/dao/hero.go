package dao

type Hero struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Health   int    `json:"health"`
	Strength int    `json:"strength"`
	MainHand int    `json:"mainHand"`
}
