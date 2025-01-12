package sonex

import (
	"context"
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/config"
	lib_common "sonex_api_backend_go/internal/lib/common"
	contract_common "sonex_api_backend_go/internal/lib/contract/common"
	"sonex_api_backend_go/internal/lib/contract/uniswap_v3"
	"sonex_api_backend_go/internal/lib/node_caller"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// Implements adapter of the SONEX dApp
type StaticCaller struct {
	c config.Config
}

func NewStaticCaller(config config.Config) (*StaticCaller, error) {
	return &StaticCaller{
		c: config,
	}, nil
}

// Batch query ERC20 balances for given accounts and tokens
func (staticCaller *StaticCaller) BatchBalanceOf(ctx context.Context, chainId int64, addresses []common.Address, tokens []common.Address) (map[string]map[string]*big.Int, error) {
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
	erc20ABI, err := abi.JSON(strings.NewReader(contract_common.Erc20ABI))
	if err != nil {
		return nil, err
	}

	result := map[string]map[string]*big.Int{}

	callData := []contract_common.Multicall3Call3{}
	for _, address := range addresses {
		result[address.Hex()] = map[string]*big.Int{}
		balanceOfData, err := erc20ABI.Pack("balanceOf", address)
		if err != nil {
			return nil, err
		}
		for _, token := range tokens {
			callData = append(callData, contract_common.Multicall3Call3{
				Target:       token,
				AllowFailure: true,
				CallData:     balanceOfData,
			})
		}
	}

	multicallData, err := multicallABI.Pack("aggregate3", callData)
	if err != nil {
		return nil, err
	}
	var multicallResults []contract_common.Multicall3Result
	err = nodeCaller.StaticCall(ctx, multicallABI, common.HexToAddress(chain.Multicall), "aggregate3", multicallData, &multicallResults)
	if err != nil {
		return nil, err
	}

	for i, multicallResult := range multicallResults {
		if !multicallResult.Success {
			continue
		}
		address := addresses[i/len(tokens)]
		token := tokens[i%len(tokens)]
		var balance *big.Int
		err = erc20ABI.UnpackIntoInterface(&balance, "balanceOf", multicallResult.ReturnData)
		if err != nil {
			return nil, err
		}
		result[address.Hex()][token.Hex()] = balance
	}

	return result, nil
}

// Batch query pools' states
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
	var multicallResults []contract_common.Multicall3Result

	callDatas := []contract_common.Multicall3Call3{}
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
	feeData, err := uniswapV3PoolABI.Pack("fee")
	if err != nil {
		return nil, err
	}
	liquidityData, err := uniswapV3PoolABI.Pack("liquidity")
	if err != nil {
		return nil, err
	}
	for _, pool := range pools {
		callDatas = append(callDatas, contract_common.Multicall3Call3{
			Target:       pool,
			AllowFailure: true,
			CallData:     slot0Data,
		})
		callDatas = append(callDatas, contract_common.Multicall3Call3{
			Target:       pool,
			AllowFailure: true,
			CallData:     feeGrowthGlobal0X128Data,
		})
		callDatas = append(callDatas, contract_common.Multicall3Call3{
			Target:       pool,
			AllowFailure: true,
			CallData:     feeGrowthGlobal1X128Data,
		})
		callDatas = append(callDatas, contract_common.Multicall3Call3{
			Target:       pool,
			AllowFailure: true,
			CallData:     feeData,
		})
		callDatas = append(callDatas, contract_common.Multicall3Call3{
			Target:       pool,
			AllowFailure: true,
			CallData:     liquidityData,
		})
	}
	data, err := multicallABI.Pack("aggregate3", callDatas)
	if err != nil {
		return nil, err
	}
	err = nodeCaller.StaticCall(ctx, multicallABI, common.HexToAddress(chain.Multicall), "aggregate3", data, &multicallResults)
	if err != nil {
		return nil, err
	}
	result := []*Pool{}
	fieldNumber := 5
	for i, addr := range pools {
		multicallResult := multicallResults[i*fieldNumber]
		if !multicallResult.Success {
			_, errMsg := lib_common.DecodeRevert(multicallResult.ReturnData)
			return nil, fmt.Errorf("%s", fmt.Sprintf("Multicall slot0 failed: %s, pool: %s", errMsg, addr))
		}
		var slot0 Slot0
		err = uniswapV3PoolABI.UnpackIntoInterface(&slot0, "slot0", multicallResult.ReturnData)
		if err != nil {
			return nil, err
		}

		multicallResult = multicallResults[i*fieldNumber+1]
		if !multicallResult.Success {
			_, errMsg := lib_common.DecodeRevert(multicallResult.ReturnData)
			return nil, fmt.Errorf("%s", fmt.Sprintf("Multicall feeGrowthGlobal0X128 failed: %s, pool: %s", errMsg, addr))
		}
		var feeGrowthGlobal0X128 *big.Int
		err = uniswapV3PoolABI.UnpackIntoInterface(&feeGrowthGlobal0X128, "feeGrowthGlobal0X128", multicallResult.ReturnData)
		if err != nil {
			return nil, err
		}

		multicallResult = multicallResults[i*fieldNumber+2]
		if !multicallResult.Success {
			_, errMsg := lib_common.DecodeRevert(multicallResult.ReturnData)
			return nil, fmt.Errorf("%s", fmt.Sprintf("Multicall feeGrowthGlobal1X128 failed: %s, pool: %s", errMsg, addr))
		}
		var feeGrowthGlobal1X128 *big.Int
		err = uniswapV3PoolABI.UnpackIntoInterface(&feeGrowthGlobal1X128, "feeGrowthGlobal1X128", multicallResult.ReturnData)
		if err != nil {
			return nil, err
		}

		multicallResult = multicallResults[i*fieldNumber+3]
		if !multicallResult.Success {
			_, errMsg := lib_common.DecodeRevert(multicallResult.ReturnData)
			return nil, fmt.Errorf("%s", fmt.Sprintf("Multicall fee failed: %s, pool: %s", errMsg, addr))
		}
		var fee *big.Int
		err = uniswapV3PoolABI.UnpackIntoInterface(&fee, "fee", multicallResult.ReturnData)
		if err != nil {
			return nil, err
		}

		multicallResult = multicallResults[i*fieldNumber+4]
		if !multicallResult.Success {
			_, errMsg := lib_common.DecodeRevert(multicallResult.ReturnData)
			return nil, fmt.Errorf("%s", fmt.Sprintf("Multicall liquidity failed: %s, pool: %s", errMsg, addr))
		}
		var liquidity *big.Int
		err = uniswapV3PoolABI.UnpackIntoInterface(&liquidity, "liquidity", multicallResult.ReturnData)
		if err != nil {
			return nil, err
		}

		pool := &Pool{
			Address:              addr,
			Slot0:                &slot0,
			FeeGrowthGlobal0X128: feeGrowthGlobal0X128,
			FeeGrowthGlobal1X128: feeGrowthGlobal1X128,
			Fee:                  fee,
			Liquidity:            liquidity,
		}
		result = append(result, pool)
	}

	return result, nil
}

