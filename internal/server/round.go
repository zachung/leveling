package server

import (
	"leveling/internal/server/contract"
	"leveling/internal/server/service"
	"sync"
)

type Round struct {
	isDone bool
	heroes []*contract.IHero
	keys   map[*contract.Client]int
	events chan func()
}

func NewRound(heroes []*contract.IHero) *Round {
	return &Round{
		heroes: heroes,
		keys:   make(map[*contract.Client]int),
	}
}

func (r *Round) round(dt float64) {
	r.events = make(chan func())
	done := make(chan int)

	count := len(r.heroes)
	var wg sync.WaitGroup
	countSurvived := count
	for i, h := range r.heroes {
		if (*h).IsDie() {
			countSurvived--
			continue
		}
		wg.Add(1)
		go r.attackRound(dt, &wg, h, i+1)
	}
	wg.Wait()

	go func() {
		for {
			select {
			case event := <-r.events:
				event()
			case <-done:
				return
			}
		}
	}()
	done <- 0
}

func (r *Round) attackRound(dt float64, wg *sync.WaitGroup, self *contract.IHero, nextInx int) {
	defer wg.Done()

	count := len(r.heroes)
	if (*self).IsDie() {
		return
	}
	// 選擇攻擊目標
	for {
		if nextInx == count {
			nextInx = 0
		}
		target := r.heroes[nextInx]
		if self == target {
			return
		}
		if !(*target).IsDie() {
			(*self).Attack(dt, []*contract.IHero{target})
			break
		}
		nextInx++
	}
}

func (r *Round) AddHero(client *contract.Client, hero *contract.IHero) {
	r.keys[client] = len(r.heroes)
	r.heroes = append(r.heroes, hero)
}

func (r *Round) RemoveHero(client *contract.Client) {
	r.events <- func() {
		i := r.keys[client]
		r.heroes = append(r.heroes[:i], r.heroes[i+1:]...)
		service.Logger().Info("hero leaved\n")
	}
}
