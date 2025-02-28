package config

type PufferConfig struct {
	Pools              [][]string `json:","`
	Points             int64      `json:","`
	SubTitle           string     `json:","`
	LiquidityUSDLimits float64    `json:","`
	DailyLimits        int64      `json:","`
}
