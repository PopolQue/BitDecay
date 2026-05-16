package constants

// Quarters and Rounds
const (
	TotalRounds    = 12
	RoundsPerQuarter = 3
)

// Asset Prices and Properties
const (
	PriceTent           = 100
	PriceGlampingTent   = 250
	PriceCaravan        = 500
	PriceBungalow       = 1000
	PriceLuxBungalow    = 2000
	PriceGenerator      = 300
	PriceWaterTank      = 300
	PriceSportField     = 800
	PriceFirePit        = 200
	PriceSauna          = 600
	PriceStage          = 1200
	PriceLandExpansion  = 1500
)

const (
	AreaTent           = 10
	AreaGlampingTent   = 15
	AreaCaravan        = 20
	AreaBungalow       = 40
	AreaLuxBungalow    = 50
	AreaGenerator      = 15
	AreaWaterTank      = 15
	AreaSportField     = 100
	AreaFirePit        = 10
	AreaSauna          = 30
	AreaStage          = 60
)

const (
	CapacityTent           = 2
	CapacityGlampingTent   = 4
	CapacityCaravan        = 4
	CapacityBungalow       = 6
	CapacityLuxBungalow    = 8
)

// Action Points
const (
	MaxActionPoints = 3
)

// Asset Upgrade Costs
const (
	PriceUpgradeToGlamping = 150
	PriceUpgradeToLuxus    = 1000
)

// Board Constants
const (
	MapRadius = 6
	HexSize   = 40
)

// Resource Generation
const (
	GeneratorPowerOutput = 50
	WaterTankOutput      = 50
	InitialMoney         = 2000
	InitialArea          = 200
)

// Events
type EventEffect string

const (
	EffectNPCHalf         EventEffect = "NPCHalf"
	EffectNPCDouble       EventEffect = "NPCDouble"
	EffectNoTents         EventEffect = "NoTents"
	EffectPreferTents     EventEffect = "PreferTents"
	EffectExtraThirst     EventEffect = "ExtraThirst"
	EffectNoWaterGen      EventEffect = "NoWaterGen"
	EffectNoPowerGen      EventEffect = "NoPowerGen"
)
