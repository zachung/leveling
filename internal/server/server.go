package server

import (
	"fmt"
	"io"
	"leveling/internal/constract"
	"leveling/internal/hero"
	"leveling/internal/repository"
	"leveling/internal/server/message"
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
		speed:    4,
		stopChan: make(chan bool),
	}

	return &server
}

func (s *Server) Start() {
	s.write("Server initialing\n")
	s.gameInitial()
	s.write("Server started\n")
	s.gameStart()
}

func (s *Server) gameStart() {
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
				s.write("Game finished\n")
				return
			}
		}
	}()

	// listen for client
	go func() {
		s.write("Listening for client\n")
		server := constract.Server(s)
		message.NewMessenger(&server)
	}()

	<-s.stopChan
}

func (s *Server) gameInitial() {
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
	s.write("Stopping the application...\n")
	go func() {
		time.Sleep(time.Second)
		s.stopChan <- true
	}()
}

func (s *Server) SetConsole(writer *io.Writer) {
	s.console = writer
}

func (s *Server) write(message string) {
	(*s.console).Write([]byte(message))
}

func (s *Server) Log(format string, args ...any) {
	(*s.console).Write([]byte(fmt.Sprintf(format, args...)))
}
