package model

import (
	"math"
)

type HardwareDef struct {
	ID            string
	Name          string
	Tier          int
	BaseCost      float64
	CostScaling   float64 // per-purchase multiplier (default 1.15)
	BaseBPS       float64 // bits/second contribution
	EntropyWeight float64 // entropy/second (negative = reduction)
	Description   string
}

var AllHardware = []HardwareDef{
	{ID: "logic_gate", Name: "Logic Gate", Tier: 1, BaseCost: 10, CostScaling: 1.15, BaseBPS: 0.1, EntropyWeight: 0.01},
	{ID: "alu", Name: "ALU", Tier: 1, BaseCost: 100, CostScaling: 1.15, BaseBPS: 0.8, EntropyWeight: 0.05},
	{ID: "ecc_memory", Name: "ECC Memory", Tier: 1, BaseCost: 80, CostScaling: 1.12, BaseBPS: 0.0, EntropyWeight: -0.08},
	{ID: "heatsink", Name: "Heatsink & Fan", Tier: 1, BaseCost: 120, CostScaling: 1.12, BaseBPS: 0.0, EntropyWeight: -0.12},
	{ID: "quantum_core", Name: "Quantum Core", Tier: 2, BaseCost: 5000, CostScaling: 1.18, BaseBPS: 15.0, EntropyWeight: 0.80},
	{ID: "neural_link", Name: "Neural Link", Tier: 2, BaseCost: 12000, CostScaling: 1.18, BaseBPS: 35.0, EntropyWeight: 1.50},
	{ID: "ai_kernel", Name: "AI Kernel", Tier: 3, BaseCost: 500000, CostScaling: 1.20, BaseBPS: 250.0, EntropyWeight: 3.00},
	{ID: "temp_buffer", Name: "Temporal Buffer", Tier: 3, BaseCost: 1200000, CostScaling: 1.20, BaseBPS: 120.0, EntropyWeight: -2.50},
}

func CurrentCost(def HardwareDef, owned int) float64 {
	return def.BaseCost * math.Pow(def.CostScaling, float64(owned))
}
