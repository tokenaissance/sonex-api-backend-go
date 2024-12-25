package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	Debug bool `json:",env=DEBUG,default=false"`
	rest.RestConf
	Demains    []string      `json:",optional,default=[]"`
	PostgreSQL PgConfig      `json:","`
	Chains     []ChainConfig `json:","`
	Sonexes    []SonexConfig `json:",optional"`
	Xpools     []XpoolConfig `json:",optional"`
	Farm       []FarmConfig  `json:",optional"`
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

func (config *Config) GetSonexConfig(chainId int64) *SonexConfig {
	for _, c := range config.Sonexes {
		if c.ChainId == chainId {
			return &c
		}
	}
	return nil
}
