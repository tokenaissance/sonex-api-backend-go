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

func (position *Position) CurrentTick() *big.Int {
	Q96 := new(big.Float).SetInt(new(big.Int).Lsh(big.NewInt(1), 96))

	sqrtPrice := new(big.Float).Quo(new(big.Float).SetInt(position.Pool.Slot0.SqrtPriceX96), Q96)
	price := Pow(sqrtPrice, big.NewFloat(2))
	value, _ := new(big.Float).Quo(Log(price), Log(big.NewFloat(1.0001))).Int(nil)
	return value
}

func (position *Position) Amounts() (*big.Int, *big.Int) {
	Q96 := new(big.Float).SetInt(new(big.Int).Lsh(big.NewInt(1), 96))
	liquidity := new(big.Float).SetInt(position.Liquidity)
	sqrtPriceA := Sqrt(Pow(big.NewFloat(1.0001), new(big.Float).SetInt(position.TickLower)))
	sqrtPriceB := Sqrt(Pow(big.NewFloat(1.0001), new(big.Float).SetInt(position.TickUpper)))
	currentTick := position.CurrentTick()
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
