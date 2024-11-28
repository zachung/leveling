package server

import (
	"leveling/internal/server/contract"
	"leveling/internal/server/message"
	"leveling/internal/server/service"
	"sync"
)

type Round struct {
	isDone bool
	heroes map[*contract.IHero]bool
	keys   map[*message.Client]*contract.IHero
	events chan func()
}

func NewRound(heroes []*contract.IHero) *Round {
	h := map[*contract.IHero]bool{}
	for _, hero := range heroes {
		h[hero] = false
	}

	return &Round{
		heroes: h,
		keys:   make(map[*message.Client]*contract.IHero),
		events: make(chan func()),
	}
}

func (r *Round) round(dt float64) {
	count := len(r.heroes)
	var wg sync.WaitGroup
	countSurvived := count
	var heroes []*contract.IHero
	for h := range r.heroes {
		heroes = append(heroes, h)
	}
	for i, h := range heroes {
		if (*h).IsDie() {
			countSurvived--
			continue
		}
		wg.Add(1)
		go r.attackRound(dt, &wg, heroes, h, i+1)
	}
	wg.Wait()

	for {
		select {
		case event := <-r.events:
			event()
		default:
			return
		}
	}
}

func (r *Round) attackRound(dt float64, wg *sync.WaitGroup, heroes []*contract.IHero, self *contract.IHero, nextInx int) {
	defer wg.Done()

	count := len(heroes)
	if (*self).IsDie() {
		return
	}
	// 選擇攻擊目標
	for {
		if nextInx == count {
			nextInx = 0
		}
		target := heroes[nextInx]
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
	c := (*client).(*message.Client)
	go func() {
		r.events <- func() {
			r.keys[c] = hero
			r.heroes[hero] = false
			service.Logger().Info("hero %d added\n", r.keys[c])
		}
	}()
}

func (r *Round) RemoveHero(client *contract.Client) {
	c := (*client).(*message.Client)
	go func() {
		r.events <- func() {
			hero := r.keys[c]
			delete(r.heroes, hero)
			delete(r.keys, c)
			service.Logger().Info("hero %d leaved\n", hero)
		}
	}()
}
