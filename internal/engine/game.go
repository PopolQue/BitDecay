package engine

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/popolque/firstbitengi/internal/audio"
	"github.com/popolque/firstbitengi/internal/input"
	"github.com/popolque/firstbitengi/internal/model"
	"github.com/popolque/firstbitengi/internal/persist"
	"github.com/popolque/firstbitengi/internal/render"
	"github.com/popolque/firstbitengi/internal/ui"
)

func LoadAssets() error {
	return ui.LoadFont("fonts/BPdotsLight.otf")
}

type Game struct {
	engine   *GameEngine
	renderer *render.Renderer
	input    *input.InputSystem
	// uiState  *UIAnimationState
	// glitch   *GlitchSystem
}

func NewGame() *Game {
	state, err := persist.Load("save.json")
	if err != nil {
		state = model.NewGameState()
	}

	return &Game{
		engine:   NewGameEngine(state),
		renderer: render.NewRenderer(),
		input:    input.NewInputSystem(),
	}
}

func (g *Game) Update() error {
	audio.Update()
	g.input.Poll()
	g.engine.Update(g.input)
	// g.uiState.Animate()
	// g.glitch.Step(g.engine.state.Corruption)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.renderer.Draw(screen, g.engine.State())
}

func (g *Game) Layout(outsideW, outsideH int) (int, int) {
	return 1280, 768
}
