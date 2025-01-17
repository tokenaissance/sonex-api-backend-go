package quoter

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/config"
	lib_common "sonex_api_backend_go/internal/lib/common"
	contract_common "sonex_api_backend_go/internal/lib/contract/common"
	"sonex_api_backend_go/internal/lib/contract/uniswap_v3"
	"sonex_api_backend_go/internal/lib/node_caller"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Implements adapter of the "Quoter" contract.
type Quoter struct {
	chainConf config.ChainConfig
	sonexConf config.SonexConfig
}

func NewQuoter(chainConf config.ChainConfig, sonexConf config.SonexConfig) (*Quoter, error) {
	if sonexConf.Quoter == "" || !common.IsHexAddress(sonexConf.Quoter) {
		return nil, fmt.Errorf("%s", "Invalid quoter address")
	}

	return &Quoter{
		chainConf: chainConf,
		sonexConf: sonexConf,
	}, nil
}

// Finds available paths to swapping `from` token for `to` token
func FindPaths(pools []Pool, from common.Address, to common.Address, maxLength int64) []Path {
	contains := func(tokens []common.Address, address common.Address) bool {
		for _, token := range tokens {
			if token == address {
				return true
			}
		}
		return false
	}

	graph := make(map[string][]Pool)
	for _, pool := range pools {
		graph[pool.Token0.Hex()] = append(graph[pool.Token0.Hex()], pool)
		graph[pool.Token1.Hex()] = append(graph[pool.Token1.Hex()], pool)
	}

	var paths []Path
	visited := make(map[string]bool)

	var dfs func(currentToken common.Address, path Path)
	dfs = func(currentToken common.Address, path Path) {
		visited[currentToken.Hex()] = true

		if currentToken.Hex() == to.Hex() {
			pathCopy := Path{
				Pools:  append([]Pool{}, path.Pools...),
				Tokens: append([]common.Address{}, path.Tokens...),
				Fees:   append([]int64{}, path.Fees...),
			}
			paths = append(paths, pathCopy)
		} else {
			for _, pool := range graph[currentToken.Hex()] {
				var nextToken common.Address
				var fee int64

				if pool.Token0 == currentToken {
					nextToken = pool.Token1
					fee = pool.Fee
				} else {
					nextToken = pool.Token0
					fee = pool.Fee
				}

				if !visited[nextToken.Hex()] && (maxLength == 0 || int64(len(path.Pools)) < maxLength) {
					if !contains(path.Tokens, nextToken) {
						path.Pools = append(path.Pools, pool)
						path.Tokens = append(path.Tokens, nextToken)
						path.Fees = append(path.Fees, fee)

						dfs(nextToken, path)

						path.Pools = path.Pools[:len(path.Pools)-1]
						path.Tokens = path.Tokens[:len(path.Tokens)-1]
						path.Fees = path.Fees[:len(path.Fees)-1]
					}
				}
			}
		}

		visited[currentToken.Hex()] = false
	}
	dfs(from, Path{Tokens: []common.Address{from}, Fees: []int64{}})

	return paths
}

// Returns the amount out received for a given exact input swap without executing the swap.
func (quoter *Quoter) QuoteExactInput(ctx context.Context, path []byte, amountIn *big.Int) (*QuoteExactInputResult, error) {
	quoterAddress := common.HexToAddress(quoter.sonexConf.Quoter)

	abi, err := abi.JSON(strings.NewReader(uniswap_v3.IquoterV2MetaData.ABI))
	if err != nil {
		return nil, err
	}
	data, err := abi.Pack("quoteExactInput", path, amountIn)
	if err != nil {
		return nil, err
	}
	msg := ethereum.CallMsg{
		To:   &quoterAddress,
		Data: data,
	}
	client, err := ethclient.DialContext(ctx, quoter.chainConf.GetNode())
	if err != nil {
		return nil, err
	}
	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, err
	}

	var output QuoteExactInputResult
	err = abi.UnpackIntoInterface(&output, "quoteExactInput", result)
	if err != nil {
		return nil, err
	}
	return &output, nil
}

