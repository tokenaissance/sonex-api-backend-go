package sonex

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Pool struct {
	Address              common.Address
	Slot0                *Slot0
	FeeGrowthGlobal0X128 *big.Int
	FeeGrowthGlobal1X128 *big.Int
	Fee                  *big.Int
}
