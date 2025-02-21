package config

type ACSConfig struct {
	Pools              [][]string `json:","`
	Tokens             []string   `json:","`
	Points             int64      `json:","`
	SubTitle           string     `json:","`
	Graphql            string     `json:","`
	Host               string     `json:","`
	DefiId             int64      `json:","`
	ApiSecret          string     `json:","`
	ASTR               string     `json:","`
	LiquidityUSDLimits float64    `json:","`
	ACSDailyLimits     int64      `json:","`
}
