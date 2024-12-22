package observers

import (
	contract2 "leveling/internal/contract"
	"leveling/internal/server/contract"
)

type PlayerGetHurt struct{}

func (PlayerGetHurt) OnNotify(hero contract.IHero, event contract2.Message) {
	switch event.(type) {
	case contract2.StateChangeEvent:
		// TODO: 實作通知
	}
}
