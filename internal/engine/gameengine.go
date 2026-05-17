package engine

import (
	"image"
	"math"
	"time"

	"github.com/popolque/firstbitengi/internal/audio"
	"github.com/popolque/firstbitengi/internal/input"
	"github.com/popolque/firstbitengi/internal/model"
	"github.com/popolque/firstbitengi/internal/persist"
	"github.com/popolque/firstbitengi/internal/ui"
)

const GameTickMs = 100.0        // ms per game logic tick
const UpdateMs = 1000.0 / 60.0 // ~16.67ms per Ebitengine Update

type GameEngine struct {
	state      *model.GameState
	accumMs    float64
	autosaver  *persist.Autosaver
	alarmTimer float64
}

func NewGameEngine(state *model.GameState) *GameEngine {
	ge := &GameEngine{
		state:     state,
		autosaver: persist.NewAutosaver("save.json", 30*time.Second),
	}
	ge.autosaver.Start()
	return ge
}

func (ge *GameEngine) Update(in *input.InputSystem) {
	ge.accumMs += UpdateMs
	for ge.accumMs >= GameTickMs {
		ge.accumMs -= GameTickMs
		ge.gameTick(GameTickMs / 1000.0)
	}

	if in != nil {
		ge.handleInputs(in)
	}

	ge.autosaver.MaybeSnapshot(ge.state)
}

func (ge *GameEngine) handleInputs(in *input.InputSystem) {
	if ge.state.RebootPending {
		if in.Clicked {
			rect := image.Rect(ui.WidgetX+50, ui.WidgetY+150, ui.WidgetX+ui.WidgetWidth-50, ui.WidgetY+ui.WidgetHeight-150)
			if in.MousePos.In(rect) {
				if in.MousePos.X < rect.Min.X+rect.Dx()/2 { // Confirm
					ge.Reboot()
				} else { // Abort
					ge.state.RebootPending = false
				}
			}
		}
		return
	}

	if in.Clicked {
		// Tab switching
		if in.MousePos.In(ui.Tab1Rect) {
			ge.state.ActiveTab = "HARDWARE"
			ge.state.ScrollOffset = 0
		} else if in.MousePos.In(ui.Tab2Rect) {
			ge.state.ActiveTab = "UPGRADES"
			ge.state.ScrollOffset = 0
		} else if in.MousePos.In(ui.Tab3Rect) {
			ge.state.ActiveTab = "SYSTEM"
			ge.state.ScrollOffset = 0
		}

		// Tab-specific logic
		switch ge.state.ActiveTab {
		case "HARDWARE":
			if in.MousePos.In(ui.ListRect) {
				rowHeight := 60
				y := in.MousePos.Y - ui.ListRect.Min.Y + ge.state.ScrollOffset
				idx := y / rowHeight
				if idx >= 0 && idx < len(model.AllHardware) {
					ge.PurchaseHardware(model.AllHardware[idx].ID)
				}
			}
		case "UPGRADES":
			if in.MousePos.In(ui.ListRect) {
				rowHeight := 60
				y := in.MousePos.Y - ui.ListRect.Min.Y + ge.state.ScrollOffset
				idx := y / rowHeight
				if idx >= 0 && idx < len(model.AllUpgrades) {
					ge.PurchaseUpgrade(model.AllUpgrades[idx].ID)
				}
			}
		case "SYSTEM":
			if in.ClickerPressed() {
				ge.state.Bits += ge.manualClickValue()
				ge.state.TotalBitsEarned += ge.manualClickValue()
				ge.state.ClickerFlash = true
			}
			if in.RebootTriggered() && ge.state.TotalBitsEarned >= 1_000_000 {
				ge.state.RebootPending = true
			}
		}

		// Scroll logic (works in both list tabs)
		if (ge.state.ActiveTab == "HARDWARE" || ge.state.ActiveTab == "UPGRADES") && in.MousePos.In(ui.ListRect) {
			ge.state.ScrollOffset -= in.ScrollDelta * 20
			if ge.state.ScrollOffset < 0 {
				ge.state.ScrollOffset = 0
			}
		}
	} else if in.ClickerPressed() && ge.state.ActiveTab == "SYSTEM" {
		// Handle Space key for clicker even if not clicking
		ge.state.Bits += ge.manualClickValue()
		ge.state.TotalBitsEarned += ge.manualClickValue()
		ge.state.ClickerFlash = true
	}
}

