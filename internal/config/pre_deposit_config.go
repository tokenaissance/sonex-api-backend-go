package config

type PreDepositConfig struct {
	ChainId          int64   `json:","`
	WASTR            string  `json:","`
	VASTR            string  `json:","`
	BridgeStaker     string  `json:","`
	RewardClaimer    string  `json:","`
	SignerPrivateKey string  `json:","`
	PointsRate       float64 `json:","`
	SubTitle         string  `json:","`
	RewardToken      string  `json:","`
	RewardRate       float64 `json:","`
}
