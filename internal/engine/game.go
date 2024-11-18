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
	heroes   []*hero.Hero
}

func NewGame() Game {
	return Game{false, make([]*hero.Hero, 0)}
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
	}
	for _, data := range heroData {
		g.heroes = append(g.heroes, newHero(data))
	}
}

func (g *Game) gameLoop() {
	heroes := g.heroes
	iHero := constract.IHero(heroes[1])
	heroes[0].Attack(&iHero)
	for _, h := range heroes {
		if h.IsDie() {
			g.isFinish = true
			return
		}
	}
	g.heroes = append(heroes[1:], heroes[0])
}

func newHero(s string) *hero.Hero {
	data := entity.Hero{}
	err := json.Unmarshal([]byte(s), &data)
	if err != nil {
		return nil
	}
	char := hero.New(data)

	return char
}
