package common

import (
	"math/big"

	"github.com/ahmetb/go-linq/v3"
)

type CompareableBigInt struct {
	Value *big.Int
}

func (liquidity CompareableBigInt) CompareTo(other linq.Comparable) int {
	l, _ := other.(CompareableBigInt)
	return liquidity.Value.Cmp(l.Value)
}
