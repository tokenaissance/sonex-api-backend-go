package xpool

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type PoolInfo struct {
	StakingToken         common.Address
	StakingTokenName     string
	StakingTokenSymbol   string
	StakingTokenDecimals int64
	RewardToken          common.Address
	RewardTokenName      string
	RewardTokenSymbol    string
	RewardTokenDecimals  int64
	RewardPerBlock       *big.Int
	TotalStaked          *big.Int
	StakedNumber         *big.Int
	IsDestroied          bool
}
