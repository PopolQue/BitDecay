package models

import "github.com/google/uuid"

type ScreenState int

const (
	ScreenMenu ScreenState = iota
	ScreenInGame
)

type GameState struct {
	Screen          ScreenState
	Round           int
	Quarter         int
	ActionPoints    int
	Resources       Resources
	Assets          []Asset
	GuestPool       []NPC
	ActiveGuests    []NPC
	ActiveEvent     *Event
	Board           map[string]int // "q,r" -> OwnerIndex (-1 for unowned)
	SelectedNPCID   uuid.UUID
	SelectedAssetID uuid.UUID
	SelectedTile    string // "q,r"
}
