package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	log "github.com/sirupsen/logrus"
	contract2 "leveling/internal/client/contract"
	"leveling/internal/client/service"
	"leveling/internal/client/ui/object"
	"leveling/internal/contract"
	"leveling/internal/server/utils"
	"time"
)

type World struct {
	event         contract.WorldEvent
	currentTarget string
	heroes        map[string]*contract.Hero
	roles         map[string]*object.Role
}

type ListEntry string

func newWorld() *World {
	w := &World{
		heroes: make(map[string]*contract.Hero),
		roles:  make(map[string]*object.Role),
	}

	service.EventBus().AddObserver(contract2.OnWorldChanged, func() {
		state := service.EventBus().GetWorldState()
		for _, hero := range state.Heroes {
			w.heroes[hero.Name] = &hero
			if w.roles[hero.Name] == nil {
				w.roles[hero.Name] = object.NewRole(hero.Name)
			}
			service.Chat().Info("server %v\n", hero.Position)
		}
	})
	// select hero
	service.EventBus().AddObserver(contract2.OnSelectTarget, func() {
		w.selectNext()
	})

	return w
}

var a float64
var lastUpdate time.Time

func (w *World) Update() {
	now := utils.Now()
	milliseconds := int32(now.Sub(lastUpdate).Milliseconds())
	dt := milliseconds * 1
	if dt < 16 {
		return
	}
	dv := float64(dt) / 1000
	for _, hero := range w.heroes {
		// 在 server 真正回傳實際位置之前，預判位置
		// FIXME: client 預判的位置總會超出 server
		hero.Position[0] += hero.Vector[0] * dv * 160
		hero.Position[1] += hero.Vector[1] * dv * 160
		w.roles[hero.Name].Position = hero.Position

		a += dv
		if hero.Name == "Brian" && a > 1 && (hero.Vector[0] != 0 || hero.Vector[1] != 0) {
			service.Chat().Info("client %v\n", hero.Position)
			a = 0
		}
	}
	lastUpdate = now
}

func (w *World) Draw(dst *ebiten.Image) {
	for _, hero := range w.heroes {
		w.roles[hero.Name].Draw(dst)
	}
}

func (w *World) selectNext() {
	count := len(w.heroes)
	if count == 0 {
		return
	}
	stateEvent := service.EventBus().GetState()
	selfName := stateEvent.Hero.Name
	var curSelect string
	if stateEvent.Hero.Target != nil {
		curSelect = stateEvent.Hero.Target.Name
	}
	isFound := false
	heroes := make([]contract.Hero, 0)
	for _, hero := range w.heroes {
		heroes = append(heroes, *hero)
	}
	i := 0
	r := 0
	for {
		if r == 1 && i >= count {
			break
		}
		if i >= count {
			i = 0
			r = 1
		}
		if curSelect == "" {
			isFound = true
		}
		if !isFound {
			hero := heroes[i]
			if hero.Name == curSelect {
				isFound = true
			}
		} else {
			hero := heroes[i]
			if hero.Name != selfName && hero.Health > 0 {
				selectTarget(hero.Name)
				break
			}
		}
		i++
	}
}

func selectTarget(name string) {
	log.Infof("%v\n", name)
	event := contract.SelectTargetEvent{
		Event: contract.Event{
			Type: contract.SelectTarget,
		},
		Name: name,
	}
	service.Controller().Send(event)
}
