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

type CompareableFloat struct {
	Value *big.Float
}

func (compareable CompareableFloat) CompareTo(other linq.Comparable) int {
	l, _ := other.(CompareableFloat)
	return compareable.Value.Cmp(l.Value)
}
