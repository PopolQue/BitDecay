package models

import "github.com/google/uuid"

type AssetType string

const (
	AssetTent           AssetType = "Tent"
	AssetGlampingTent   AssetType = "GlampingTent"
	AssetCaravan        AssetType = "Caravan"
	AssetBungalow       AssetType = "Bungalow"
	AssetLuxBungalow    AssetType = "LuxBungalow"
	AssetGenerator      AssetType = "Generator"
	AssetWaterTank      AssetType = "WaterTank"
	AssetSportField     AssetType = "SportField"
	AssetFirePit        AssetType = "FirePit"
	AssetSauna          AssetType = "Sauna"
	AssetStage          AssetType = "Stage"
)

type Asset struct {
	ID          uuid.UUID
	Type        AssetType
	Capacity    int
	Occupants   []uuid.UUID
	IsUpgraded  bool
	Q, R        int // Axial coordinates
}
