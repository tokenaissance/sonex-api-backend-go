package config

type PreDepositConfig struct {
	ChainId          int64   `json:","`
	WASTR            string  `json:","`
	VASTR            string  `json:","`
	ASTR             string  `json:","`
	BridgeStaker     string  `json:","`
	RewardClaimer    string  `json:","`
	SignerPrivateKey string  `json:","`
	PointsRate       float64 `json:","`
	PointsUserCap    int64   `json:","`
	PointsTotalCap   int64   `json:","`
	SubTitle         string  `json:","`
	RewardToken      string  `json:","`
	RewardDaily      int64   `json:","`
	RewardPrice      int64   `json:","`
}
