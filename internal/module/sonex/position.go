package sonex

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Position struct {
	Factory                  common.Address
	TokenId                  *big.Int
	Slot0                    *Slot0
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

func (position *Position) Pool() common.Address {
	return ComputePoolAddress(position.Factory, position.Token0, position.Token1, position.Fee)
}

func (position *Position) Amounts() (*big.Int, *big.Int) {
	Q96 := new(big.Float).SetInt(new(big.Int).Lsh(big.NewInt(1), 96))

	liquidity := new(big.Float).SetInt(position.Liquidity)
	sqrtPriceA := Sqrt(Pow(big.NewFloat(1.0001), new(big.Float).SetInt(position.TickLower)))
	sqrtPriceB := Sqrt(Pow(big.NewFloat(1.0001), new(big.Float).SetInt(position.TickUpper)))
	currentTick := GetTickAtSqrtPrice(position.Slot0.SqrtPriceX96)
	sqrtPrice := new(big.Float).Quo(new(big.Float).SetInt(position.Slot0.SqrtPriceX96), Q96)

	var amount0 *big.Int
	var amount1 *big.Int

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
