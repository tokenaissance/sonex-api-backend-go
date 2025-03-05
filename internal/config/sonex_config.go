package config

// Represents Sonex configuration of a blockchain.
type SonexConfig struct {
	ChainId         int64  `json:","`
	Factory         string `json:","`
	PositionManager string `json:","`
	Quoter          string `json:",optional"`
	SwapRouter      string `json:",optional"`
	XToken          string `json:","`
	XTokenLogo      string `json:",optional"`
}
