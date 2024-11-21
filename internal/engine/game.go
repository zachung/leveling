package engine

import (
	"fmt"
	"leveling/internal/constract"
	"leveling/internal/hero"
	"leveling/internal/repository"
	"leveling/internal/utils"
	"time"
)

const MaxDt = 0.016

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
		speed:    400000,
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
				return
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
	seconds := now.Sub(g.lastTime).Seconds()
	dt := seconds * float64(g.speed)

	defer func() {
		g.lastTime = now
	}()

	for {
		// 達到結束條件
		if g.round.IsDone() {
			g.isFinish = true
			return
		}
		var roundDt float64
		if dt > MaxDt {
			dt -= MaxDt
			roundDt = MaxDt
		} else {
			roundDt = dt
		}
		g.gameUpdate(roundDt)
		if roundDt == dt {
			break
		}
	}
}

func (g *Game) gameUpdate(dt float64) {
	g.round.round(dt)
}
