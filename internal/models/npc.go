package models

import "github.com/google/uuid"

type NPCType string

const (
	NPCHippie NPCType = "Hippie"
	NPCFamily NPCType = "Family"
	NPCSnob   NPCType = "Snob"
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
