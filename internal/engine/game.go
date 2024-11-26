package engine

import (
	"fmt"
	"io"
	"leveling/internal/constract"
	"leveling/internal/hero"
	"leveling/internal/message"
	"leveling/internal/repository"
	"leveling/internal/utils"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const MaxDt = 0.016

type Game struct {
	isFinish bool
	round    *Round
	lastTime time.Time
	speed    int
	console  *io.Writer
	stopChan chan bool
}

func NewGame() *constract.Game {
	var game constract.Game
	game = &Game{
		isFinish: false,
		lastTime: utils.Now(),
		speed:    4,
		stopChan: make(chan bool),
	}

	return &game
}

func (g *Game) Start() {
	g.write("Game initialing\n")
	g.gameInitial()
	g.write("Game started\n")
	g.gameStart()
}

func (g *Game) gameStart() {
	// handle sigint
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		g.Stop()
		os.Exit(1)
	}()

	// start game loop
	go func() {
		for {
			g.gameLoop()
			if g.isFinish {
				g.write("Game finished\n")
				return
			}
		}
	}()

	// listen for client
	go func() {
		g.write("Listening for client\n")
		game := constract.Game(g)
		message.NewMessenger(&game)
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
	g.write("Stopping the application...\n")
	go func() {
		time.Sleep(time.Second)
		g.stopChan <- true
	}()
}

func (g *Game) SetConsole(writer *io.Writer) {
	g.console = writer
}

func (g *Game) write(message string) {
	(*g.console).Write([]byte(message))
}

func (g *Game) Log(format string, args ...any) {
	(*g.console).Write([]byte(fmt.Sprintf(format, args...)))
}
