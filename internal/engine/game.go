package engine

import (
	"fmt"
	"leveling/internal/constract"
	"leveling/internal/hero"
	"leveling/internal/repository"
	"leveling/internal/utils"
	"time"
)

type Game struct {
	isFinish bool
	round    *Round
	lastTime time.Time
	speed    int
}

func NewGame() Game {
	return Game{
		isFinish: false,
		lastTime: utils.Now(),
		speed:    4,
	}
}

func (g *Game) Start() {
	fmt.Println("Game initialing")
	g.gameInitial()
	fmt.Println("Game started")
	g.gameStart()
	fmt.Println("Game finished")
}

func (g *Game) gameStart() {
	done := make(chan bool)
	go func() {
		for {
			g.gameLoop()
			if g.isFinish {
				done <- true
			}
		}
	}()
	<-done
}

func (g *Game) gameInitial() {
	var heroes []*constract.IHero
	for _, data := range repository.GetHeroData() {
		heroes = append(heroes, hero.New(data))
	}
	g.round = NewRound(heroes)
}

func (g *Game) gameLoop() {
	// time
	now := utils.Now()
	defer func() {
		g.lastTime = now
	}()

	seconds := now.Sub(g.lastTime).Seconds()

	g.gameUpdate(seconds * float64(g.speed))
	g.gameRender()
}

func (g *Game) gameUpdate(dt float64) {
	g.round.round(dt)
}

func (g *Game) gameRender() {
	// 達到結束條件
	if g.round.IsDone() {
		g.isFinish = true
		return
	}
}
