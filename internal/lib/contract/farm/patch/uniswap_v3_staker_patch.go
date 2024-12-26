package patch

import (
	"context"
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/lib/contract/farm/patch/math"
	v3PoolContract "sonex_api_backend_go/internal/lib/contract/uniswap_v3"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeromicro/go-zero/core/logx"
)

// ComputeTotalTokensInPool calculates total amount of token0 and token1 in the incentive pool
func GetTotalTokensInPool(ctx context.Context, ethClient *ethclient.Client, positionManagerAddress, poolAddress common.Address, tokenIds []int64) (*big.Int, *big.Int, error) {

	// tokenIds
	logx.Info("tokenIds: ", tokenIds)
	logx.Info("positionManagerAddress: ", positionManagerAddress)
	logx.Info("poolAddress: ", poolAddress)

	totalToken0 := big.NewInt(0)
	totalToken1 := big.NewInt(0)

	pool, err := v3PoolContract.NewUniswapV3Pool(poolAddress, ethClient)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create pool contract: %v", err)
	}

	positionManager, err := v3PoolContract.NewNonfungiblePositionManager(positionManagerAddress, ethClient)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create position manager: %v", err)
	}

	slot0, err := pool.Slot0(&bind.CallOpts{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get slot0: %v", err)
	}

	for _, tokenId := range tokenIds {
		position, err := positionManager.Positions(&bind.CallOpts{}, big.NewInt(tokenId))
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get position %d: %v", tokenId, err)
		}

		sqrtRatioAX96, err := math.GetSqrtRatioAtTick(position.TickLower)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get sqrtRatioA for position %d: %v", tokenId, err)
		}

		sqrtRatioBX96, err := math.GetSqrtRatioAtTick(position.TickUpper)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get sqrtRatioB for position %d: %v", tokenId, err)
		}

		amount0, err := math.GetAmount0ForLiquidity(sqrtRatioAX96, slot0.SqrtPriceX96, position.Liquidity)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to calculate amount0 for position %d: %v", tokenId, err)
		}
		totalToken0.Add(totalToken0, amount0)

		amount1, err := math.GetAmount1ForLiquidity(slot0.SqrtPriceX96, sqrtRatioBX96, position.Liquidity)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to calculate amount1 for position %d: %v", tokenId, err)
		}
		totalToken1.Add(totalToken1, amount1)
	}

	return totalToken0, totalToken1, nil
}
