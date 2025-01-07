package quoter

import "math/big"

type QuoteExactInputResult struct {
	AmountOut                   *big.Int   `json:"amountOut"`
	SqrtPriceX96AfterList       []*big.Int `json:"sqrtPriceX96AfterList"`
	InitializedTicksCrossedList []uint32   `json:"initializedTicksCrossedList"`
	GasEstimate                 *big.Int   `json:"gasEstimate"`
	Error                       error
}

type QuoteExactOutputResult struct {
	AmountIn                    *big.Int   `json:"amountIn"`
	SqrtPriceX96AfterList       []*big.Int `json:"sqrtPriceX96AfterList"`
	InitializedTicksCrossedList []uint32   `json:"initializedTicksCrossedList"`
	GasEstimate                 *big.Int   `json:"gasEstimate"`
	Error                       error
}
