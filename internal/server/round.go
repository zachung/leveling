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
	isDone bool
	heroes map[string]contract.IHero
	keys   map[*connection.Client]contract.IHero
	events chan func()
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
	mu := &sync.Mutex{}
	heroes := make(map[string]contract.IHero)
	for _, h := range r.heroes {
		wg.Add(1)
		go func() {
			defer wg.Done()

			if h.Update(dt) {
				mu.Lock()
				heroes[h.GetName()] = h
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	for {
		select {
		case event := <-r.events:
			event()
		default:
			r.broadcastHeroes(heroes)
			return
		}
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
				Hero: contract2.Hero{
					Name:   hero.GetName(),
					Health: hero.GetHealth(),
				},
			}
			hero.Subject().Notify(event)
			service.Logger().Info("%s arrived, current %d.\n", hero.GetName(), len(r.keys))
			r.broadcastHeroes(r.heroes)
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
			client.Close()
			service.Logger().Info("Bye bye %s, now %d.\n", hero.GetName(), len(r.keys))
		}
	}()
}

func (r *Round) broadcastHeroes(iHeroes map[string]contract.IHero) {
	var heroes = make(map[string]contract2.Hero)
	for _, hero := range iHeroes {
		heroes[hero.GetName()] = hero.GetState()
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
