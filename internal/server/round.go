package server

import (
	"leveling/internal/server/contract"
	"sync"
)

type Round struct {
	isDone bool
	heroes []*contract.IHero
}

func NewRound(heroes []*contract.IHero) *Round {
	return &Round{
		isDone: false,
		heroes: heroes,
	}
}

func (r *Round) round(dt float64) {
	count := len(r.heroes)
	var wg sync.WaitGroup
	countSurvived := count
	defer func() {
		if countSurvived <= 1 {
			r.isDone = true
		}
	}()

	for i, h := range r.heroes {
		if (*h).IsDie() {
			countSurvived--
			continue
		}
		wg.Add(1)
		go r.attackRound(dt, &wg, h, i+1)
	}
	wg.Wait()
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
			//(*self).Attack(dt, []*contract.IHero{target})
			break
		}
		nextInx++
	}
}

func (r *Round) IsDone() bool {
	return r.isDone
}
