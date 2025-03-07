package config

// Represents Farm configuration of a blockchain
type FarmConfig struct {
	ChainId              int64             `json:","`
	Staker               string            `json:","`
	PointsHoldTimeLimits int64             `json:","`
	PointsUSDLimits      float64           `json:","`
	Points               int64             `json:","`
	SubTitle             string            `json:","`
	XRewardToken         string            `json:","`
	Boosts               []FarmBoostConfig `json:","`
}

type FarmBoostConfig struct {
	Tokens               []string `json:","`
	PointsHoldTimeLimits int64    `json:","`
	PointsUSDLimits      float64  `json:","`
	Points               int64    `json:","`
}
