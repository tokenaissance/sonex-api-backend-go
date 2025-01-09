package sonex

import (
	"math/big"
	"sonex_api_backend_go/internal/types"

	"github.com/ethereum/go-ethereum/common"
)

type Position struct {
	Factory                  common.Address
	TokenId                  *big.Int
	Pool                     *Pool
	LowerTick                *Tick
	UpperTick                *Tick
	Owner                    common.Address `json:"owner"`
	Nonce                    *big.Int       `json:"nonce"`
	Operator                 common.Address `json:"operator"`
	Token0                   common.Address `json:"token0"`
	Token1                   common.Address `json:"token1"`
	Fee                      *big.Int       `json:"fee"`
	TickLower                *big.Int       `json:"tickLower"`
	TickUpper                *big.Int       `json:"tickUpper"`
	Liquidity                *big.Int       `json:"liquidity"`
	FeeGrowthInside0LastX128 *big.Int       `json:"feeGrowthInside0LastX128"`
	FeeGrowthInside1LastX128 *big.Int       `json:"feeGrowthInside1LastX128"`
	TokensOwed0              *big.Int       `json:"tokensOwed0"`
	TokensOwed1              *big.Int       `json:"tokensOwed1"`
}

func (position *Position) getFeeGrowthInsideX128() (*big.Int, *big.Int) {
	currentTick := GetTickAtSqrtPrice(position.Pool.Slot0.SqrtPriceX96)

	var feeGrowthBelow0X128 *big.Int
	var feeGrowthBelow1X128 *big.Int
	var feeGrowthAbove0X128 *big.Int
	var feeGrowthAbove1X128 *big.Int

	if currentTick.Cmp(position.TickLower) >= 0 {
		feeGrowthBelow0X128 = position.LowerTick.FeeGrowthOutside0X128
		feeGrowthBelow1X128 = position.LowerTick.FeeGrowthOutside1X128
	} else {
		feeGrowthBelow0X128 = new(big.Int).Sub(position.Pool.FeeGrowthGlobal0X128, position.LowerTick.FeeGrowthOutside0X128)
		feeGrowthBelow1X128 = new(big.Int).Sub(position.Pool.FeeGrowthGlobal1X128, position.LowerTick.FeeGrowthOutside1X128)
	}

	if currentTick.Cmp(position.TickUpper) < 0 {
		feeGrowthAbove0X128 = position.UpperTick.FeeGrowthOutside0X128
		feeGrowthAbove1X128 = position.UpperTick.FeeGrowthOutside1X128
	} else {
		feeGrowthAbove0X128 = new(big.Int).Sub(position.Pool.FeeGrowthGlobal0X128, position.UpperTick.FeeGrowthOutside0X128)
		feeGrowthAbove1X128 = new(big.Int).Sub(position.Pool.FeeGrowthGlobal1X128, position.UpperTick.FeeGrowthOutside1X128)
	}

	feeGrowthInside0X128 := new(big.Int).Sub(new(big.Int).Sub(position.Pool.FeeGrowthGlobal0X128, feeGrowthBelow0X128), feeGrowthAbove0X128)
	feeGrowthInside1X128 := new(big.Int).Sub(new(big.Int).Sub(position.Pool.FeeGrowthGlobal1X128, feeGrowthBelow1X128), feeGrowthAbove1X128)

	return feeGrowthInside0X128, feeGrowthInside1X128
}

func (position *Position) Amounts() (*big.Int, *big.Int) {
	Q96 := new(big.Float).SetInt(new(big.Int).Lsh(big.NewInt(1), 96))

	liquidity := new(big.Float).SetInt(position.Liquidity)
	sqrtPriceA := Sqrt(Pow(big.NewFloat(1.0001), new(big.Float).SetInt(position.TickLower)))
	sqrtPriceB := Sqrt(Pow(big.NewFloat(1.0001), new(big.Float).SetInt(position.TickUpper)))
	currentTick := position.Pool.Slot0.Tick
	// currentTick := GetTickAtSqrtPrice(position.Pool.Slot0.SqrtPriceX96)
	sqrtPrice := new(big.Float).Quo(new(big.Float).SetInt(position.Pool.Slot0.SqrtPriceX96), Q96)

	amount0 := big.NewInt(0)
	amount1 := big.NewInt(0)

	if currentTick.Cmp(position.TickLower) < 0 {
		amount0, _ = new(big.Float).Mul(liquidity, new(big.Float).Quo(new(big.Float).Sub(sqrtPriceB, sqrtPriceA), new(big.Float).Mul(sqrtPriceA, sqrtPriceB))).Int(nil)
	} else if currentTick.Cmp(position.TickUpper) >= 0 {
		amount1, _ = new(big.Float).Mul(liquidity, new(big.Float).Sub(sqrtPriceB, sqrtPriceA)).Int(nil)
	} else {
		amount0, _ = new(big.Float).Mul(liquidity, new(big.Float).Quo(new(big.Float).Sub(sqrtPriceB, sqrtPrice), new(big.Float).Mul(sqrtPrice, sqrtPriceB))).Int(nil)
		amount1, _ = new(big.Float).Mul(liquidity, new(big.Float).Sub(sqrtPrice, sqrtPriceA)).Int(nil)
	}

	return amount0, amount1
}

func (position *Position) UnclaimedFees() (*big.Int, *big.Int) {
	Q128 := new(big.Int).Lsh(big.NewInt(1), 128)

	feeGrowthDelta0 := new(big.Int).Sub(position.Pool.FeeGrowthGlobal0X128, position.FeeGrowthInside0LastX128)
	feeGrowthDelta1 := new(big.Int).Sub(position.Pool.FeeGrowthGlobal1X128, position.FeeGrowthInside1LastX128)

	fee0 := new(big.Int).Mul(feeGrowthDelta0, position.Liquidity)
	fee1 := new(big.Int).Mul(feeGrowthDelta1, position.Liquidity)

	fee0 = new(big.Int).Div(fee0, Q128)
	fee1 = new(big.Int).Div(fee1, Q128)

	return fee0, fee1
}

type PositionValues struct {
	Pnl            float64                      `json:"pnl"`
	PositionValues []*types.PositionValueRecord `json:"positionValues"`
	TotalPosition  float64                      `json:"totalPosition"`
	TotalClaimable float64                      `json:"totalClaimable"`
}
