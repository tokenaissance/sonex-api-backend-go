package patch

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"

	"sonex_api_backend_go/internal/config"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	commonContract "sonex_api_backend_go/internal/lib/contract/common"
	"sonex_api_backend_go/internal/lib/contract/farm/patch/math"
	v3Contract "sonex_api_backend_go/internal/lib/contract/uniswap_v3"
	nodeCaller "sonex_api_backend_go/internal/lib/node_caller"
)

const batchSize = 50

type PositionResult struct {
	TokenId  int64
	Position *struct {
		Nonce                    *big.Int
		Operator                 common.Address
		Token0                   common.Address
		Token1                   common.Address
		Fee                      *big.Int
		TickLower                *big.Int
		TickUpper                *big.Int
		Liquidity                *big.Int
		FeeGrowthInside0LastX128 *big.Int
		FeeGrowthInside1LastX128 *big.Int
		TokensOwed0              *big.Int
		TokensOwed1              *big.Int
	}
	Error error
}

// 获取质押池中所有token的总量
func GetTotalTokensInPool(ctx context.Context, conf config.Config, chainId int64, poolAddress common.Address, tokenIds []int64) (totalToken0 *big.Int, totalToken1 *big.Int, err error) {

	totalToken0 = big.NewInt(0)
	totalToken1 = big.NewInt(0)

	positionManagerAddress := conf.GetSonexConfig(chainId).PositionManager
	chainConfig := conf.GetChainConfig(chainId)
	multicallAddress := chainConfig.Multicall
	if positionManagerAddress == "" || multicallAddress == "" {
		err = fmt.Errorf("invalid config")
		return
	}

	client, err := ethclient.Dial(chainConfig.GetNode())
	if err != nil {
		return
	}
	pool, err := v3Contract.NewUniswapV3Pool(poolAddress, client)
	if err != nil {
		err = fmt.Errorf("failed to create pool contract: %v", err)
		return
	}

	slot0, err := pool.Slot0(&bind.CallOpts{})
	if err != nil {
		err = fmt.Errorf("failed to get slot0: %v", err)
		return
	}

	positionABI, err := abi.JSON(strings.NewReader(v3Contract.NonfungiblePositionManagerMetaData.ABI))
	if err != nil {
		return
	}

	multicallABI, err := abi.JSON(strings.NewReader(commonContract.MulticallMetaData.ABI))
	if err != nil {
		return
	}

	// 批量查询positions
	resultChan := make(chan PositionResult, len(tokenIds))
	var wg sync.WaitGroup

	for i := 0; i < len(tokenIds); i += batchSize {
		end := i + batchSize
		if end > len(tokenIds) {
			end = len(tokenIds)
		}

		wg.Add(1)
		go func(batch []int64) {
			defer wg.Done()

			calls := make([]commonContract.Multicall3Call3, len(batch))
			for j, tokenId := range batch {
				data, err := positionABI.Pack("positions", big.NewInt(tokenId))
				if err != nil {
					resultChan <- PositionResult{TokenId: tokenId, Error: err}
					continue
				}
				calls[j] = commonContract.Multicall3Call3{
					Target:       common.HexToAddress(positionManagerAddress),
					AllowFailure: true,
					CallData:     data,
				}
			}

			caller, _ := nodeCaller.NewNodeCaller(*chainConfig) // 从配置获取
			var result struct {
				ReturnData []struct {
					Success    bool
					ReturnData []byte
				}
			}

			callData, _ := multicallABI.Pack("aggregate3", calls)
			err := caller.StaticCall(ctx, multicallABI, common.HexToAddress(multicallAddress), "aggregate3", callData, &result)
			if err != nil {
				for _, tokenId := range batch {
					resultChan <- PositionResult{TokenId: tokenId, Error: err}
				}
				return
			}

			for j, tokenId := range batch {
				if !result.ReturnData[j].Success {
					resultChan <- PositionResult{TokenId: tokenId, Error: fmt.Errorf("multicall positions failed")}
					continue
				}

				var position struct {
					Nonce                    *big.Int
					Operator                 common.Address
					Token0                   common.Address
					Token1                   common.Address
					Fee                      *big.Int
					TickLower                *big.Int
					TickUpper                *big.Int
					Liquidity                *big.Int
					FeeGrowthInside0LastX128 *big.Int
					FeeGrowthInside1LastX128 *big.Int
					TokensOwed0              *big.Int
					TokensOwed1              *big.Int
				}

				err = positionABI.UnpackIntoInterface(&position, "positions", result.ReturnData[j].ReturnData)
				if err != nil {
					resultChan <- PositionResult{TokenId: tokenId, Error: err}
					continue
				}
				resultChan <- PositionResult{TokenId: tokenId, Position: &position}
			}
		}(tokenIds[i:end])
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 计算Token总量
	for result := range resultChan {
		if result.Error != nil {
			return nil, nil, result.Error
		}

		pos := result.Position
		amount0, amount1, err := calculateAmounts(pos.TickLower, pos.TickUpper, pos.Liquidity, slot0.SqrtPriceX96)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to calculate amounts for position %d: %v", result.TokenId, err)
		}

		totalToken0.Add(totalToken0, amount0)
		totalToken1.Add(totalToken1, amount1)
	}

	return totalToken0, totalToken1, nil
}

// 计算单个position的token数量
func calculateAmounts(tickLower, tickUpper, liquidity, currentSqrtPrice *big.Int) (*big.Int, *big.Int, error) {
	sqrtRatioAX96, err := math.GetSqrtRatioAtTick(tickLower)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get sqrtRatioA: %v", err)
	}

	sqrtRatioBX96, err := math.GetSqrtRatioAtTick(tickUpper)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get sqrtRatioB: %v", err)
	}

	amount0, err := math.GetAmount0ForLiquidity(sqrtRatioAX96, currentSqrtPrice, liquidity)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to calculate amount0: %v", err)
	}

	amount1, err := math.GetAmount1ForLiquidity(currentSqrtPrice, sqrtRatioBX96, liquidity)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to calculate amount1: %v", err)
	}

	return amount0, amount1, nil
}
