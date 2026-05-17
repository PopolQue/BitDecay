package model

import "github.com/google/uuid"

type GameState struct {
	ID            uuid.UUID `json:"id"`
	Bits          float64   `json:"bits"`
	Entropy       float64   `json:"entropy"`
	Corruption    float64   `json:"corruption"`
	ClickerFlash  bool      `json:"-"` // UI state only
}

func NewGameState() *GameState {
	return &GameState{
		ID: uuid.New(),
	}
}
