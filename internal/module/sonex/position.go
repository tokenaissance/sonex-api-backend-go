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

func (position *Position) Amounts() (*big.Int, *big.Int) {
	Q96 := new(big.Float).SetInt(new(big.Int).Lsh(big.NewInt(1), 96)) // Q96 float
	liquidityInput := position.Liquidity
	sqrtPriceX96Input := position.Pool.Slot0.SqrtPriceX96
	tickLowerInput := position.TickLower
	tickUpperInput := position.TickUpper

	liquidity := new(big.Float).SetInt(liquidityInput)
	sqrtRatioA := Sqrt(Pow(big.NewFloat(1.0001), new(big.Float).SetInt(tickLowerInput))) // sqrtRatioA
	sqrtRatioB := Sqrt(Pow(big.NewFloat(1.0001), new(big.Float).SetInt(tickUpperInput))) // sqrtRatioB
	currentTick := GetTickAtSqrtPrice(sqrtPriceX96Input)                                 // currentTick
	sqrtPrice := new(big.Float).Quo(new(big.Float).SetInt(sqrtPriceX96Input), Q96)       // sqrtPrice

	amount0 := big.NewInt(0)
	amount1 := big.NewInt(0)

	if currentTick.Cmp(tickLowerInput) < 0 {
		amount0, _ = new(big.Float).Mul(liquidity, new(big.Float).Quo(new(big.Float).Sub(sqrtRatioB, sqrtRatioA), new(big.Float).Mul(sqrtRatioA, sqrtRatioB))).Int(nil)
	} else if currentTick.Cmp(tickUpperInput) >= 0 {
		amount1, _ = new(big.Float).Mul(liquidity, new(big.Float).Sub(sqrtRatioB, sqrtRatioA)).Int(nil)
	} else {
		amount0, _ = new(big.Float).Mul(liquidity, new(big.Float).Quo(new(big.Float).Sub(sqrtRatioB, sqrtPrice), new(big.Float).Mul(sqrtPrice, sqrtRatioB))).Int(nil)
		amount1, _ = new(big.Float).Mul(liquidity, new(big.Float).Sub(sqrtPrice, sqrtRatioA)).Int(nil)
	}

	return amount0, amount1
}

func Amounts(liquidityInput *big.Int, sqrtPriceX96Input *big.Int, tickLowerInput *big.Int, tickUpperInput *big.Int) (*big.Int, *big.Int) {
	Q96 := new(big.Float).SetInt(new(big.Int).Lsh(big.NewInt(1), 96)) // Q96 float

	liquidity := new(big.Float).SetInt(liquidityInput)
	sqrtRatioA := Sqrt(Pow(big.NewFloat(1.0001), new(big.Float).SetInt(tickLowerInput))) // sqrtRatioA
	sqrtRatioB := Sqrt(Pow(big.NewFloat(1.0001), new(big.Float).SetInt(tickUpperInput))) // sqrtRatioB
	currentTick := GetTickAtSqrtPrice(sqrtPriceX96Input)                                 // currentTick
	sqrtPrice := new(big.Float).Quo(new(big.Float).SetInt(sqrtPriceX96Input), Q96)       // sqrtPrice

	amount0 := big.NewInt(0)
	amount1 := big.NewInt(0)

	if currentTick.Cmp(tickLowerInput) < 0 {
		amount0, _ = new(big.Float).Mul(liquidity, new(big.Float).Quo(new(big.Float).Sub(sqrtRatioB, sqrtRatioA), new(big.Float).Mul(sqrtRatioA, sqrtRatioB))).Int(nil)
	} else if currentTick.Cmp(tickUpperInput) >= 0 {
		amount1, _ = new(big.Float).Mul(liquidity, new(big.Float).Sub(sqrtRatioB, sqrtRatioA)).Int(nil)
	} else {
		amount0, _ = new(big.Float).Mul(liquidity, new(big.Float).Quo(new(big.Float).Sub(sqrtRatioB, sqrtPrice), new(big.Float).Mul(sqrtPrice, sqrtRatioB))).Int(nil)
		amount1, _ = new(big.Float).Mul(liquidity, new(big.Float).Sub(sqrtPrice, sqrtRatioA)).Int(nil)
	}

	return amount0, amount1
}

func (position *Position) UnclaimedFees() (*big.Int, *big.Int) {
	Q128 := new(big.Int).Lsh(big.NewInt(1), 128)
	Q256 := new(big.Int).Lsh(big.NewInt(1), 256)

	subIn256 := func(x *big.Int, y *big.Int) *big.Int {
		difference := new(big.Int).Sub(x, y)
		if difference.Sign() < 0 {
			difference = new(big.Int).Add(difference, Q256)
		}
		return difference
	}

	currentTick := position.Pool.Slot0.Tick
	var tickLowerFeeGrowthBelow0 *big.Int
	var tickLowerFeeGrowthBelow1 *big.Int
	var tickUpperFeeGrowthAbove0 *big.Int
	var tickUpperFeeGrowthAbove1 *big.Int

	if currentTick.Cmp(position.TickUpper) >= 0 {
		tickUpperFeeGrowthAbove0 = subIn256(position.Pool.FeeGrowthGlobal0X128, position.UpperTick.FeeGrowthOutside0X128)
		tickUpperFeeGrowthAbove1 = subIn256(position.Pool.FeeGrowthGlobal1X128, position.UpperTick.FeeGrowthOutside1X128)
	} else {
		tickUpperFeeGrowthAbove0 = position.UpperTick.FeeGrowthOutside0X128
		tickUpperFeeGrowthAbove1 = position.UpperTick.FeeGrowthOutside1X128
	}

	if currentTick.Cmp(position.TickLower) >= 0 {
		tickLowerFeeGrowthBelow0 = position.LowerTick.FeeGrowthOutside0X128
		tickLowerFeeGrowthBelow1 = position.LowerTick.FeeGrowthOutside1X128
	} else {
		tickLowerFeeGrowthBelow0 = subIn256(position.Pool.FeeGrowthGlobal0X128, position.LowerTick.FeeGrowthOutside0X128)
		tickLowerFeeGrowthBelow1 = subIn256(position.Pool.FeeGrowthGlobal1X128, position.LowerTick.FeeGrowthOutside1X128)
	}

	frT10 := subIn256(subIn256(position.Pool.FeeGrowthGlobal0X128, tickLowerFeeGrowthBelow0), tickUpperFeeGrowthAbove0)
	frT11 := subIn256(subIn256(position.Pool.FeeGrowthGlobal1X128, tickLowerFeeGrowthBelow1), tickUpperFeeGrowthAbove1)

	fee0 := new(big.Int).Div(new(big.Int).Mul(subIn256(frT10, position.FeeGrowthInside0LastX128), position.Liquidity), Q128)
	fee1 := new(big.Int).Div(new(big.Int).Mul(subIn256(frT11, position.FeeGrowthInside1LastX128), position.Liquidity), Q128)

	return fee0, fee1
}

type PositionValues struct {
	Pnl            float64                      `json:"pnl"`
	PositionValues []*types.PositionValueRecord `json:"positionValues"`
	TotalPosition  float64                      `json:"totalPosition"`
	TotalClaimable float64                      `json:"totalClaimable"`
}
