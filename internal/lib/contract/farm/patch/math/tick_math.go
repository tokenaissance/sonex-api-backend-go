package math

import (
	"fmt"
	"math/big"
)

var (
	MinTick = big.NewInt(-887272)
	MaxTick = big.NewInt(887272)
	Q96     = new(big.Int).Lsh(big.NewInt(1), 96)
	Q128    = new(big.Int).Lsh(big.NewInt(1), 128)

	// Precomputed constants from TickMath.sol
	tickHexValues = []string{
		"0xfffcb933bd6fad37aa2d162d1a594001",
		"0xfff97272373d413259a46990580e213a",
		"0xfff2e50f5f656932ef12357cf3c7fdcc",
		"0xffe5caca7e10e4e61c3624eaa0941cd0",
		"0xffcb9843d60f6159c9db58835c926644",
		"0xff973b41fa98c081472e6896dfb254c0",
		"0xff2ea16466c96a3843ec78b326b52861",
		"0xfe5dee046a99a2a811c461f1969c3053",
		"0xfcbe86c7900a88aedcffc83b479aa3a4",
		"0xf987a7253ac413176f2b074cf7815e54",
		"0xf3392b0822b70005940c7a398e4b70f3",
		"0xe7159475a2c29b7443b29c7fa6e889d9",
		"0xd097f3bdfd2022b8845ad8f792aa5825",
		"0xa9f746462d870fdf8a65dc1f90e061e5",
		"0x70d869a156d2a1b890bb3df62baf32f7",
		"0x31be135f97d08fd981231505542fcfa6",
		"0x9aa508b5b7a84e1c677de54f3e99bc9",
		"0x5d6af8dedb81196699c329225ee604",
		"0x2216e584f5fa1ea926041bedfe98",
		"0x48a170391f7dc42444e8fa2",
	}
)

// GetSqrtRatioAtTick calculates sqrt(1.0001^tick) * 2^96
func GetSqrtRatioAtTick(tick *big.Int) (*big.Int, error) {
	if tick.Cmp(MinTick) < 0 || tick.Cmp(MaxTick) > 0 {
		return nil, fmt.Errorf("tick out of range")
	}

	absTick := new(big.Int)
	if tick.Sign() < 0 {
		absTick.Neg(tick)
	} else {
		absTick.Set(tick)
	}

	ratio := new(big.Int)
	if absTick.Bit(0) != 0 {
		ratio.SetString("0xfffcb933bd6fad37aa2d162d1a594001", 0)
	} else {
		ratio.SetString("0x100000000000000000000000000000000", 0)
	}

	// Bit operation loop
	for i := 0; i < len(tickHexValues); i++ {
		if absTick.Bit(int(i)) != 0 {
			tmp, ok := new(big.Int).SetString(tickHexValues[i], 0)
			if !ok {
				return nil, fmt.Errorf("failed to parse hex value at bit %d", i)
			}
			ratio.Mul(ratio, tmp)
			ratio.Rsh(ratio, 128)
		}
	}

	if tick.Sign() > 0 {
		maxUint256 := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))
		ratio = new(big.Int).Div(maxUint256, ratio)
	}

	// Q128.128 to Q96 conversion
	result := new(big.Int).Rsh(ratio, 32)
	if new(big.Int).And(ratio, new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 32), big.NewInt(1))).Sign() > 0 {
		result.Add(result, big.NewInt(1))
	}

	return result, nil
}

// GetAmount0ForLiquidity calculates amount0 for given liquidity and price range
func GetAmount0ForLiquidity(sqrtRatioAX96, sqrtRatioBX96, liquidity *big.Int) (*big.Int, error) {
	if sqrtRatioAX96.Cmp(sqrtRatioBX96) > 0 {
		sqrtRatioAX96, sqrtRatioBX96 = sqrtRatioBX96, sqrtRatioAX96
	}

	numerator := new(big.Int).Sub(sqrtRatioBX96, sqrtRatioAX96)
	numerator.Mul(numerator, liquidity)

	denominator := new(big.Int).Mul(sqrtRatioBX96, sqrtRatioAX96)
	if denominator.Sign() == 0 {
		return nil, fmt.Errorf("division by zero")
	}

	amount0 := new(big.Int).Mul(numerator, Q96)
	amount0.Div(amount0, denominator)

	return amount0, nil
}

// GetAmount1ForLiquidity calculates amount1 for given liquidity and price range
func GetAmount1ForLiquidity(sqrtRatioAX96, sqrtRatioBX96, liquidity *big.Int) (*big.Int, error) {
	if sqrtRatioAX96.Cmp(sqrtRatioBX96) > 0 {
		sqrtRatioAX96, sqrtRatioBX96 = sqrtRatioBX96, sqrtRatioAX96
	}

	numerator := new(big.Int).Sub(sqrtRatioBX96, sqrtRatioAX96)
	amount1 := new(big.Int).Mul(liquidity, numerator)
	amount1.Div(amount1, Q96)

	return amount1, nil
}
