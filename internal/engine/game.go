package engine

import (
	"fmt"
	"leveling/internal/constract"
	"leveling/internal/hero"
	"leveling/internal/weapons"
)

type Game struct {
}

func NewGame() Game {
	return Game{}
}

func (g Game) Start() {
	fmt.Println("Game started")
	g.gameLoop()
	fmt.Println("Game finished")
}

func (g Game) gameLoop() {
	heroes := []*hero.Hero{
		newBrian(),
		newTaras(),
	}
attackRound:
	for {
		iHero := constract.IHero(heroes[1])
		heroes[0].Attack(&iHero)
		for _, h := range heroes {
			if h.IsDie() {
				break attackRound
			}
		}
		heroes = append(heroes[1:], heroes[0])
	}
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
