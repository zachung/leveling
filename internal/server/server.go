package server

import (
	"io"
	"leveling/internal/hero"
	"leveling/internal/repository"
	"leveling/internal/server/constract"
	"leveling/internal/server/message"
	"leveling/internal/server/service"
	"leveling/internal/utils"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const MaxDt = 0.016

type Server struct {
	isFinish bool
	round    *Round
	lastTime time.Time
	speed    int
	console  *io.Writer
	stopChan chan bool
}

func NewServer() *constract.Server {
	var server constract.Server
	server = &Server{
		isFinish: false,
		lastTime: utils.Now(),
		speed:    1,
		stopChan: make(chan bool),
	}

	return &server
}

func (s *Server) Start() {
	s.gameInitial()
	s.gameStart()
}

func (s *Server) gameStart() {
	service.Logger().Info("Server started\n")
	// handle sigint
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		s.Stop()
		os.Exit(1)
	}()

	// start game loop
	go func() {
		for {
			s.gameLoop()
			if s.isFinish {
				service.Logger().Info("Game finished\n")
				return
			}
		}
	}()

	<-s.stopChan
}

func (s *Server) gameInitial() {
	service.GetLocator().
		SetLogger(service.NewConsole())
	service.Logger().Info("Server initialing\n")

	// listen for client
	go func() {
		service.Logger().Info("Listening for client\n")
		message.NewMessenger()
	}()

	var heroes []*constract.IHero
	for _, data := range repository.GetHeroData() {
		heroes = append(heroes, hero.New(data))
	}
	s.round = NewRound(heroes)
}

func (s *Server) gameLoop() {
	// time
	now := utils.Now()
	seconds := now.Sub(s.lastTime).Seconds()
	dt := seconds * float64(s.speed)

	defer func() {
		s.lastTime = now
	}()

	for {
		// 達到結束條件
		if s.round.IsDone() {
			s.isFinish = true
			return
		}
		var roundDt float64
		if dt > MaxDt {
			dt -= MaxDt
			roundDt = MaxDt
		} else {
			roundDt = dt
		}
		s.gameUpdate(roundDt)
		if roundDt == dt {
			break
		}
	}
}

func (s *Server) gameUpdate(dt float64) {
	s.round.round(dt)
}

func (s *Server) Stop() {
	service.Logger().Info("Stopping the application...\n")
	go func() {
		time.Sleep(time.Second)
		s.stopChan <- true
	}()
}
