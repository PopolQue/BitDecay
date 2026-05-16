package game

import (
	"math/rand"
	"time"

	"github.com/popolque/firstbitengi/internal/constants"
	"github.com/popolque/firstbitengi/internal/models"
	"github.com/google/uuid"
)

func (e *Engine) refreshNPCPool() {
	e.State.GuestPool = []models.NPC{}
	numNPCs := 4 
	if e.State.Quarter == 1 || e.State.Quarter == 4 {
		numNPCs -= 1
	} else {
		numNPCs += 1
	}
	if e.State.ActiveEvent != nil {
		if e.State.ActiveEvent.Effect == string(constants.EffectNPCHalf) {
			numNPCs /= 2
		} else if e.State.ActiveEvent.Effect == string(constants.EffectNPCDouble) {
			numNPCs *= 2
		}
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
		npc.Name = "Hippie"
		npc.IncomePerNight = 20
		npc.PowerNeed = 5
		npc.WaterNeed = 5
	case models.NPCFamily:
		npc.Name = "Family"
		npc.IncomePerNight = 50
		npc.PowerNeed = 15
		npc.WaterNeed = 20
		npc.SpecialNeeds = []models.AssetType{models.AssetSportField}
	case models.NPCSnob:
		npc.Name = "Snob"
		npc.IncomePerNight = 100
		npc.PowerNeed = 30
		npc.WaterNeed = 15
		npc.SpecialNeeds = []models.AssetType{models.AssetSauna}
	}
	return npc
}
