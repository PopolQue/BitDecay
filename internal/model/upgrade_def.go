package model

type UpgradeType int

const (
	UpgradeClickBoost UpgradeType = iota
	UpgradeProductionBoost
	UpgradeEntropyReduction
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
}
