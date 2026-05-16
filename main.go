package main

import (
	"log"

	"github.com/popolque/firstbitengi/internal/game"
	"github.com/popolque/firstbitengi/internal/models"
	"github.com/popolque/firstbitengi/internal/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	engine   *game.Engine
	renderer *ui.Renderer
}

func (g *Game) Update() error {
	if g.engine.State.Screen == models.ScreenMenu {
		if inpututil.IsKeyJustPressed(ebiten.KeyS) {
			g.engine.State.Screen = models.ScreenInGame
		}
		return nil
	}

	// Mouse Selection
	g.renderer.HandleMouseInput()

	// Debounced inputs for buying and accepting guests
	if inpututil.IsKeyJustPressed(ebiten.Key1) || 
	   inpututil.IsKeyJustPressed(ebiten.Key2) || 
	   inpututil.IsKeyJustPressed(ebiten.Key3) || 
	   inpututil.IsKeyJustPressed(ebiten.Key4) || 
	   inpututil.IsKeyJustPressed(ebiten.KeyB) || 
	   inpututil.IsKeyJustPressed(ebiten.KeyU) || 
	   inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.renderer.HandleInput()
	}

	// Space to end turn
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.engine.EndTurn()
	}

	// Back to menu with ESC
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.engine.State.Screen = models.ScreenMenu
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.renderer.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 600
}

func main() {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Campers - Prototype")

	engine := game.NewEngine()
	renderer := ui.NewRenderer(engine)

	game := &Game{
		engine:   engine,
		renderer: renderer,
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
