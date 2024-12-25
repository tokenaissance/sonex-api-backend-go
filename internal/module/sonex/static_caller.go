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

func (staticCaller *StaticCaller) BatchGetSlot0s(ctx context.Context, chainId int64, pools []common.Address) ([]*Slot0, error) {
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
	data, err := uniswapV3PoolABI.Pack("slot0")
	if err != nil {
		return nil, err
	}
	for _, pool := range pools {
		callDatas = append(callDatas, contract_common.Multicall3Call{
			Target:   pool,
			CallData: data,
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
	slot0s := []*Slot0{}
	for i, pool := range pools {
		var slot0 Slot0
		err = uniswapV3PoolABI.UnpackIntoInterface(&slot0, "slot0", multicallResult.ReturnData[i])
		if err != nil {
			return nil, err
		}
		slot0.Pool = pool
		slot0s = append(slot0s, &slot0)
	}

	return slot0s, nil
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
		data, err = positionManagerABI.Pack("positions", tokenId)
		if err != nil {
			return nil, err
		}
		callDatas = append(callDatas, contract_common.Multicall3Call{
			Target:   positionManager,
			CallData: data,
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

	// Get slot0
	poolMap := map[string]*Slot0{}
	for _, position := range positions {
		poolMap[position.Pool().Hex()] = nil
	}
	pools := []common.Address{}
	for key := range poolMap {
		pools = append(pools, common.HexToAddress(key))
	}
	slot0s, err := staticCaller.BatchGetSlot0s(ctx, chainId, pools)
	if err != nil {
		return nil, err
	}
	for i, pool := range pools {
		poolMap[pool.Hex()] = slot0s[i]
	}
	for _, position := range positions {
		position.Slot0 = poolMap[position.Pool().Hex()]
	}

	return positions, nil
}
