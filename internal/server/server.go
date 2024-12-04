package server

import (
	"io"
	"leveling/internal/server/connection"
	"leveling/internal/server/contract"
	"leveling/internal/server/entity"
	"leveling/internal/server/repository"
	"leveling/internal/server/service"
	"leveling/internal/server/utils"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const MaxDtMs = 16

type Server struct {
	round    *Round
	lastTime time.Time
	speed    int32
	console  *io.Writer
	stopChan chan bool
}

func NewServer() *contract.Server {
	var server contract.Server
	server = &Server{
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
		}
	}()

	<-s.stopChan
}

func (s *Server) gameInitial() {
	server := contract.Server(s)
	locator := service.GetLocator().
		SetServer(&server).
		SetLogger(service.NewConsole())
	service.Logger().Info("Server initialing\n")

	// listen for client
	go func() {
		hub := connection.NewHub()
		locator.SetHub(hub)
		go (*hub).Run()
		connection.NewMessenger()
		service.Logger().Info("Listening for client\n")
	}()

	var heroes []*contract.IHero
	for _, data := range repository.GetHeroData() {
		heroes = append(heroes, entity.New(data, nil))
	}
	s.round = NewRound(heroes)
}

func (s *Server) gameLoop() {
	// time
	now := utils.Now()
	milliseconds := int32(now.Sub(s.lastTime).Milliseconds())
	dt := milliseconds * s.speed

	defer func() {
		s.lastTime = now
	}()

	for {
		if dt < MaxDtMs {
			time.Sleep(time.Duration(MaxDtMs-dt) * time.Millisecond)
			return
		}
		dt -= MaxDtMs
		s.gameUpdate(float64(MaxDtMs) / 1000)
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

func (s *Server) NewClientConnect(client *contract.Client) *contract.IHero {
	c := *client
	data := repository.GetHeroByName(c.GetName())
	newHero := entity.New(data, client)
	s.round.AddHero(client, newHero)

	return newHero
}

func (s *Server) LeaveClientConnect(client *contract.Client) {
	s.round.RemoveHero(client)
}
