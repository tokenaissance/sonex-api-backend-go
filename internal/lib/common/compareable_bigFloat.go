package common

import (
	"math/big"

	"github.com/ahmetb/go-linq/v3"
)

type CompareableBigFloat struct {
	Value *big.Float
}

func (liquidity CompareableBigFloat) CompareTo(other linq.Comparable) int {
	l, _ := other.(CompareableBigFloat)
	return liquidity.Value.Cmp(l.Value)
}
