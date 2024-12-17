package component

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
	contract2 "leveling/internal/client/contract"
)

type ListBox struct {
	list []string
	x, y float64
}

func NewListBox(x, y float64) *ListBox {
	return &ListBox{
		x: x,
		y: y,
	}
}

func (l *ListBox) Clear() {
	l.list = []string{}
}

func (l *ListBox) AppendItem(item string) {
	l.list = append(l.list, item)
}

func (l *ListBox) Draw(dst *ebiten.Image) {
	for i, str := range l.list {
		textOp := &text.DrawOptions{}
		textOp.GeoM.Translate(l.x, l.y+float64(i*contract2.NormalFontSize))
		textOp.ColorScale.ScaleWithColor(color.White)
		text.Draw(dst, str, &text.GoTextFace{
			Source: contract2.UiFaceSource,
			Size:   contract2.NormalFontSize,
		}, textOp)
	}
}
