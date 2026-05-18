package engine

import (
	"fmt"
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
			rect := image.Rect(ui.WidgetX+100, ui.WidgetY+100, ui.WidgetX+ui.WidgetWidth-100, ui.WidgetY+ui.WidgetHeight-100)
			if in.MousePos.In(rect) {
				audio.PlayClick()
				if in.MousePos.X < rect.Min.X+rect.Dx()/2 {
					ge.Reboot()
				} else {
					ge.state.RebootPending = false
				}
			}
		}
		return
	}

	if ge.state.PacketActive && in.Clicked && in.MousePos.In(ui.PacketRect) {
		reward := math.Max(1024, ge.state.Bits*0.1)
		ge.state.Bits += reward
		ge.state.TotalBitsEarned += reward
		ge.state.PacketActive = false
		ge.state.LogMessage(fmt.Sprintf("[REWARD] PACKET_HARVEST: +%.0f bits", reward))
		audio.PlayClick()
		return
	}

	if in.Clicked {
		if in.MousePos.In(ui.ClickerRegion) {
			ge.PerformManualClick()
		}

		if in.MousePos.In(ui.ListRect) {
			rowHeight := 50
			y := in.MousePos.Y - ui.ListRect.Min.Y + ge.state.ScrollOffset
			idx := y / rowHeight
			if idx >= 0 && idx < len(model.AllHardware) {
				ge.PurchaseHardware(model.AllHardware[idx].ID)
			}
		}

		if in.MousePos.In(ui.RebootBtnRect) && ge.state.TotalBitsEarned >= 1_000_000 {
			ge.state.RebootPending = true
			audio.PlayClick()
		}

		if in.MousePos.In(ui.ListRect) {
			ge.state.ScrollOffset -= in.ScrollDelta * 20
			if ge.state.ScrollOffset < 0 {
				ge.state.ScrollOffset = 0
			}
		}
	} else if in.ClickerPressed() {
		ge.PerformManualClick()
	}
}

func (ge *GameEngine) PerformManualClick() {
	ge.state.Bits += ge.manualClickValue()
	ge.state.TotalBitsEarned += ge.manualClickValue()
	ge.state.ClickerFlash = true
	audio.PlayClick()
}

func (ge *GameEngine) Reboot() {
	gain := math.Log10(ge.state.TotalBitsEarned/1_000_000) * 0.1
	if gain < 0 {
		gain = 0
	}
	ge.state.GHzMultiplier += gain
	ge.state.RebootCount++

	ge.state.Bits = 0
	ge.state.TotalBitsEarned = 0
	ge.state.Entropy = 0
	ge.state.Corruption = 0
	ge.state.Hardware = make(map[string]int)
	ge.state.Upgrades = make(map[string]bool)
	ge.state.RebootPending = false
	ge.state.Sanitize() // Re-adds basic rack shelf
	ge.state.LogMessage(fmt.Sprintf("[SYSTEM] PURGE_COMPLETE: +%.3fX GHz", gain))
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

	if ge.state.Bits < cost {
		return
	}

	// Constraint Check: Space
	if target.SpaceImpact > 0 && ge.state.SpaceUsage+target.SpaceImpact > ge.state.SpaceCapacity {
		ge.state.LogMessage("[ERROR] INSUFFICIENT_RACK_SPACE")
		return
	}

	// Constraint Check: Power (Allow buying PSUs always, but check consumption for others)
	if target.WattsImpact > 0 && ge.state.PowerUsage+target.WattsImpact > ge.state.PowerCapacity {
		ge.state.LogMessage("[ERROR] POWER_OVERLOAD_PREVENTED")
		return
	}

	ge.state.Bits -= cost
	ge.state.Hardware[id]++
	audio.PlayClick()
	ge.state.LogMessage(fmt.Sprintf("[INSTALL] %s", target.Name))
}

