package engine

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/popolque/firstbitengi/internal/model"
)

type InputProvider interface {
	ClickerPressed() bool
}

type Game struct {
	engine   *GameEngine
	// renderer *render.Renderer
	// input    *input.InputSystem
	// uiState  *UIAnimationState
	// glitch   *GlitchSystem
}

func NewGame() *Game {
	state := model.NewGameState()
	return &Game{
		engine: NewGameEngine(state),
	}
}

func (g *Game) Update() error {
	// g.input.Poll()
	g.engine.Update(nil) // Pass input provider when ready
	// g.uiState.Animate()
	// g.glitch.Step(g.engine.state.Corruption)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// g.renderer.Draw(screen, g.engine.state, g.uiState, g.glitch)
}

func (g *Game) Layout(outsideW, outsideH int) (int, int) {
	return 1280, 768
}
