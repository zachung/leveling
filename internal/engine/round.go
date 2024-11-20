package engine

import "leveling/internal/constract"

type Round struct {
	isDone bool
	heroes []*constract.IHero
}

func NewRound(heroes []*constract.IHero) *Round {
	return &Round{
		isDone: false,
		heroes: heroes,
	}
}

func (r *Round) round(dt float64) {
	for i, h := range r.heroes {
		// FIXME: use goroutine
		// 選擇攻擊目標
		var target *constract.IHero
		for _, otherHero := range append(r.heroes[i+1:], r.heroes[:i]...) {
			if h != otherHero && !(*otherHero).IsDie() {
				target = otherHero
				break
			}
		}
		if target == nil {
			// 找不到目標，遊戲結束
			r.isDone = true
			return
		}
		(*h).Attack(dt, []*constract.IHero{target})
	}
}

func (r *Round) IsDone() bool {
	return r.isDone
}
