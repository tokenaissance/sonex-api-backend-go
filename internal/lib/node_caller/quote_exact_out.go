package node_caller

import "math/big"

type QuoteExactOut struct {
	AmountOut                   *big.Int   `json:"amountOut"`
	SqrtPriceX96AfterList       []*big.Int `json:"sqrtPriceX96AfterList"`
	InitializedTicksCrossedList []uint32   `json:"initializedTicksCrossedList"`
	GasEstimate                 *big.Int   `json:"gasEstimate"`
}
