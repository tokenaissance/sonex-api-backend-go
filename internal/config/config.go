package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	Debug bool `json:",env=DEBUG,default=false"`
	rest.RestConf
	PostgreSQL PgConfig      `json:","`
	Chains     []ChainConfig `json:","`
	Sonexes    []SonexConfig `json:",optional"`
	Xpools     []XpoolConfig `json:",optional"`
	Redis      RedisConfig   `json:","`
}

func (config *Config) GetChainConfig(chainId int64) *ChainConfig {
	for _, c := range config.Chains {
		if c.ChainId == chainId {
			return &c
		}
	}
	return nil
}
