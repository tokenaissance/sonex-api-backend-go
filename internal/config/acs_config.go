package config

type ACSConfig struct {
	Pools              [][]string `json:","`
	Tokens             []string   `json:","`
	Points             int64      `json:","`
	SubTitle           string     `json:","`
	Host               string     `json:","`
	ApiSecret          string     `json:","`
	ASTR               string     `json:","`
	LiquidityUSDLimits float64    `json:","`
}
