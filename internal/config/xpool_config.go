package config

// Represents Xpool configuration of a blockchain.
type XpoolConfig struct {
	ChainId   int64  `json:","`
	Factory   string `json:","`
	BlockTime int64  `json:","`
}
