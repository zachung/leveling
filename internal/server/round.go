package server

import (
	contract2 "leveling/internal/contract"
	"leveling/internal/server/connection"
	"leveling/internal/server/contract"
	hero2 "leveling/internal/server/entity"
	"leveling/internal/server/observers"
	"leveling/internal/server/service"
	"sync"
)

type Round struct {
	isDone       bool
	heroes       map[string]*contract.IHero
	keys         map[*connection.Client]*contract.IHero
	events       chan func()
	roundChanged bool
}

func NewRound(heroes []*contract.IHero) *Round {
	r := &Round{
		heroes: map[string]*contract.IHero{},
		keys:   make(map[*connection.Client]*contract.IHero),
		events: make(chan func()),
	}
	round := contract.Round(r)

	subject := hero2.NewRoleSubject()
	hurt := contract.Observer(observers.EnemyGetHurt{})
	subject.AddObserver(hurt)
	for _, hero := range heroes {
		iHero := *hero
		iHero.SetSubject(subject)
		iHero.SetRound(&round)
		r.heroes[iHero.GetName()] = hero
	}

	return r
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
	c := (*client).(*connection.Client)
	go func() {
		r.events <- func() {
			round := contract.Round(r)

			// observer
			subject := contract.Subject(&hero2.RoleSubject{})
			hurt := contract.Observer(observers.PlayerGetHurt{})
			subject.AddObserver(hurt)
			// setup hero
			h := (*hero).(*hero2.Hero)
			h.SetSubject(subject)
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
			h.Subject().Notify(h, event)
			r.roundChanged = true
			service.Logger().Info("%s arrived, current %d.\n", h.GetName(), len(r.keys))
		}
	}()
}

func (r *Round) RemoveHero(client *contract.Client) {
	c := (*client).(*connection.Client)
	go func() {
		r.events <- func() {
			hero := r.keys[c]
			delete(r.heroes, (*hero).GetName())
			delete(r.keys, c)
			r.roundChanged = true
			(*client).Close()
			service.Logger().Info("Bye bye %s, now %d.\n", (*hero).GetName(), len(r.keys))
		}
	}()
}

func (r *Round) broadcastHeroes() {
	var heroes []contract2.Hero
	for _, hero := range r.heroes {
		h := (*hero).(*hero2.Hero)
		elems := contract2.Hero{
			Name:   h.GetName(),
			Health: h.GetHealth(),
		}
		target := h.GetTarget()
		if target != nil {
			elems.Target = &contract2.Hero{
				Name:   (*target).GetName(),
				Health: (*target).GetHealth(),
			}
		}
		heroes = append(heroes, elems)
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
