package quoter

import (
	"context"
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/contract/uniswap_v3"
	"strings"

	"github.com/ahmetb/go-linq/v3"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Quoter struct {
	chain config.ChainConfig
}

func NewQuoter(chain config.ChainConfig) (*Quoter, error) {
	if chain.Quoter == "" || !common.IsHexAddress(chain.Quoter) {
		return nil, fmt.Errorf("%s", "Invalid quoter address")
	}

	return &Quoter{
		chain: chain,
	}, nil
}

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

				if pool.Token0.Hex() == currentToken.Hex() {
					nextToken = pool.Token1
					fee = pool.Fee
				} else {
					nextToken = pool.Token0
					fee = pool.Fee
				}

				if !visited[nextToken.Hex()] {
					if !contains(path.Tokens, common.HexToAddress(pool.Token1.Hex())) {
						path.Pools = append(path.Pools, pool)
						path.Tokens = append(path.Tokens, common.HexToAddress(pool.Token1.Hex()))
						path.Fees = append(path.Fees, fee)

						if maxLength != 0 && len(path.Pools) > int(maxLength) {
							return
						}
					}

					dfs(nextToken, path)

					path.Pools = path.Pools[:len(path.Pools)-1]
					path.Tokens = path.Tokens[:len(path.Tokens)-1]
					path.Fees = path.Fees[:len(path.Fees)-1]
				}
			}
		}

		visited[currentToken.Hex()] = false
	}
	dfs(from, Path{Tokens: []common.Address{common.HexToAddress(from.Hex())}, Fees: []int64{}})

	linq.From(paths).OrderByT(func(item Path) int { return len(item.Tokens) }).ToSlice(&paths)
	return paths
}

func (quoter *Quoter) QuoteExactInput(ctx context.Context, path []byte, amountIn *big.Int) (*QuoteExactInputResult, error) {
	quoterAddress := common.HexToAddress(quoter.chain.Quoter)

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
	client, err := ethclient.DialContext(ctx, quoter.chain.GetNode())
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

func (quoter *Quoter) QuoteExactOutput(ctx context.Context, path []byte, amountOut *big.Int) (*QuoteExactOutputResult, error) {
	quoterAddress := common.HexToAddress(quoter.chain.Quoter)

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
	client, err := ethclient.DialContext(ctx, quoter.chain.GetNode())
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
