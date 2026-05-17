package engine

import (
	"github.com/popolque/firstbitengi/internal/model"
)

const GameTickMs = 100.0        // ms per game logic tick
const UpdateMs = 1000.0 / 60.0 // ~16.67ms per Ebitengine Update

type GameEngine struct {
	state   *model.GameState
	accumMs float64
}

func NewGameEngine(state *model.GameState) *GameEngine {
	return &GameEngine{
		state: state,
	}
}

func (ge *GameEngine) Update(in InputProvider) {
	ge.accumMs += UpdateMs
	for ge.accumMs >= GameTickMs {
		ge.accumMs -= GameTickMs
		ge.gameTick(GameTickMs / 1000.0)
	}

	if in != nil && in.ClickerPressed() {
		ge.state.Bits += 1.0 // Manual click value placeholder
		ge.state.ClickerFlash = true
	}
}

func (ge *GameEngine) gameTick(dt float64) {
	// Logic ticks here
	ge.state.ClickerFlash = false
}

func (ge *GameEngine) State() *model.GameState {
	return ge.state
}
