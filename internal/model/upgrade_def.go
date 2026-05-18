package model

type UpgradeType int

const (
	UpgradeClickBoost UpgradeType = iota
	UpgradeProductionBoost
	UpgradeEntropyReduction
	UpgradePowerEfficiency   // Reduces power draw of compute
	UpgradeCoolingEfficiency // Increases fan/cooling effectiveness
)

type UpgradeDef struct {
	ID          string
	Name        string
	Cost        float64
	Type        UpgradeType
	Multiplier  float64
	Description string
	PrereqID    string // Optional prerequisite upgrade
}

var AllUpgrades = []UpgradeDef{
	{
		ID:          "overclock_1",
		Name:        "Overclock Tier 1",
		Cost:        500,
		Type:        UpgradeClickBoost,
		Multiplier:  2.0,
		Description: "Doubles manual bit generation.",
	},
	{
		ID:          "logic_opt_1",
		Name:        "Logic Optimization",
		Cost:        1000,
		Type:        UpgradeProductionBoost,
		Multiplier:  1.1,
		Description: "Increases all hardware production by 10%.",
	},
	{
		ID:          "shielding_1",
		Name:        "EM Shielding",
		Cost:        2500,
		Type:        UpgradeEntropyReduction,
		Multiplier:  0.9,
		Description: "Reduces entropy generation by 10%.",
	},
	// Infrastructure Upgrades
	{
		ID:          "volt_mod_1",
		Name:        "Voltage Undervolting",
		Cost:        5000,
		Type:        UpgradePowerEfficiency,
		Multiplier:  0.8,
		Description: "Reduces CPU/GPU power draw by 20%.",
	},
	{
		ID:          "thermal_paste",
		Name:        "Liquid Metal Paste",
		Cost:        8000,
		Type:        UpgradeCoolingEfficiency,
		Multiplier:  1.25,
		Description: "Increases cooling effectiveness by 25%.",
	},
	{
		ID:          "copper_heatsinks",
		Name:        "Copper Fin Arrays",
		Cost:        15000,
		Type:        UpgradeCoolingEfficiency,
		Multiplier:  1.5,
		Description: "Increases cooling effectiveness by 50%.",
	},
	{
		ID:          "datacenter_link",
		Name:        "Backplane Sync",
		Cost:        50000,
		Type:        UpgradeProductionBoost,
		Multiplier:  1.5,
		Description: "Increases all bit production by 50%.",
	},
}
