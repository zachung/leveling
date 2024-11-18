package engine

import (
	"fmt"
	"leveling/internal/constract"
	"leveling/internal/hero"
	"leveling/internal/weapons"
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
	g.heroes = append(g.heroes, newBrian())
	g.heroes = append(g.heroes, newTaras())
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

func newBrian() *hero.Hero {
	char := &hero.Hero{Name: "Brian", Health: 100, Strength: 6}
	weapon := weapons.NewWeapon(constract.Sword)
	char.Hold(&weapon)

	return char
}

func newTaras() *hero.Hero {
	char := &hero.Hero{Name: "Taras", Health: 100, Strength: 8}
	TarasWeapon := weapons.NewWeapon(constract.Axe)
	char.Hold(&TarasWeapon)

	return char
}