func (ge *GameEngine) Reboot() {
	gain := math.Log10(ge.state.TotalBitsEarned/1_000_000) * 0.1
	if gain < 0 {
		gain = 0
	}
	ge.state.GHzMultiplier += gain
	ge.state.RebootCount++

	// Reset state
	ge.state.Bits = 0
	ge.state.TotalBitsEarned = 0
	ge.state.Entropy = 0
	ge.state.Corruption = 0
	ge.state.Hardware = make(map[string]int)
	ge.state.Upgrades = make(map[string]bool)
	ge.state.RebootPending = false
}

func (ge *GameEngine) PurchaseHardware(id string) {
	var target model.HardwareDef
	found := false
	for _, h := range model.AllHardware {
		if h.ID == id {
			target = h
			found = true
			break
		}
	}

	if !found {
		return
	}

	owned := ge.state.Hardware[id]
	cost := model.CurrentCost(target, owned)

	if ge.state.Bits >= cost {
		ge.state.Bits -= cost
		ge.state.Hardware[id]++
	}
}

func (ge *GameEngine) PurchaseUpgrade(id string) {
	var target model.UpgradeDef
	found := false
	for _, u := range model.AllUpgrades {
		if u.ID == id {
			target = u
			found = true
			break
		}
	}

	if !found || ge.state.Upgrades[id] {
		return
	}

	if ge.state.Bits >= target.Cost {
		ge.state.Bits -= target.Cost
		ge.state.Upgrades[id] = true
	}
}

func (ge *GameEngine) gameTick(dt float64) {
	// 1. Calculate Multipliers from Upgrades
	prodMult := 1.0
	entropyMult := 1.0
	for id, owned := range ge.state.Upgrades {
		if owned {
			for _, u := range model.AllUpgrades {
				if u.ID == id {
					switch u.Type {
					case model.UpgradeProductionBoost:
						prodMult *= u.Multiplier
					case model.UpgradeEntropyReduction:
						entropyMult *= u.Multiplier
					}
				}
			}
		}
	}

	// 2. Calculate Production
	bps := 0.0
	entropyDelta := 0.0

	for _, def := range model.AllHardware {
		count := ge.state.Hardware[def.ID]
		if count > 0 {
			// Basic production
			bps += float64(count) * def.BaseBPS * prodMult
			// Entropy weight
			entropyDelta += float64(count) * def.EntropyWeight * entropyMult
		}
	}

	// Apply GHz Multiplier
	bps *= ge.state.GHzMultiplier

	// Apply Corruption Penalty
	corruptPenalty := math.Min(ge.state.Corruption/200.0, 0.5)
	bps *= (1.0 - corruptPenalty)

	// Update Bits
	earned := bps * dt
	ge.state.Bits += earned
	ge.state.TotalBitsEarned += earned

	// 2. Entropy & Corruption
	ge.state.Entropy += entropyDelta * dt
	if ge.state.Entropy < 0 {
		ge.state.Entropy = 0
	}
	if ge.state.Entropy > 100 {
		ge.state.Entropy = 100
	}

	if ge.state.Entropy > 50 {
		corruptDelta := (ge.state.Entropy - 50) * 0.002
		ge.state.Corruption += corruptDelta * dt
	}

	// 3. Decay
	if ge.state.Corruption > 75 {
		decayRate := ge.state.Bits * (ge.state.Corruption - 75) * 0.0001
		ge.state.Bits -= decayRate * dt
	}

	// Clamp and housekeeping
	if ge.state.Corruption > 100 {
		ge.state.Corruption = 100
	}
	if ge.state.Bits < 0 {
		ge.state.Bits = 0
	}

	// 4. Audio Triggers
	if ge.state.Corruption > 90 {
		ge.alarmTimer += dt
		if ge.alarmTimer >= 3.0 {
			audio.PlayAlarm()
			ge.alarmTimer = 0
		}
	} else {
		ge.alarmTimer = 0
	}

	ge.state.ClickerFlash = false
}

func (ge *GameEngine) manualClickValue() float64 {
	baseClick := 1.0
	for id, owned := range ge.state.Upgrades {
		if owned {
			for _, u := range model.AllUpgrades {
				if u.ID == id && u.Type == model.UpgradeClickBoost {
					baseClick *= u.Multiplier
				}
			}
		}
	}
	return baseClick * ge.state.GHzMultiplier
}

func (ge *GameEngine) State() *model.GameState {
	return ge.state
}
