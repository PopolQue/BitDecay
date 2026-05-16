package models

import "github.com/google/uuid"

type NPCType string

const (
	NPCHippie NPCType = "Hippie"
	NPCFamily NPCType = "Family"
	NPCSnob   NPCType = "Snob"
)

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

type NPC struct {
	ID               uuid.UUID
	Name             string
	Type             NPCType
	StayDuration     int // in rounds
	GroupSize        int
	IncomePerNight   int
	PowerNeed        int
	WaterNeed        int
	SpecialNeeds     []AssetType
	AssignedAssetID  uuid.UUID
}

type Asset struct {
	ID          uuid.UUID
	Type        AssetType
	Capacity    int
	Occupants   []uuid.UUID
}

type Resources struct {
	Money int
	Area  int // Total available area
	UsedArea int
	Power int
	Water int
}

type GameState struct {
	Round        int
	Quarter      int
	Resources    Resources
	Assets       []Asset
	GuestPool    []NPC
	ActiveGuests []NPC
}
