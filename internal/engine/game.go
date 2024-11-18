package engine

import (
	"encoding/json"
	"fmt"
	"leveling/internal/constract"
	"leveling/internal/entity"
	"leveling/internal/hero"
)

type Game struct {
	isFinish bool
	heroes   []*constract.IHero
}

func NewGame() Game {
	return Game{false, make([]*constract.IHero, 0)}
}

func (g *Game) IsFinish() bool {
	return g.isFinish
}

func (g *Game) Start() {
	fmt.Println("Game initialing")
	g.gameInitial()
	fmt.Println("Game started")
	for {
		g.gameLoop()
		if g.isFinish {
			break
		}
	}
	fmt.Println("Game finished")
}

func (g *Game) gameInitial() {
	heroData := []string{
		`{"name": "Brian", "Health": 100, "Strength": 6, "mainHand": 0}`,
		`{"name": "Taras", "Health": 100, "Strength": 8, "mainHand": 2}`,
		`{"name": "Sin", "Health": 100, "Strength": 8, "mainHand": 1}`,
	}
	for _, data := range heroData {
		g.heroes = append(g.heroes, newHero(data))
	}
}

func (g *Game) gameLoop() {
	heroes := g.heroes
	// 多個 hero 進入攻擊視野
	(*heroes[0]).Attack(g.heroes[1:])
	// 輪番檢查死亡狀態
	for i := len(heroes) - 1; i >= 0; i-- {
		if (*heroes[i]).IsDie() {
			heroes = append(heroes[:i], heroes[i+1:]...)
		}
	}
	// 達到結束條件
	if len(heroes) <= 1 {
		g.isFinish = true
		return
	}
	// 輪到下一個 hero
	g.heroes = append(heroes[1:], heroes[0])
}

func newHero(s string) *constract.IHero {
	data := entity.Hero{}
	err := json.Unmarshal([]byte(s), &data)
	if err != nil {
		return nil
	}
	iHero := constract.IHero(hero.New(data))

	return &iHero
}
