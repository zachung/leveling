package server

import (
	contract2 "leveling/internal/contract"
	"leveling/internal/server/connection"
	"leveling/internal/server/contract"
	hero2 "leveling/internal/server/entity"
	"leveling/internal/server/observers"
	"leveling/internal/server/repository"
	"leveling/internal/server/service"
	"sync"
)

type Round struct {
	isDone       bool
	heroes       map[string]contract.IHero
	keys         map[*connection.Client]contract.IHero
	events       chan func()
	roundChanged bool
}

func NewRound(heroes []contract.IHero) *Round {
	r := &Round{
		heroes: map[string]contract.IHero{},
		keys:   make(map[*connection.Client]contract.IHero),
		events: make(chan func()),
	}

	for _, hero := range heroes {
		hero.SetRound(r)
		r.heroes[hero.GetName()] = hero
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

func (r *Round) updateEntity(dt float64, wg *sync.WaitGroup, self contract.IHero) {
	defer wg.Done()

	// 選擇攻擊目標
	if self.Update(dt) {
		r.roundChanged = true
	}
}

func (r *Round) AddHero(client contract.Client) contract.IHero {
	c := client.(*connection.Client)
	data := repository.GetHeroByName(client.GetName())
	// observer
	subject := contract.Subject(&hero2.RoleSubject{})
	hurt := contract.Observer(observers.NewPlayerListener(client))
	subject.AddObserver(hurt)
	hero := hero2.NewRole(data, subject, client)

	go func() {
		r.events <- func() {
			// setup hero
			hero.SetSubject(subject)
			hero.SetRound(r)

			r.keys[c] = hero
			r.heroes[hero.GetName()] = hero
			event := contract2.StateChangeEvent{
				Event: contract2.Event{
					Type: contract2.StateChange,
				},
				Name:   hero.GetName(),
				Health: hero.GetHealth(),
			}
			hero.Subject().Notify(event)
			r.roundChanged = true
			service.Logger().Info("%s arrived, current %d.\n", hero.GetName(), len(r.keys))
		}
	}()

	return hero
}

func (r *Round) RemoveHero(client contract.Client) {
	c := client.(*connection.Client)
	go func() {
		r.events <- func() {
			hero := r.keys[c]
			delete(r.heroes, hero.GetName())
			delete(r.keys, c)
			r.roundChanged = true
			client.Close()
			service.Logger().Info("Bye bye %s, now %d.\n", hero.GetName(), len(r.keys))
		}
	}()
}

func (r *Round) broadcastHeroes() {
	var heroes []contract2.Hero
	for _, hero := range r.heroes {
		heroes = append(heroes, hero.GetState())
	}

	event := contract2.WorldEvent{
		Event: contract2.Event{
			Type: contract2.World,
		},
		Heroes: heroes,
	}
	service.Hub().Broadcast(event)
}

func (r *Round) GetHero(name string) contract.IHero {
	return r.heroes[name]
}
