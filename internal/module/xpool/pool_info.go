package xpool

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// PoolInfo 池子信息
type PoolInfo struct {
	StakingToken         common.Address
	StakingTokenName     string
	StakingTokenSymbol   string
	StakingTokenDecimals int64
	StakingTokenLogo     string
	RewardToken          common.Address
	RewardTokenName      string
	RewardTokenSymbol    string
	RewardTokenDecimals  int64
	RewardTokenLogo      string
	RewardPerBlock       *big.Int
	TotalStaked          *big.Int
	StakedNumber         *big.Int
	IsDestroied          bool
}

// AddressPoolInfo 用户在池子中的质押状态信息
type AddressPoolInfo struct {
	Address      common.Address // 用户钱包地址
	EarnedReward *big.Int       // 用户已获得的奖励
	Height       *big.Int       // 用户获得奖励的区块高度
	Staked       *big.Int       // 用户已质押的代币数量
	Claimed      *big.Int       // 用户已领取的奖励
}
