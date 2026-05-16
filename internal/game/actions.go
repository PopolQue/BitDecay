package game

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/popolque/firstbitengi/internal/constants"
	"github.com/popolque/firstbitengi/internal/models"
)

func (e *Engine) BuyTile(q, r int) bool {
	if e.State.ActionPoints <= 0 {
		return false
	}
	tileKey := fmt.Sprintf("%d,%d", q, r)
	owner, exists := e.State.Board[tileKey]
	if !exists || owner != -1 {
		return false
	}
	if !e.checkAdjacency(q, r) {
		return false
	}
	if e.State.Resources.Money >= 100 {
		e.State.Resources.Money -= 100
		e.State.Board[tileKey] = 0 // Player 0
		e.State.ActionPoints--
		return true
	}
	return false
}

func (e *Engine) checkAdjacency(q, r int) bool {
	neighbors := [][]int{
		{q + 1, r}, {q - 1, r}, {q, r + 1}, {q, r - 1}, {q + 1, r - 1}, {q - 1, r + 1},
	}
	for _, n := range neighbors {
		key := fmt.Sprintf("%d,%d", n[0], n[1])
		if owner, exists := e.State.Board[key]; exists && owner == 0 {
			return true
		}
	}
	return false
}

func (e *Engine) BuyAsset(aType models.AssetType, q, r int) bool {
	if e.State.ActionPoints <= 0 {
		return false
	}
	// Check if tile is owned and empty
	tileKey := fmt.Sprintf("%d,%d", q, r)
	if owner, exists := e.State.Board[tileKey]; !exists || owner != 0 {
		return false
	}
	for _, a := range e.State.Assets {
		if a.Q == q && a.R == r {
			return false // Already has an asset
		}
	}

	price := 0
	capacity := 0
	switch aType {
	case models.AssetTent:
		price, capacity = constants.PriceTent, constants.CapacityTent
	case models.AssetCaravan:
		price, capacity = constants.PriceCaravan, constants.CapacityCaravan
	case models.AssetBungalow:
		price, capacity = constants.PriceBungalow, constants.CapacityBungalow
	case models.AssetGenerator:
		price = constants.PriceGenerator
	case models.AssetWaterTank:
		price = constants.PriceWaterTank
	case models.AssetSportField:
		price = constants.PriceSportField
	case models.AssetSauna:
		price = constants.PriceSauna
	}

	if e.State.Resources.Money >= price {
		e.State.Resources.Money -= price
		e.State.Assets = append(e.State.Assets, models.Asset{
			ID:       uuid.New(),
			Type:     aType,
			Capacity: capacity,
			Q:        q,
			R:        r,
		})
		e.State.ActionPoints--
		return true
	}
	return false
}
