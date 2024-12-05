package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"leveling/internal/contract"
	"time"
)

type SkillView struct {
	grid   *tview.Grid
	skill1 *tview.TextView
	skill2 *tview.TextView
	skill3 *tview.TextView
}

type State struct {
	view       *tview.Grid
	healthView *tview.TextView
	app        *tview.Application
	skillView  *SkillView
	targetView *tview.TextView
}

func newState(app *tview.Application) *State {
	healthView := tview.NewTextView()
	healthView.SetTitleAlign(tview.AlignLeft)
	healthView.SetChangedFunc(func() {
		app.Draw()
	})
	skillView := newSkillView()
	targetView := tview.NewTextView()
	targetView.SetTitleAlign(tview.AlignLeft)

	grid := tview.NewGrid().
		SetRows(-1).
		SetColumns(-1, -1, -1).
		SetBorders(true).
		AddItem(healthView, 0, 0, 1, 1, 0, 0, false).
		AddItem(skillView.grid, 0, 1, 1, 1, 0, 0, false).
		AddItem(targetView, 0, 2, 1, 1, 0, 0, false)

	return &State{
		view:       grid,
		healthView: healthView,
		skillView:  skillView,
		targetView: targetView,
		app:        app,
	}
}

func newSkillView() *SkillView {
	skill1 := tview.NewTextView().
		SetText("１").
		SetDynamicColors(true)
	skill2 := tview.NewTextView().
		SetText("２").
		SetDynamicColors(true)
	skill3 := tview.NewTextView().
		SetText("３").
		SetDynamicColors(true)

	grid := tview.NewGrid().
		SetRows(-1, 1).
		SetColumns(2, 2, 2, -1).
		AddItem(skill1, 1, 0, 1, 1, 0, 0, false).
		AddItem(skill2, 1, 1, 1, 1, 0, 0, false).
		AddItem(skill3, 1, 2, 1, 1, 0, 0, false)
	grid.SetBackgroundColor(tcell.ColorRed)

	return &SkillView{
		grid:   grid,
		skill1: skill1,
		skill2: skill2,
		skill3: skill3,
	}
}

func (s *State) UpdateState(event contract.StateChangeEvent) {
	// self health
	s.healthView.SetText(fmt.Sprintf("%v: %d", event.Name, event.Health))
	if event.Damage > 0 {
		s.healthView.SetBorderColor(tcell.ColorRed)
	}
	go func() {
		time.Sleep(100 * time.Millisecond)
		s.healthView.SetBorderColor(tcell.ColorWhite)
		s.app.Draw()
	}()

	// auto attack
	if event.IsAutoAttack {
		s.skillView.skill1.SetBackgroundColor(tcell.ColorRed)
	} else {
		s.skillView.skill1.SetBackgroundColor(tcell.ColorBlack)
	}

	// auto attack
	if event.Action.Id == 2 {
		s.skillView.skill2.SetBackgroundColor(tcell.ColorRed)
	} else {
		s.skillView.skill2.SetBackgroundColor(tcell.ColorBlack)
	}
	// target
	s.targetView.SetText(fmt.Sprintf("%v: %d", event.Target.Name, event.Target.Health))
}
