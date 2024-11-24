package engine

import (
	"leveling/internal/constract"
	"leveling/internal/hero"
	"leveling/internal/repository"
	"leveling/internal/ui"
	"leveling/internal/utils"
	"time"
)

const MaxDt = 0.016

type Game struct {
	isFinish bool
	round    *Round
	lastTime time.Time
	speed    int
	ui       *ui.UI
	stopChan chan bool
}

func NewGame() *constract.Game {
	game := &Game{
		isFinish: false,
		lastTime: utils.Now(),
		speed:    4,
		stopChan: make(chan bool),
	}
	igame := constract.Game(game)
	game.ui = ui.NewUi(&igame)

	return &igame
}

func (g *Game) Start() {
	g.ui.Run()
	ui.Logger().Info("Game initialing")
	g.gameInitial()
	ui.Logger().Info("Game started")
	g.gameStart()
}

func (g *Game) gameStart() {
	go func() {
		for {
			g.gameLoop()
			if g.isFinish {
				ui.Logger().Info("Game finished")
				return
			}
		}
	}()
	<-g.stopChan
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

func (g *Game) Stop() {
	ui.Logger().Info("Stopping the application...")
	go func() {
		time.Sleep(time.Second)
		g.ui.Stop()
		g.stopChan <- true
	}()
}

func (g *Game) UI() *constract.UI {
	uii := constract.UI(g.ui)

	return &uii
}
