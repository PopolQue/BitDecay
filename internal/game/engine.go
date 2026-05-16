package game

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/popolque/firstbitengi/internal/constants"
	"github.com/popolque/firstbitengi/internal/models"
	"github.com/google/uuid"
)

type Engine struct {
	State models.GameState
}

func NewEngine() *Engine {
	engine := &Engine{
		State: models.GameState{
			Round:        1,
			Quarter:      1,
			ActionPoints: constants.MaxActionPoints,
			Resources: models.Resources{
				Money: constants.InitialMoney,
				Area:  constants.InitialArea,
			},
			Board: make(map[string]int),
		},
	}

	radius := constants.MapRadius
	for q := -radius; q <= radius; q++ {
		for r := -radius; r <= radius; r++ {
			if abs(q)+abs(r)+abs(-q-r) <= radius*2 {
				engine.State.Board[fmt.Sprintf("%d,%d", q, r)] = -1
			}
		}
	}

	engine.State.Board["0,0"] = 0 
	return engine
}

func (e *Engine) EndTurn() {
	e.calculateResources()
	e.updateGuests()
	e.State.ActionPoints = constants.MaxActionPoints
	e.State.Round++
	if e.State.Round > constants.TotalRounds {
		return
	}
	newQuarter := ((e.State.Round - 1) / constants.RoundsPerQuarter) + 1
	if newQuarter != e.State.Quarter {
		e.State.Quarter = newQuarter
		e.triggerQuarterlyEvent()
	}
	e.refreshNPCPool()
}

func (e *Engine) triggerQuarterlyEvent() {
	events := []models.Event{
		{Name: "Dauerregen", Description: "NPCs pro Zug halbieren sich", Effect: string(constants.EffectNPCHalf)},
		{Name: "Tourismusboom", Description: "NPCs pro Zug verdoppeln sich", Effect: string(constants.EffectNPCDouble)},
		{Name: "Sturmwarnung", Description: "NPCs wollen nicht in Zelten schlafen", Effect: string(constants.EffectNoTents)},
		{Name: "Blackout", Description: "Generatoren liefern keinen Strom", Effect: string(constants.EffectNoPowerGen)},
		{Name: "Dürre", Description: "Wassertanks liefern kein Wasser", Effect: string(constants.EffectNoWaterGen)},
	}
	rand.Seed(time.Now().UnixNano())
	selected := events[rand.Intn(len(events))]
	e.State.ActiveEvent = &selected
}

func (e *Engine) calculateResources() {
	income := 0
	powerGen := 0
	waterGen := 0
	powerCons := 0
	waterCons := 0
	noPower := e.State.ActiveEvent != nil && e.State.ActiveEvent.Effect == string(constants.EffectNoPowerGen)
	noWater := e.State.ActiveEvent != nil && e.State.ActiveEvent.Effect == string(constants.EffectNoWaterGen)

	for _, asset := range e.State.Assets {
		if asset.Type == models.AssetGenerator && !noPower {
			powerGen += constants.GeneratorPowerOutput
		}
		if asset.Type == models.AssetWaterTank && !noWater {
			waterGen += constants.WaterTankOutput
		}
	}
	for _, guest := range e.State.ActiveGuests {
		income += guest.IncomePerNight
		powerCons += guest.PowerNeed
		waterCons += guest.WaterNeed
	}
	e.State.Resources.Money += income
	e.State.Resources.Power = powerGen - powerCons
	e.State.Resources.Water = waterGen - waterCons
}

func (e *Engine) updateGuests() {
	var remainingGuests []models.NPC
	for _, guest := range e.State.ActiveGuests {
		guest.StayDuration--
		if guest.StayDuration > 0 {
			remainingGuests = append(remainingGuests, guest)
		} else {
			for i, asset := range e.State.Assets {
				if asset.ID == guest.AssignedAssetID {
					for j, occupantID := range asset.Occupants {
						if occupantID == guest.ID {
							e.State.Assets[i].Occupants = append(asset.Occupants[:j], asset.Occupants[j+1:]...)
							break
						}
					}
				}
			}
		}
	}
	e.State.ActiveGuests = remainingGuests
}

func (e *Engine) UpgradeAsset(assetID uuid.UUID) bool {
	if e.State.ActionPoints <= 0 {
		return false
	}
	for i, asset := range e.State.Assets {
		if asset.ID == assetID && !asset.IsUpgraded {
			cost := 0
			newCapacity := 0
			newType := asset.Type
			if asset.Type == models.AssetTent {
				cost = constants.PriceUpgradeToGlamping
				newCapacity = constants.CapacityGlampingTent
				newType = models.AssetGlampingTent
			} else if asset.Type == models.AssetBungalow {
				cost = constants.PriceUpgradeToLuxus
				newCapacity = constants.CapacityLuxBungalow
				newType = models.AssetLuxBungalow
			}
			if cost > 0 && e.State.Resources.Money >= cost {
				e.State.Resources.Money -= cost
				e.State.Assets[i].IsUpgraded = true
				e.State.Assets[i].Type = newType
				e.State.Assets[i].Capacity = newCapacity
				e.State.ActionPoints--
				return true
			}
		}
	}
	return false
}

func (e *Engine) AcceptGuest(npcID uuid.UUID, assetID uuid.UUID) bool {
	if e.State.ActionPoints <= 0 {
		return false
	}
	var npc models.NPC
	npcIdx := -1
	for i, n := range e.State.GuestPool {
		if n.ID == npcID {
			npc = n
			npcIdx = i
			break
		}
	}
	if npcIdx == -1 {
		return false
	}
	for i, asset := range e.State.Assets {
		if asset.ID == assetID {
			if len(asset.Occupants) == 0 && asset.Capacity >= npc.GroupSize {
				if !e.isPlacementAllowed(npc, asset) {
					return false
				}
				for _, need := range npc.SpecialNeeds {
					found := false
					for _, a := range e.State.Assets {
						if a.Type == need {
							found = true
							break
						}
					}
					if !found {
						return false
					}
				}
				npc.AssignedAssetID = assetID
				e.State.Assets[i].Occupants = append(e.State.Assets[i].Occupants, npc.ID)
				e.State.ActiveGuests = append(e.State.ActiveGuests, npc)
				e.State.GuestPool = append(e.State.GuestPool[:npcIdx], e.State.GuestPool[npcIdx+1:]...)
				e.State.ActionPoints--
				return true
			}
		}
	}
	return false
}

func (e *Engine) isPlacementAllowed(npc models.NPC, asset models.Asset) bool {
	if e.State.ActiveEvent != nil && e.State.ActiveEvent.Effect == string(constants.EffectNoTents) {
		if asset.Type == models.AssetTent || asset.Type == models.AssetGlampingTent {
			return false
		}
	}
	switch npc.Type {
	case models.NPCHippie:
		return asset.Type == models.AssetTent || asset.Type == models.AssetCaravan || asset.Type == models.AssetBungalow
	case models.NPCFamily:
		return asset.Type == models.AssetGlampingTent || asset.Type == models.AssetCaravan || asset.Type == models.AssetBungalow
	case models.NPCSnob:
		return asset.Type == models.AssetGlampingTent || asset.Type == models.AssetBungalow || asset.Type == models.AssetLuxBungalow
	}
	return false
}