// Batch version of [Quoter.QuoteExactInput]
func (quoter *Quoter) BatchQuoteExactInputs(ctx context.Context, paths [][]byte, amountIn *big.Int) (map[string]*QuoteExactInputResult, error) {
	quoterAddress := common.HexToAddress(quoter.sonexConf.Quoter)

	quoterABI, err := abi.JSON(strings.NewReader(uniswap_v3.IquoterV2MetaData.ABI))
	if err != nil {
		return nil, err
	}
	multicallABI, err := abi.JSON(strings.NewReader(contract_common.MulticallABI))
	if err != nil {
		return nil, err
	}
	nodeCaller, err := node_caller.NewNodeCaller(quoter.chainConf)
	if err != nil {
		return nil, err
	}
	var multicallResults []contract_common.Multicall3Result

	multicallData := []contract_common.Multicall3Call3{}
	for _, path := range paths {
		data, err := quoterABI.Pack("quoteExactInput", path, amountIn)
		if err != nil {
			return nil, err
		}
		multicallData = append(multicallData, contract_common.Multicall3Call3{
			Target:       quoterAddress,
			AllowFailure: true,
			CallData:     data,
		})
	}
	data, err := multicallABI.Pack("aggregate3", multicallData)
	if err != nil {
		return nil, err
	}
	err = nodeCaller.StaticCall(ctx, multicallABI, common.HexToAddress(quoter.chainConf.Multicall), "aggregate3", data, &multicallResults)
	if err != nil {
		return nil, err
	}

	result := map[string]*QuoteExactInputResult{}
	for i, path := range paths {
		var output QuoteExactInputResult
		multicallResult := multicallResults[i]
		if !multicallResult.Success {
			_, errMsg := lib_common.DecodeRevert(multicallResult.ReturnData)
			output = QuoteExactInputResult{Error: errors.New(errMsg)}
		} else {
			err = quoterABI.UnpackIntoInterface(&output, "quoteExactInput", multicallResult.ReturnData)
			if err != nil {
				return nil, err
			}
		}
		result[common.Bytes2Hex(path)] = &output
	}

	return result, nil
}

// Returns the amount in required for a given exact output swap without executing the swap
func (quoter *Quoter) QuoteExactOutput(ctx context.Context, path []byte, amountOut *big.Int) (*QuoteExactOutputResult, error) {
	quoterAddress := common.HexToAddress(quoter.sonexConf.Quoter)

	abi, err := abi.JSON(strings.NewReader(uniswap_v3.IquoterV2MetaData.ABI))
	if err != nil {
		return nil, err
	}
	data, err := abi.Pack("quoteExactOutput", path, amountOut)
	if err != nil {
		return nil, err
	}
	msg := ethereum.CallMsg{
		To:   &quoterAddress,
		Data: data,
	}
	client, err := ethclient.DialContext(ctx, quoter.chainConf.GetNode())
	if err != nil {
		return nil, err
	}
	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, err
	}

	var output QuoteExactOutputResult
	err = abi.UnpackIntoInterface(&output, "quoteExactOutput", result)
	if err != nil {
		return nil, err
	}
	return &output, nil
}

// Batch version of [Quoter.QuoteExactOutput]
func (quoter *Quoter) BatchQuoteExactOutputs(ctx context.Context, paths [][]byte, amountIn *big.Int) (map[string]*QuoteExactOutputResult, error) {
	quoterAddress := common.HexToAddress(quoter.sonexConf.Quoter)

	quoterABI, err := abi.JSON(strings.NewReader(uniswap_v3.IquoterV2MetaData.ABI))
	if err != nil {
		return nil, err
	}
	multicallABI, err := abi.JSON(strings.NewReader(contract_common.MulticallABI))
	if err != nil {
		return nil, err
	}
	nodeCaller, err := node_caller.NewNodeCaller(quoter.chainConf)
	if err != nil {
		return nil, err
	}
	var multicallResults []contract_common.Multicall3Result

	multicallData := []contract_common.Multicall3Call3{}
	for _, path := range paths {
		data, err := quoterABI.Pack("quoteExactOutput", path, amountIn)
		if err != nil {
			return nil, err
		}
		multicallData = append(multicallData, contract_common.Multicall3Call3{
			Target:       quoterAddress,
			AllowFailure: true,
			CallData:     data,
		})
	}
	data, err := multicallABI.Pack("aggregate3", multicallData)
	if err != nil {
		return nil, err
	}
	err = nodeCaller.StaticCall(ctx, multicallABI, common.HexToAddress(quoter.chainConf.Multicall), "aggregate3", data, &multicallResults)
	if err != nil {
		return nil, err
	}

	result := map[string]*QuoteExactOutputResult{}
	for i, path := range paths {
		var output QuoteExactOutputResult
		multicallResult := multicallResults[i]
		if !multicallResult.Success {
			_, errMsg := lib_common.DecodeRevert(multicallResult.ReturnData)
			output = QuoteExactOutputResult{Error: errors.New(errMsg)}
		} else {
			err = quoterABI.UnpackIntoInterface(&output, "quoteExactOutput", multicallResult.ReturnData)
			if err != nil {
				return nil, err
			}
		}
		result[common.Bytes2Hex(path)] = &output
	}

	return result, nil
}
