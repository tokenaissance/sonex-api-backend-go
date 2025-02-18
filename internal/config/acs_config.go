package config

type ACSConfig struct {
	Tokens    []string `json:","`
	Points    int64    `json:","`
	SubTitle  string   `json:","`
	Host      string   `json:","`
	ApiSecret string   `json:","`
}
