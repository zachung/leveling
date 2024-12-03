package server

import (
	contract2 "leveling/internal/contract"
	"leveling/internal/server/contract"
	hero2 "leveling/internal/server/hero"
	"leveling/internal/server/message"
	"leveling/internal/server/service"
	"sync"
)

type Round struct {
	isDone       bool
	heroes       map[string]*contract.IHero
	keys         map[*message.Client]*contract.IHero
	events       chan func()
	roundChanged bool
}

func NewRound(heroes []*contract.IHero) *Round {
	h := map[string]*contract.IHero{}
	for _, hero := range heroes {
		h[(*hero).GetName()] = hero
	}

	return &Round{
		heroes: h,
		keys:   make(map[*message.Client]*contract.IHero),
		events: make(chan func()),
	}
}

func (r *Round) round(dt float64) {
	var wg sync.WaitGroup
	for _, h := range r.heroes {
		wg.Add(1)
		go r.updateEntity(dt, &wg, h)
	}
	wg.Wait()

	for {
		select {
		case event := <-r.events:
			event()
		default:
			r.afterRound()
			return
		}
	}
}

func (r *Round) afterRound() {
	if r.roundChanged {
		r.broadcastHeroes()
	}
	r.roundChanged = false
}

func (r *Round) updateEntity(dt float64, wg *sync.WaitGroup, self *contract.IHero) {
	defer wg.Done()

	// 選擇攻擊目標
	if (*self).Update(dt) {
		r.roundChanged = true
	}
}

func (r *Round) AddHero(client *contract.Client, hero *contract.IHero) {
	c := (*client).(*message.Client)
	go func() {
		r.events <- func() {
			h := (*hero).(*hero2.Hero)
			round := contract.Round(r)
			h.SetRound(&round)
			r.keys[c] = hero
			r.heroes[h.GetName()] = hero
			event := contract2.StateChangeEvent{
				Event: contract2.Event{
					Type: contract2.StateChange,
				},
				Name:   h.GetName(),
				Health: h.GetHealth(),
			}
			c.Send(event)
			r.roundChanged = true
			service.Logger().Info("%s arrived, current %d.\n", h.GetName(), len(r.keys))
		}
	}()
}

func (r *Round) RemoveHero(client *contract.Client) {
	c := (*client).(*message.Client)
	go func() {
		r.events <- func() {
			hero := r.keys[c]
			delete(r.heroes, (*hero).GetName())
			delete(r.keys, c)
			r.roundChanged = true
			service.Logger().Info("Bye bye %s, now %d.\n", (*hero).GetName(), len(r.keys))
		}
	}()
}

func (r *Round) broadcastHeroes() {
	var heroes []contract2.Hero
	for _, hero := range r.heroes {
		h := (*hero).(*hero2.Hero)
		heroes = append(heroes, contract2.Hero{h.GetName(), h.GetHealth()})
	}

	event := contract2.WorldEvent{
		Event: contract2.Event{
			Type: contract2.World,
		},
		Heroes: heroes,
	}
	service.Hub().Broadcast(event)
}

func (r *Round) GetHero(name string) *contract.IHero {
	return r.heroes[name]
}
