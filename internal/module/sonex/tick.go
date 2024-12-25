package sonex

import "math/big"

type Tick struct {
	LiquidityGross                 *big.Int `json:"liquidityGross"`
	LiquidityNet                   *big.Int `json:"liquidityNet"`
	FeeGrowthOutside0X128          *big.Int `json:"feeGrowthOutside0X128"`
	FeeGrowthOutside1X128          *big.Int `json:"feeGrowthOutside1X128"`
	TickCumulativeOutside          *big.Int `json:"tickCumulativeOutside"`
	SecondsPerLiquidityOutsideX128 *big.Int `json:"secondsPerLiquidityOutsideX128"`
	SecondsOutside                 uint32   `json:"secondsOutside"`
	Initialized                    bool     `json:"initialized"`
}
