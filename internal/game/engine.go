package game

import (
	"math/rand"
	"time"

	"firstEbitengi/internal/constants"
	"firstEbitengi/internal/models"
	"github.com/google/uuid"
)

type Engine struct {
	State models.GameState
}

func NewEngine() *Engine {
	return &Engine{
		State: models.GameState{
			Round:   1,
			Quarter: 1,
			Resources: models.Resources{
				Money: constants.InitialMoney,
				Area:  constants.InitialArea,
			},
		},
	}
}

func (e *Engine) EndTurn() {
	// 1. Calculate Income and Resource Consumption
	e.calculateResources()

	// 2. Update Active Guests (decrement stay duration)
	e.updateGuests()

	// 3. Advance Round and Quarter
	e.State.Round++
	if e.State.Round > constants.TotalRounds {
		// End Game Logic could go here
		return
	}

	e.State.Quarter = ((e.State.Round - 1) / constants.RoundsPerQuarter) + 1

	// 4. Refresh NPC Pool
	e.refreshNPCPool()
}

func (e *Engine) calculateResources() {
	income := 0
	powerGen := 0
	waterGen := 0
	powerCons := 0
	waterCons := 0

	for _, asset := range e.State.Assets {
		if asset.Type == models.AssetGenerator {
			powerGen += constants.GeneratorPowerOutput
		}
		if asset.Type == models.AssetWaterTank {
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
			// Guest leaves, free up asset
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

func (e *Engine) refreshNPCPool() {
	e.State.GuestPool = []models.NPC{}
	
	// Base number of NPCs depends on season (GDD logic)
	numNPCs := 4 // Placeholder base
	if e.State.Quarter == 1 || e.State.Quarter == 4 {
		numNPCs -= 1
	} else {
		numNPCs += 1
	}

	for i := 0; i < numNPCs; i++ {
		e.State.GuestPool = append(e.State.GuestPool, generateRandomNPC())
	}
}

func generateRandomNPC() models.NPC {
	rand.Seed(time.Now().UnixNano())
	types := []models.NPCType{models.NPCHippie, models.NPCFamily, models.NPCSnob}
	t := types[rand.Intn(len(types))]

	npc := models.NPC{
		ID:           uuid.New(),
		Type:         t,
		StayDuration: rand.Intn(3) + 1,
		GroupSize:    rand.Intn(4) + 1,
	}

	switch t {
	case models.NPCHippie:
		npc.Name = "Peaceful Paul"
		npc.IncomePerNight = 20
		npc.PowerNeed = 5
		npc.WaterNeed = 5
	case models.NPCFamily:
		npc.Name = "The Millers"
		npc.IncomePerNight = 50
		npc.PowerNeed = 15
		npc.WaterNeed = 20
	case models.NPCSnob:
		npc.Name = "Lord Fancy"
		npc.IncomePerNight = 100
		npc.PowerNeed = 30
		npc.WaterNeed = 15
	}

	return npc
}

func (e *Engine) BuyAsset(aType models.AssetType) bool {
	price := 0
	area := 0
	capacity := 0

	switch aType {
	case models.AssetTent:
		price, area, capacity = constants.PriceTent, constants.AreaTent, constants.CapacityTent
	case models.AssetCaravan:
		price, area, capacity = constants.PriceCaravan, constants.AreaCaravan, constants.CapacityCaravan
	case models.AssetBungalow:
		price, area, capacity = constants.PriceBungalow, constants.AreaBungalow, constants.CapacityBungalow
	case models.AssetGenerator:
		price, area = constants.PriceGenerator, constants.AreaGenerator
	case models.AssetWaterTank:
		price, area = constants.PriceWaterTank, constants.AreaWaterTank
	}

	if e.State.Resources.Money >= price && e.State.Resources.Area-e.State.Resources.UsedArea >= area {
		e.State.Resources.Money -= price
		e.State.Resources.UsedArea += area
		e.State.Assets = append(e.State.Assets, models.Asset{
			ID:       uuid.New(),
			Type:     aType,
			Capacity: capacity,
		})
		return true
	}
	return false
}

func (e *Engine) AcceptGuest(npcID uuid.UUID, assetID uuid.UUID) bool {
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
				// Assignment logic (simple for now: one group per asset)
				npc.AssignedAssetID = assetID
				e.State.Assets[i].Occupants = append(e.State.Assets[i].Occupants, npc.ID)
				e.State.ActiveGuests = append(e.State.ActiveGuests, npc)
				
				// Remove from pool
				e.State.GuestPool = append(e.State.GuestPool[:npcIdx], e.State.GuestPool[npcIdx+1:]...)
				return true
			}
		}
	}
	return false
}
