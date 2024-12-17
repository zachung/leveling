package ui

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
	"leveling/internal/client/contract"
	"log"
)

var (
	mplusFaceSource *text.GoTextFaceSource
)

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s
}

type Writer struct {
	screen *ebiten.Image
	text   []byte
}

func NewWriter(screen *ebiten.Image) *Writer {
	return &Writer{screen: screen}
}

func (w *Writer) Write(p []byte) (n int, err error) {
	w.text = append(w.text, p...)

	op := &text.DrawOptions{}
	op.GeoM.Translate(0, 60)
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(w.screen, string(w.text), &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   contract.NormalFontSize,
	}, op)

	return len(p), nil
}
