package server

import (
	"io"
	"leveling/internal/server/connection"
	"leveling/internal/server/contract"
	"leveling/internal/server/entity"
	"leveling/internal/server/observers"
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

func NewServer() contract.Server {
	server := &Server{
		lastTime: utils.Now(),
		speed:    1,
		stopChan: make(chan bool),
	}

	return server
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
	locator := service.GetLocator().
		SetServer(s).
		SetLogger(service.NewConsole())
	service.Logger().Info("Server initialing\n")

	// listen for client
	go func() {
		hub := connection.NewHub()
		locator.SetHub(hub)
		go hub.Run()
		connection.NewMessenger()
		service.Logger().Info("Listening for client\n")
	}()

	var heroes []contract.IHero
	for _, data := range repository.GetHeroData() {
		subject := entity.NewRoleSubject()
		hero := entity.NewRole(data, subject, nil)
		hurt := contract.Observer(observers.NewEnemyListener(hero))
		subject.AddObserver(hurt)
		heroes = append(heroes, hero)
	}
	s.round = NewRound(heroes)
}

func (s *Server) gameLoop() {
	// time
	now := utils.Now()
	milliseconds := int32(now.Sub(s.lastTime).Milliseconds())
	dt := milliseconds * s.speed

	defer func() {
		s.lastTime = now.Add(-time.Duration(dt) * time.Millisecond)
	}()

	for {
		if dt < MaxDtMs {
			time.Sleep(time.Duration(MaxDtMs-dt) * time.Millisecond)
			return
		}
		dt -= MaxDtMs
		s.round.round(float64(MaxDtMs) / 1000)
	}
}

func (s *Server) Stop() {
	service.Logger().Info("Stopping the application...\n")
	go func() {
		time.Sleep(time.Second)
		s.stopChan <- true
	}()
}

func (s *Server) NewClientConnect(client contract.Client) contract.IHero {
	return s.round.AddHero(client)
}

func (s *Server) LeaveClientConnect(client contract.Client) {
	s.round.RemoveHero(client)
}
