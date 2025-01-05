package contract

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
	"log"
)

var UiFaceSource *text.GoTextFaceSource
var UiTextFace *text.GoTextFace

const (
	ScreenWidth    = 640
	ScreenHeight   = 480
	NormalFontSize = 16
)

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal(err)
	}
	UiFaceSource = s

	UiTextFace = &text.GoTextFace{
		Source: s,
		Size:   float64(NormalFontSize),
	}
}