// Batch get the positions about specified NFT token IDs.
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
	var multicallResults []contract_common.Multicall3Result

	// Get positions
	callDatas := []contract_common.Multicall3Call3{}
	data, err := positionManagerABI.Pack("factory")
	if err != nil {
		return nil, err
	}
	callDatas = append(callDatas, contract_common.Multicall3Call3{
		Target:       positionManager,
		AllowFailure: true,
		CallData:     data,
	})
	for _, tokenId := range tokenIds {
		positionsData, err := positionManagerABI.Pack("positions", tokenId)
		if err != nil {
			return nil, err
		}
		callDatas = append(callDatas, contract_common.Multicall3Call3{
			Target:       positionManager,
			AllowFailure: true,
			CallData:     positionsData,
		})

		ownerOfData, err := positionManagerABI.Pack("ownerOf", tokenId)
		if err != nil {
			return nil, err
		}
		callDatas = append(callDatas, contract_common.Multicall3Call3{
			Target:       positionManager,
			AllowFailure: true,
			CallData:     ownerOfData,
		})
	}
	data, err = multicallABI.Pack("aggregate3", callDatas)
	if err != nil {
		return nil, err
	}
	err = nodeCaller.StaticCall(ctx, multicallABI, common.HexToAddress(chain.Multicall), "aggregate3", data, &multicallResults)
	if err != nil {
		return nil, err
	}
	multicallResult := multicallResults[0]
	if !multicallResult.Success {
		_, errMsg := lib_common.DecodeRevert(multicallResult.ReturnData)
		return nil, fmt.Errorf("%s", fmt.Sprintf("Multicall failed: %s", errMsg))
	}
	var factory common.Address
	err = positionManagerABI.UnpackIntoInterface(&factory, "factory", multicallResult.ReturnData)
	if err != nil {
		return nil, err
	}

	positions := []*Position{}
	fieldNumber := 2
	for i, tokenId := range tokenIds {
		multicallResult := multicallResults[i*fieldNumber+1]
		if !multicallResult.Success {
			_, errMsg := lib_common.DecodeRevert(multicallResult.ReturnData)
			return nil, fmt.Errorf("%s", fmt.Sprintf("Multicall failed: %s, token id: %s", errMsg, tokenId.String()))
		}
		var position Position
		err = positionManagerABI.UnpackIntoInterface(&position, "positions", multicallResult.ReturnData)
		if err != nil {
			return nil, err
		}

		multicallResult = multicallResults[i*fieldNumber+2]
		if !multicallResult.Success {
			_, errMsg := lib_common.DecodeRevert(multicallResult.ReturnData)
			return nil, fmt.Errorf("%s", fmt.Sprintf("Multicall failed: %s, token id: %s", errMsg, tokenId.String()))
		}
		var owner common.Address
		err = positionManagerABI.UnpackIntoInterface(&owner, "ownerOf", multicallResult.ReturnData)
		if err != nil {
			return nil, err
		}

		position.Factory = factory
		position.Owner = owner
		position.TokenId = tokenId
		positions = append(positions, &position)
	}

	// Get pool
	poolMap := map[string]*Pool{}
	callDatas = []contract_common.Multicall3Call3{}
	for _, position := range positions {
		poolAddress := ComputePoolAddress(position.Factory, position.Token0, position.Token1, position.Fee)
		poolMap[poolAddress.Hex()] = nil

		ticksData, err := uniswapV3PoolABI.Pack("ticks", position.TickLower)
		if err != nil {
			return nil, err
		}
		callDatas = append(callDatas, contract_common.Multicall3Call3{
			Target:       poolAddress,
			AllowFailure: true,
			CallData:     ticksData,
		})

		ticksData, err = uniswapV3PoolABI.Pack("ticks", position.TickUpper)
		if err != nil {
			return nil, err
		}
		callDatas = append(callDatas, contract_common.Multicall3Call3{
			Target:       poolAddress,
			AllowFailure: true,
			CallData:     ticksData,
		})
	}
	data, err = multicallABI.Pack("aggregate3", callDatas)
	if err != nil {
		return nil, err
	}
	err = nodeCaller.StaticCall(ctx, multicallABI, common.HexToAddress(chain.Multicall), "aggregate3", data, &multicallResults)
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

		multicallResult := multicallResults[i*2]
		if !multicallResult.Success {
			_, errMsg := lib_common.DecodeRevert(multicallResult.ReturnData)
			return nil, fmt.Errorf("%s", fmt.Sprintf("Multicall failed: %s, token id: %s", errMsg, position.TokenId.String()))
		}
		var lowerTick Tick
		err = uniswapV3PoolABI.UnpackIntoInterface(&lowerTick, "ticks", multicallResult.ReturnData)
		if err != nil {
			return nil, err
		}
		position.LowerTick = &lowerTick

		multicallResult = multicallResults[i*2+1]
		if !multicallResult.Success {
			_, errMsg := lib_common.DecodeRevert(multicallResult.ReturnData)
			return nil, fmt.Errorf("%s", fmt.Sprintf("Multicall failed: %s, token id: %s", errMsg, position.TokenId.String()))
		}
		var upperTick Tick
		err = uniswapV3PoolABI.UnpackIntoInterface(&upperTick, "ticks", multicallResult.ReturnData)
		if err != nil {
			return nil, err
		}
		position.UpperTick = &upperTick
	}

	return positions, nil
}
