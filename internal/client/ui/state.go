package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	contract2 "leveling/internal/client/contract"
	"leveling/internal/client/service"
	"leveling/internal/contract"
	"strconv"
)

var (
	skillBoxSize     = float32(40)
	skillStrokeWidth = float32(2)
	colorRed         = color.RGBA{R: 0xff, A: 0xff}
)

type State struct {
	skills    []*SkillBox
	isWarning bool
	event     *contract.StateChangeEvent
}

func newState() *State {
	const skillCount = 10
	x := (screenWidth - (skillBoxSize+skillStrokeWidth*2)*skillCount) / 2
	y := screenHeight - skillBoxSize - skillStrokeWidth*2

	var skills []*SkillBox
	for i := 0; i < skillCount; i++ {
		skills = append(skills, &SkillBox{
			x:    x + (skillBoxSize+skillStrokeWidth)*float32(i) + skillStrokeWidth,
			y:    y,
			text: strconv.Itoa(i),
		})
	}
	s := &State{
		skills: skills,
	}

	return s
}

func (s *State) UpdateState(event contract.StateChangeEvent) {
	s.event = &event
}

func (s *State) Draw(dst *ebiten.Image) {
	event := service.EventBus().GetState()

	// auto attack
	if event.IsAutoAttack {
		s.skills[0].isEnabled = true
	} else {
		s.skills[0].isEnabled = false
	}
	// skill 1
	if event.Action.Id == 2 {
		s.skills[1].isEnabled = true
	} else {
		s.skills[1].isEnabled = false
	}

	for _, skill := range s.skills {
		drawSkillBox(dst, skill)
	}
}

type SkillBox struct {
	x, y      float32
	text      string
	isEnabled bool
}

func drawSkillBox(dst *ebiten.Image, box *SkillBox) {
	var clr color.Color
	if box.isEnabled {
		clr = colorRed
	} else {
		clr = color.White
	}
	vector.StrokeRect(dst, box.x, box.y, skillBoxSize, skillBoxSize, skillStrokeWidth, clr, false)
	textOp := &text.DrawOptions{}
	// half number align to center
	textOp.GeoM.Translate(float64(box.x+skillBoxSize/4), float64(box.y))
	textOp.ColorScale.ScaleWithColor(color.White)
	text.Draw(dst, box.text, &text.GoTextFace{
		Source: contract2.UiFaceSource,
		Size:   float64(skillBoxSize),
	}, textOp)
}
