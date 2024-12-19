package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	Debug bool `json:",env=DEBUG,default=false"`
	rest.RestConf
	PostgreSQL PgConfig      `json:","`
	Chains     []ChainConfig `json:","`
	Sonexes    []SonexConfig `json:",optional"`
	XPools     []XPoolConfig `json:",optional"`
}
