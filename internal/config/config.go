package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	Debug bool `json:",env=DEBUG,default=false"`
	rest.RestConf
	Domains    []string      `json:",optional,default=[]"`
	PostgreSQL PgConfig      `json:","`
	Chains     []ChainConfig `json:","`
	Sonexes    []SonexConfig `json:",optional"`
	Xpools     []XpoolConfig `json:",optional"`
	Farms      []FarmConfig  `json:",optional"`
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

func (config *Config) GetXpoolConfig(chainId int64) *XpoolConfig {
	for _, c := range config.Xpools {
		if c.ChainId == chainId {
			return &c
		}
	}
	return nil
}

func (config *Config) GetFarmConfig(chainId int64) *FarmConfig {
	for _, c := range config.Farms {
		if c.ChainId == chainId {
			return &c
		}
	}
	return nil
}