func (ge *GameEngine) gameTick(dt float64) {
	// 1. Calculate Multipliers
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

	// 2. Reset and Recompute Infrastructure Totals
	ge.state.PowerUsage = 0
	ge.state.PowerCapacity = 0
	ge.state.SpaceUsage = 0
	ge.state.SpaceCapacity = 0
	
	totalHeatGen := 0.0
	totalCooling := 10.0 // Base ambient cooling
	bps := 0.0
	entropyDelta := 0.0

	for _, def := range model.AllHardware {
		count := float64(ge.state.Hardware[def.ID])
		if count > 0 {
			// Bits & Entropy
			bps += count * def.BaseBPS * prodMult
			entropyDelta += count * def.EntropyWeight * entropyMult

			// Power
			if def.WattsImpact > 0 {
				ge.state.PowerUsage += count * def.WattsImpact
			} else {
				ge.state.PowerCapacity += count * math.Abs(def.WattsImpact)
			}

			// Thermal
			if def.ThermalImpact > 0 {
				totalHeatGen += count * def.ThermalImpact
			} else {
				totalCooling += count * math.Abs(def.ThermalImpact)
			}

			// Space
			if def.SpaceImpact > 0 {
				ge.state.SpaceUsage += count * def.SpaceImpact
			} else {
				ge.state.SpaceCapacity += count * math.Abs(def.SpaceImpact)
			}
		}
	}

	// Calculate Heat Level (0-100%)
	if totalHeatGen > 0 {
		ge.state.HeatLevel = (totalHeatGen / totalCooling) * 50.0 // 50 is nominal
	} else {
		ge.state.HeatLevel = 0
	}

	// 3. Apply Multipliers & Penalties
	bps *= ge.state.GHzMultiplier
	
	// Thermal Penalty
	if ge.state.HeatLevel > 80 {
		thermalPenalty := (ge.state.HeatLevel - 80) / 40.0 // Up to 0.5 at 100%
		bps *= (1.0 - math.Min(thermalPenalty, 0.5))
		ge.state.Corruption += (ge.state.HeatLevel - 80) * 0.005 * dt
	}

	// Power Penalty (if somehow exceeded)
	if ge.state.PowerUsage > ge.state.PowerCapacity {
		ge.state.Entropy += 2.0 * dt
	}

	// Corruption Penalty
	corruptPenalty := math.Min(ge.state.Corruption/200.0, 0.5)
	bps *= (1.0 - corruptPenalty)

	// 4. Update Economy
	earned := bps * dt
	ge.state.Bits += earned
	ge.state.TotalBitsEarned += earned

	// 5. Entropy & Corruption
	ge.state.Entropy += (entropyDelta + (ge.state.PowerUsage / 1000.0)) * dt
	if ge.state.Entropy < 0 { ge.state.Entropy = 0 }
	if ge.state.Entropy > 100 { ge.state.Entropy = 100 }

	if ge.state.Entropy > 50 {
		corruptDelta := (ge.state.Entropy - 50) * 0.002
		ge.state.Corruption += corruptDelta * dt
	}

	if ge.state.Corruption > 75 {
		decayRate := ge.state.Bits * (ge.state.Corruption - 75) * 0.0001
		ge.state.Bits -= decayRate * dt
	}

	// Clamp
	if ge.state.Corruption > 100 { ge.state.Corruption = 100 }
	if ge.state.Bits < 0 { ge.state.Bits = 0 }

	// 6. Random Events & Audio
	if ge.state.Corruption > 90 || ge.state.HeatLevel > 95 {
		ge.alarmTimer += dt
		if ge.alarmTimer >= 3.0 {
			audio.PlayAlarm()
			ge.alarmTimer = 0
		}
	} else {
		ge.alarmTimer = 0
	}

	if !ge.state.PacketActive {
		if math.Mod(ge.state.Bits*123.45, 100) < 0.5 {
			if ge.state.TotalBitsEarned > 1024 {
				ge.state.PacketActive = true
				ge.state.PacketTimer = 8.0
			}
		}
	} else {
		ge.state.PacketTimer -= dt
		if ge.state.PacketTimer <= 0 {
			ge.state.PacketActive = false
		}
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
