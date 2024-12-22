package ui

import (
	"fmt"
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"leveling/internal/client/ui/keys"
)

type Game struct {
	state *State
	world *World
	ui    *ebitenui.UI
}

func NewGame() *Game {
	container := layoutRoot()
	ui := ebitenui.UI{
		Container: container,
	}
	game := Game{
		ui: &ui,
	}
	game.state = newState()
	game.world = newWorld()
	container.AddChild(game.world.list)

	return &game
}

func (g *Game) Update() error {
	// update the UI
	g.ui.Update()
	keys.NewRune(keys.NewSwitchTarget(nil)).Execute()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.ui.Draw(screen)

	g.state.Draw(screen)
	g.world.Draw(screen)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()), 0, 0)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
