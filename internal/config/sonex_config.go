package config

type SonexConfig struct {
	ChainId         int64  `json:","`
	Factory         string `json:","`
	PositionManager string `json:","`
	Quoter          string `json:",optional"`
	SwapRouter      string `json:",optional"`
}
