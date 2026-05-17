package model

type GameState struct {
	// Economy
	Bits            float64
	TotalBitsEarned float64

	// System health
	Entropy    float64 // 0.0 – 100.0
	Corruption float64 // 0.0 – 100.0

	// Prestige
	GHzMultiplier float64
	RebootCount   int

	// Hardware owned
	Hardware map[string]int
	Upgrades map[string]bool

	// UI animation signals (not persisted)
	ClickerFlash  bool   // consumed by renderer
	ScrollOffset  int    // hardware panel scroll position
	RebootPending bool   // confirmation dialog visible
	ActiveTab     string // "HARDWARE", "UPGRADES", "SYSTEM"
}

func NewGameState() *GameState {
	gs := &GameState{}
	gs.Sanitize()
	return gs
}

func (gs *GameState) Sanitize() {
	if gs.GHzMultiplier == 0 {
		gs.GHzMultiplier = 1.0
	}
	if gs.Hardware == nil {
		gs.Hardware = make(map[string]int)
	}
	if gs.Upgrades == nil {
		gs.Upgrades = make(map[string]bool)
	}
	if gs.ActiveTab == "" {
		gs.ActiveTab = "HARDWARE"
	}
}
