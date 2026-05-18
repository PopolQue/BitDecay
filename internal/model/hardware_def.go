package model

import (
	"math"
)

type HardwareDef struct {
	ID            string
	Name          string
	Tier          int
	BaseCost      float64
	CostScaling   float64 // per-purchase multiplier
	BaseBPS       float64 // bits/second contribution
	EntropyWeight float64 // base entropy contribution

	// Infrastructure impacts
	WattsImpact   float64 // + consumes power, - provides power (PSUs)
	ThermalImpact float64 // + generates heat, - cools system (Fans)
	SpaceImpact   float64 // + consumes units, - provides units (Racks)

	Description string
}

var AllHardware = []HardwareDef{
	// Infrastructure: Tier 1
	{ID: "rack_shelf", Name: "Basic Rack Shelf", Tier: 1, BaseCost: 50, CostScaling: 1.10, SpaceImpact: -4.0, Description: "Provides 4U of rack space."},
	{ID: "psu_450", Name: "450W PSU", Tier: 1, BaseCost: 80, CostScaling: 1.15, WattsImpact: -450, ThermalImpact: 10, Description: "Provides 450W of power."},
	{ID: "fan_80", Name: "80mm Case Fan", Tier: 1, BaseCost: 30, CostScaling: 1.12, WattsImpact: 5, ThermalImpact: -20, Description: "Basic cooling. Consumes 5W."},

	// Compute: Tier 1
	{ID: "cpu_i3", Name: "Dual-Core CPU", Tier: 1, BaseCost: 150, CostScaling: 1.20, BaseBPS: 5.0, WattsImpact: 65, ThermalImpact: 40, SpaceImpact: 1, Description: "Generates bits. High heat."},
	{ID: "ram_8", Name: "8GB DDR4 RAM", Tier: 1, BaseCost: 60, CostScaling: 1.15, BaseBPS: 0.5, WattsImpact: 10, ThermalImpact: 5, SpaceImpact: 0, Description: "Low bit gen, but essential for stability."},

	// Infrastructure: Tier 2
	{ID: "rack_cabinet", Name: "42U Server Rack", Tier: 2, BaseCost: 1000, CostScaling: 1.25, SpaceImpact: -42, Description: "A full server cabinet. Huge capacity."},
	{ID: "psu_1200", Name: "1200W Platinum PSU", Tier: 2, BaseCost: 500, CostScaling: 1.20, WattsImpact: -1200, ThermalImpact: 25, Description: "High-efficiency power delivery."},
	{ID: "liquid_cooler", Name: "AIO Liquid Cooler", Tier: 2, BaseCost: 250, CostScaling: 1.15, WattsImpact: 25, ThermalImpact: -150, SpaceImpact: 1, Description: "Advanced thermal management."},

	// Compute: Tier 2
	{ID: "cpu_threadripper", Name: "64-Core Threadripper", Tier: 2, BaseCost: 5000, CostScaling: 1.30, BaseBPS: 500.0, WattsImpact: 280, ThermalImpact: 180, SpaceImpact: 1, Description: "Extreme processing power."},
	{ID: "gpu_rtx", Name: "RTX Compute Card", Tier: 2, BaseCost: 3500, CostScaling: 1.30, BaseBPS: 800.0, WattsImpact: 350, ThermalImpact: 220, SpaceImpact: 2, Description: "Massive bit throughput. Huge power draw."},
}

func CurrentCost(def HardwareDef, owned int) float64 {
	return def.BaseCost * math.Pow(def.CostScaling, float64(owned))
}
