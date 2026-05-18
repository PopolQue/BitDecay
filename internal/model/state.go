package model

import "math"

type GameState struct {
	// Economy
	Bits            float64
	TotalBitsEarned float64

	// System health
	Entropy    float64 // 0.0 – 100.0
	Corruption float64 // 0.0 – 100.0

	// Infrastructure State (New)
	PowerUsage    float64
	PowerCapacity float64
	HeatLevel     float64 // 0.0 - 100.0 (based on gen vs cooling)
	SpaceUsage    float64
	SpaceCapacity float64

	// Prestige
	GHzMultiplier float64
	RebootCount   int

	// Hardware owned
	Hardware map[string]int
	Upgrades map[string]bool

	// UI animation signals (not persisted)
	ClickerFlash  bool     // consumed by renderer
	ScrollOffset  int      // hardware panel scroll position
	RebootPending bool     // confirmation dialog visible
	MessageLog    []string // System log messages
	PacketActive  bool     // Is a random packet available to intercept?
	PacketTimer   float64  // Time remaining for the packet
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
	// Initial Rack Space so the player can start building
	if _, ok := gs.Hardware["rack_shelf"]; !ok && len(gs.Hardware) == 0 {
		gs.Hardware["rack_shelf"] = 1
	}
	if gs.Upgrades == nil {
		gs.Upgrades = make(map[string]bool)
	}
	if gs.MessageLog == nil {
		gs.MessageLog = []string{"[SYSTEM] CORE_INIT_SUCCESS"}
	}
}

func (gs *GameState) LogMessage(msg string) {
	gs.MessageLog = append([]string{msg}, gs.MessageLog...)
	if len(gs.MessageLog) > 9{ // we show only 9 messages because 10 would overlap the LogWidget in the unified UI
		gs.MessageLog = gs.MessageLog[:9]
	}
}

func (gs *GameState) GetRank() string {
	val := gs.TotalBitsEarned / 8
	switch {
	case val < 1024:
		return "NOVICE"
	case val < 1024*1024:
		return "TECHNICIAN"
	case val < 1024*1024*1024:
		return "ENGINEER"
	case val < 1024*1024*1024*1024:
		return "ARCHITECT"
	default:
		return "OVERSEER"
	}
}

func (gs *GameState) GetRebootThreshold() float64 {
	base := 10000.0
	scaling := 1.5 // Threshold increases by 1.5x each reboot
	return base * math.Pow(scaling, float64(gs.RebootCount))
}
