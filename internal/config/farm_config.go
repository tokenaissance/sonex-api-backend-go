package config

// Represents Farm configuration of a blockchain
type FarmConfig struct {
	ChainId int64  `json:","`
	Staker  string `json:","`
}
