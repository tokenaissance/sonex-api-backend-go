package config

import (
	"math/rand"
	"time"

	"github.com/ahmetb/go-linq/v3"
)

type ChainConfig struct {
	ChainId                 int64            `json:","`
	Name                    string           `json:","`
	Nodes                   map[string]int64 `json:","`
	MinConfirmations        int64            `json:",default=5"`
	MinDelayWaitMillisecond int64            `json:",default=2000"`
	MaxLogsNumber           int64            `json:",default=10000"`
	Multicall               string           `json:",optional"`
	Tokens                  []TokenConfig    `json:",optional"`
	WETH                    string           `json:",optional"`
}

type TokenConfig struct {
	Address      string `json:","`
	IsStablecoin bool   `json:",default=false"`
}

func (config *ChainConfig) GetNode() string {
	if len(config.Nodes) == 0 {
		return ""
	}
	totalWeight := linq.From(config.Nodes).SelectT(func(item linq.KeyValue) int64 { return item.Value.(int64) }).SumInts()
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	randValue := rand.Int63n(totalWeight)
	var cumulativeWeight int64
	for key, weight := range config.Nodes {
		cumulativeWeight += weight
		if randValue < cumulativeWeight {
			return key
		}
	}
	return ""
}
