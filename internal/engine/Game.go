package engine

import (
	"fmt"
	"leveling/internal/hero"
)

type Game struct {
}

func NewGame() Game {
	return Game{}
}

func (g Game) Start() {
	fmt.Println("Game started")
	g.gameLoop()
}

func (g Game) gameLoop() {
	brian := &hero.Hero{Name: "brian", Health: 100, AttackType: hero.StraightSword, Strength: 6}
	taras := &hero.Hero{Name: "taras", Health: 100, AttackType: hero.Claw, Strength: 8}
	heroes := []*hero.Hero{brian, taras}
attackRound:
	for {
		heroes[0].Attack(heroes[1])
		for _, h := range heroes {
			if h.IsDie() {
				break attackRound
			}
		}
		heroes = append(heroes[1:], heroes[0])
	}
	//fmt.Printf("The winner is %s\n", brian.Name)
	fmt.Println("Game finished")
}
