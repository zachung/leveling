package ui

import (
	"fmt"
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"leveling/internal/client/contract"
	"leveling/internal/client/service"
	"leveling/internal/client/ui/keys"
)

type Game struct {
	state    *State
	worldMap *World
	ui       *ebitenui.UI
	world    *ebiten.Image
	camera   Camera
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
	game.worldMap = newWorld()
	game.world = ebiten.NewImage(contract.ScreenWidth*2, contract.ScreenHeight*2)

	return &game
}

var keyHandler = keys.NewMove(keys.NewAction(keys.NewSwitchTarget(nil)))

func (g *Game) Update() error {
	// update the UI
	g.ui.Update()
	g.worldMap.Update()
	keyHandler.Execute()
	name := service.EventBus().GetState().Hero.Name
	if g.worldMap.heroes[name] != nil {
		g.camera.Position = g.worldMap.heroes[name].Position
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		if g.camera.ZoomFactor > -2400 {
			g.camera.ZoomFactor -= 1
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		if g.camera.ZoomFactor < 2400 {
			g.camera.ZoomFactor += 1
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.camera.Reset()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.world.Clear()
	g.worldMap.Draw(g.world)
	g.camera.Render(g.world, screen)

	g.ui.Draw(screen)
	g.state.Draw(screen)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))

	worldX, worldY := g.camera.ScreenToWorld(ebiten.CursorPosition())
	ebitenutil.DebugPrintAt(
		screen,
		fmt.Sprintf("%s\nCursor World Pos: %.2f,%.2f",
			g.camera.String(),
			worldX, worldY),
		0, contract.ScreenHeight-32,
	)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
