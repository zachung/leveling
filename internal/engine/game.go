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
	heroes   []*constract.IHero
	lastTime time.Time
	speed    int
}

func NewGame() Game {
	return Game{
		isFinish: false,
		heroes:   make([]*constract.IHero, 0),
		lastTime: utils.Now(),
		speed:    4,
	}
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
	for _, data := range repository.GetHeroData() {
		g.heroes = append(g.heroes, hero.New(data))
	}
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
	// 多個 hero 進入攻擊視野
	(*g.heroes[0]).Attack(dt, g.heroes[1:])
}

func (g *Game) gameRender() {
	heroes := g.heroes
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
