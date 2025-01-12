package config

type PreDepositConfig struct {
	ChainId            int64   `json:","`
	BridgeStaker       string  `json:","`
	RewardClaimer      string  `json:","`
	SignerPrivateKey   string  `json:","`
	PointsRate         float64 `json:","`
	SubTitle           string  `json:","`
	StakeTokenDecimals int64   `json:","`
	RewardToken        string  `json:","`
	RewardRate         float64 `json:","`
}
