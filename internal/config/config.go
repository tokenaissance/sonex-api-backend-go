package config

import "github.com/zeromicro/go-zero/rest"

// Root configuration.
type Config struct {
	Debug bool `json:",env=DEBUG,default=false"`
	rest.RestConf
	Domains       []string             `json:",optional,default=[]"`
	PostgreSQL    PgConfig             `json:","`
	Redis         RedisConfig          `json:","`
	PointsService *PointsServiceConfig `json:",optional"`
	Galxe         *GalxeConfig         `json:",optional"`
	Chains        []ChainConfig        `json:","`
	Sonexes       []SonexConfig        `json:",optional"`
	PreDeposits   []PreDepositConfig   `json:",optional"`
	Xpools        []XpoolConfig        `json:",optional"`
	Farms         []FarmConfig         `json:",optional"`
}

// Get general configuration of a blockchain of `chainId`.
func (config *Config) GetChainConfig(chainId int64) *ChainConfig {
	for _, c := range config.Chains {
		if c.ChainId == chainId {
			return &c
		}
	}
	return nil
}

// Get Sonex configuration of a blockchain of `chainId`.
func (config *Config) GetSonexConfig(chainId int64) *SonexConfig {
	for _, c := range config.Sonexes {
		if c.ChainId == chainId {
			return &c
		}
	}
	return nil
}

// Get PreDeposit configuration of a blockchain of `chainId`.
func (config *Config) GetPreDepositConfig(chainId int64) *PreDepositConfig {
	for _, c := range config.PreDeposits {
		if c.ChainId == chainId {
			return &c
		}
	}
	return nil
}

// Get Xpool configuration of a blockchain of `chainId`.
func (config *Config) GetXpoolConfig(chainId int64) *XpoolConfig {
	for _, c := range config.Xpools {
		if c.ChainId == chainId {
			return &c
		}
	}
	return nil
}

// Get Farm configuration of a blockchain of `chainId`.
func (config *Config) GetFarmConfig(chainId int64) *FarmConfig {
	for _, c := range config.Farms {
		if c.ChainId == chainId {
			return &c
		}
	}
	return nil
}
