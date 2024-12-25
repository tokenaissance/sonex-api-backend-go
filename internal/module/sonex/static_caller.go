package sonex

import (
	"context"
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/config"
	contract_common "sonex_api_backend_go/internal/lib/contract/common"
	"sonex_api_backend_go/internal/lib/contract/uniswap_v3"
	"sonex_api_backend_go/internal/lib/node_caller"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type StaticCaller struct {
	c config.Config
}

func NewStaticCaller(config config.Config) (*StaticCaller, error) {
	return &StaticCaller{
		c: config,
	}, nil
}

func (staticCaller *StaticCaller) BatchGetPools(ctx context.Context, chainId int64, pools []common.Address) ([]*Pool, error) {
	chain := staticCaller.c.GetChainConfig(chainId)
	if chain == nil {
		return nil, fmt.Errorf("%s", "Miss chain configure")
	}
	if chain.Multicall == "" || !common.IsHexAddress(chain.Multicall) {
		return nil, fmt.Errorf("%s", "Invalid multicall contract")
	}

	nodeCaller, err := node_caller.NewNodeCaller(*chain)
	if err != nil {
		return nil, err
	}
	multicallABI, err := abi.JSON(strings.NewReader(contract_common.MulticallABI))
	if err != nil {
		return nil, err
	}
	uniswapV3PoolABI, err := abi.JSON(strings.NewReader(uniswap_v3.UniswapV3PoolABI))
	if err != nil {
		return nil, err
	}
	var multicallResult struct {
		BlockNumber *big.Int `json:"blockNumber"`
		ReturnData  [][]byte `json:"returnData"`
	}

	callDatas := []contract_common.Multicall3Call{}
	slot0Data, err := uniswapV3PoolABI.Pack("slot0")
	if err != nil {
		return nil, err
	}
	feeGrowthGlobal0X128Data, err := uniswapV3PoolABI.Pack("feeGrowthGlobal0X128")
	if err != nil {
		return nil, err
	}
	feeGrowthGlobal1X128Data, err := uniswapV3PoolABI.Pack("feeGrowthGlobal1X128")
	if err != nil {
		return nil, err
	}
	for _, pool := range pools {
		callDatas = append(callDatas, contract_common.Multicall3Call{
			Target:   pool,
			CallData: slot0Data,
		})
		callDatas = append(callDatas, contract_common.Multicall3Call{
			Target:   pool,
			CallData: feeGrowthGlobal0X128Data,
		})
		callDatas = append(callDatas, contract_common.Multicall3Call{
			Target:   pool,
			CallData: feeGrowthGlobal1X128Data,
		})
	}
	data, err := multicallABI.Pack("aggregate", callDatas)
	if err != nil {
		return nil, err
	}
	err = nodeCaller.StaticCall(ctx, multicallABI, common.HexToAddress(chain.Multicall), "aggregate", data, &multicallResult)
	if err != nil {
		return nil, err
	}
	result := []*Pool{}
	for i, addr := range pools {
		var slot0 Slot0
		err = uniswapV3PoolABI.UnpackIntoInterface(&slot0, "slot0", multicallResult.ReturnData[i*3])
		if err != nil {
			return nil, err
		}
		var feeGrowthGlobal0X128 *big.Int
		err = uniswapV3PoolABI.UnpackIntoInterface(&feeGrowthGlobal0X128, "feeGrowthGlobal0X128", multicallResult.ReturnData[i*3+1])
		if err != nil {
			return nil, err
		}
		var feeGrowthGlobal1X128 *big.Int
		err = uniswapV3PoolABI.UnpackIntoInterface(&feeGrowthGlobal1X128, "feeGrowthGlobal1X128", multicallResult.ReturnData[i*3+2])
		if err != nil {
			return nil, err
		}
		pool := &Pool{
			Address:              addr,
			Slot0:                &slot0,
			FeeGrowthGlobal0X128: feeGrowthGlobal0X128,
			FeeGrowthGlobal1X128: feeGrowthGlobal1X128,
		}
		result = append(result, pool)
	}

	return result, nil
}

func (staticCaller *StaticCaller) BatchGetPositions(ctx context.Context, chainId int64, positionManager common.Address, tokenIds []*big.Int) ([]*Position, error) {
	chain := staticCaller.c.GetChainConfig(chainId)
	if chain == nil {
		return nil, fmt.Errorf("%s", "Miss chain configure")
	}
	if chain.Multicall == "" || !common.IsHexAddress(chain.Multicall) {
		return nil, fmt.Errorf("%s", "Invalid multicall contract")
	}

	nodeCaller, err := node_caller.NewNodeCaller(*chain)
	if err != nil {
		return nil, err
	}
	multicallABI, err := abi.JSON(strings.NewReader(contract_common.MulticallABI))
	if err != nil {
		return nil, err
	}
	positionManagerABI, err := abi.JSON(strings.NewReader(uniswap_v3.NonfungiblePositionManagerABI))
	if err != nil {
		return nil, err
	}
	uniswapV3PoolABI, err := abi.JSON(strings.NewReader(uniswap_v3.UniswapV3PoolABI))
	if err != nil {
		return nil, err
	}
	var multicallResult struct {
		BlockNumber *big.Int `json:"blockNumber"`
		ReturnData  [][]byte `json:"returnData"`
	}

	// Get positions
	callDatas := []contract_common.Multicall3Call{}
	data, err := positionManagerABI.Pack("factory")
	if err != nil {
		return nil, err
	}
	callDatas = append(callDatas, contract_common.Multicall3Call{
		Target:   positionManager,
		CallData: data,
	})
	for _, tokenId := range tokenIds {
		positionsData, err := positionManagerABI.Pack("positions", tokenId)
		if err != nil {
			return nil, err
		}
		callDatas = append(callDatas, contract_common.Multicall3Call{
			Target:   positionManager,
			CallData: positionsData,
		})
	}
	data, err = multicallABI.Pack("aggregate", callDatas)
	if err != nil {
		return nil, err
	}
	err = nodeCaller.StaticCall(ctx, multicallABI, common.HexToAddress(chain.Multicall), "aggregate", data, &multicallResult)
	if err != nil {
		return nil, err
	}
	var factory common.Address
	err = positionManagerABI.UnpackIntoInterface(&factory, "factory", multicallResult.ReturnData[0])
	if err != nil {
		return nil, err
	}
	positions := []*Position{}
	for i, tokenId := range tokenIds {
		var position Position
		err = positionManagerABI.UnpackIntoInterface(&position, "positions", multicallResult.ReturnData[i+1])
		if err != nil {
			return nil, err
		}
		position.Factory = factory
		position.TokenId = tokenId
		positions = append(positions, &position)
	}

	// Get pool
	poolMap := map[string]*Pool{}
	callDatas = []contract_common.Multicall3Call{}
	for _, position := range positions {
		poolAddress := ComputePoolAddress(position.Factory, position.Token0, position.Token1, position.Fee)
		poolMap[poolAddress.Hex()] = nil

		ticksData, err := uniswapV3PoolABI.Pack("ticks", position.TickLower)
		if err != nil {
			return nil, err
		}
		callDatas = append(callDatas, contract_common.Multicall3Call{
			Target:   poolAddress,
			CallData: ticksData,
		})

		ticksData, err = uniswapV3PoolABI.Pack("ticks", position.TickUpper)
		if err != nil {
			return nil, err
		}
		callDatas = append(callDatas, contract_common.Multicall3Call{
			Target:   poolAddress,
			CallData: ticksData,
		})
	}
	data, err = multicallABI.Pack("aggregate", callDatas)
	if err != nil {
		return nil, err
	}
	err = nodeCaller.StaticCall(ctx, multicallABI, common.HexToAddress(chain.Multicall), "aggregate", data, &multicallResult)
	if err != nil {
		return nil, err
	}

	poolAddresses := []common.Address{}
	for key := range poolMap {
		poolAddresses = append(poolAddresses, common.HexToAddress(key))
	}
	pools, err := staticCaller.BatchGetPools(ctx, chainId, poolAddresses)
	if err != nil {
		return nil, err
	}
	for i, pool := range pools {
		poolMap[pool.Address.Hex()] = pools[i]
	}
	for i, position := range positions {
		poolAddress := ComputePoolAddress(position.Factory, position.Token0, position.Token1, position.Fee)
		position.Pool = poolMap[poolAddress.Hex()]

		var lowerTick Tick
		err = uniswapV3PoolABI.UnpackIntoInterface(&lowerTick, "ticks", multicallResult.ReturnData[i*2])
		if err != nil {
			return nil, err
		}
		position.LowerTick = &lowerTick

		var upperTick Tick
		err = uniswapV3PoolABI.UnpackIntoInterface(&upperTick, "ticks", multicallResult.ReturnData[i*2+1])
		if err != nil {
			return nil, err
		}
		position.UpperTick = &upperTick
	}

	return positions, nil
}
