package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	contract2 "leveling/internal/client/contract"
	"leveling/internal/client/service"
	"leveling/internal/contract"
)

type World struct {
	event         contract.WorldEvent
	currentTarget string
}

type ListEntry string

func newWorld() *World {
	return &World{}
}

func selectTarget(name string) {
	event := contract.SelectTargetEvent{
		Event: contract.Event{
			Type: contract.SelectTarget,
		},
		Name: name,
	}
	service.Controller().Send(event)
}

func (w *World) UpdateWorld(event contract.WorldEvent) {
	w.event = event
}

// Focus deprecated
func (w *World) Focus() {
}

func (w *World) SelectNext() {
	if len(w.event.Heroes) == 0 {
		return
	}
	curIndex := 0
	for i, hero := range w.event.Heroes {
		if hero.Name == w.currentTarget {
			curIndex = i
		}
	}
	index := curIndex + 1
	if index >= len(w.event.Heroes) {
		index = 0
	}
	selectTarget(w.event.Heroes[index].Name)
}

func (w *World) Draw(dst *ebiten.Image) {
	event := service.EventBus().GetWorldState()
	heroes := event.Heroes
	if len(heroes) == 0 {
		return
	}
	curName := service.Connector().GetCurName()

	// world map actors
	for _, hero := range event.Heroes {
		x, y := hero.Position[0], hero.Position[1]
		//log.Infof("%v\n", hero.Position)
		name := hero.Name
		var rectClr color.RGBA
		if name == curName {
			rectClr = color.RGBA{100, 100, 100, 255}
		} else {
			rectClr = color.RGBA{255, 100, 100, 255}
		}
		vector.DrawFilledRect(dst, float32(x)-1, float32(y)-1, 20, 20, rectClr, true)
		textOp := &text.DrawOptions{}
		textOp.GeoM.Translate(x, y-20)
		textOp.ColorScale.ScaleWithColor(color.White)
		text.Draw(dst, name, &text.GoTextFace{
			Source: contract2.UiFaceSource,
			Size:   contract2.NormalFontSize,
		}, textOp)
	}
}
