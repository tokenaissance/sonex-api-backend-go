package sonex

import (
	"math/big"
	"sonex_api_backend_go/internal/types"

	"github.com/ethereum/go-ethereum/common"
)

type Pool struct {
	Address              common.Address
	Slot0                *Slot0
	FeeGrowthGlobal0X128 *big.Int
	FeeGrowthGlobal1X128 *big.Int
	Fee                  *big.Int
}

type GetPoolListParams struct {
	Order           string
	Page            int64
	Take            int64
	ChainId         int64
	Keywords        string
	Token0Address   string
	Token1Address   string
	PositionManager string
}

type PoolFromDB struct {
	ChainId int64
	Id      string
	Address string
	FeeTier int64
	Token0  *types.PoolsResponse_Token
	Token1  *types.PoolsResponse_Token
}
