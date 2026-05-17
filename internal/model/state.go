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
	return &GameState{
		GHzMultiplier: 1.0,
		Hardware:      make(map[string]int),
		Upgrades:      make(map[string]bool),
		ActiveTab:     "HARDWARE",
	}
}
