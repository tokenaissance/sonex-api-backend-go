package cache

import (
	"context"
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/config"
	lib_common "sonex_api_backend_go/internal/lib/common"
	contract_common "sonex_api_backend_go/internal/lib/contract/common"
	"sonex_api_backend_go/internal/lib/contract/uniswap_v3"
	"sonex_api_backend_go/internal/lib/node_caller"
	"sonex_api_backend_go/internal/lib/quoter"
	"sonex_api_backend_go/internal/module/sonex"
	"strings"
	"sync"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

const DEFAULT_PRICE_CACHE_EXPIRE_SECONDS int64 = 300

type TokenCacher struct {
	chainConf               config.ChainConfig
	sonexConf               config.SonexConfig
	tokens                  map[string]Token
	tokensLock              sync.RWMutex
	prices                  map[string]map[string]PriceCache
	pricesLock              sync.RWMutex
	priceCacheExpireSeconds int64
}

type Token struct {
	Address      string
	Name         string
	Symbol       string
	Decimals     int64
	IsStablecoin bool
}

type Price struct {
	From                        string
	To                          string
	Pools                       []quoter.Pool
	Path                        []byte
	Value                       *big.Int
	SqrtPriceX96AfterList       []*big.Int
	InitializedTicksCrossedList []uint32
	GasEstimate                 *big.Int
	PriceImpact                 float64
	Err                         error
}

type PriceCache struct {
	Prices    []Price
	UpdatedAt time.Time
}

func NewTokenCacher(chainConf config.ChainConfig, sonexConf config.SonexConfig) *TokenCacher {
	return &TokenCacher{
		chainConf:               chainConf,
		sonexConf:               sonexConf,
		tokens:                  map[string]Token{},
		prices:                  map[string]map[string]PriceCache{},
		priceCacheExpireSeconds: DEFAULT_PRICE_CACHE_EXPIRE_SECONDS,
	}
}

func (cacher *TokenCacher) sortPaths(ctx context.Context, paths []quoter.Path, maxNumber int64) ([]quoter.Path, error) {
	poolMap := map[string]*struct {
		Liquidity *big.Int
		Slot0     *sonex.Slot0
	}{}
	linq.From(paths).ForEachT(func(path quoter.Path) {
		linq.From(path.Pools).ForEachT(func(pool quoter.Pool) {
			poolMap[pool.Address.Hex()] = nil
		})
	})
	usedPools := []common.Address{}
	for key := range poolMap {
		usedPools = append(usedPools, common.HexToAddress(key))
	}

	if cacher.chainConf.Multicall == "" || !common.IsHexAddress(cacher.chainConf.Multicall) {
		return nil, fmt.Errorf("%s", "Invalid multicall contract")
	}

	callDatas := []contract_common.Multicall3Call{}
	poolABI, err := abi.JSON(strings.NewReader(uniswap_v3.UniswapV3PoolABI))
	if err != nil {
		return nil, err
	}
	liquidityData, err := poolABI.Pack("liquidity")
	if err != nil {
		return nil, err
	}
	slot0Data, err := poolABI.Pack("slot0")
	if err != nil {
		return nil, err
	}
	for _, pool := range usedPools {
		callDatas = append(callDatas, contract_common.Multicall3Call{
			Target:   pool,
			CallData: liquidityData,
		})
		callDatas = append(callDatas, contract_common.Multicall3Call{
			Target:   pool,
			CallData: slot0Data,
		})
	}
	multicallABI, err := abi.JSON(strings.NewReader(contract_common.MulticallABI))
	if err != nil {
		return nil, err
	}
	data, err := multicallABI.Pack("aggregate", callDatas)
	if err != nil {
		return nil, err
	}
	nodeCaller, err := node_caller.NewNodeCaller(cacher.chainConf)
	if err != nil {
		return nil, err
	}
	var result struct {
		BlockNumber *big.Int `json:"blockNumber"`
		ReturnData  [][]byte `json:"returnData"`
	}
	err = nodeCaller.StaticCall(ctx, multicallABI, common.HexToAddress(cacher.chainConf.Multicall), "aggregate", data, &result)
	if err != nil {
		return nil, err
	}
	for i, poolAddr := range usedPools {
		var liquidity *big.Int
		err = poolABI.UnpackIntoInterface(&liquidity, "liquidity", result.ReturnData[i*2])
		if err != nil {
			return nil, err
		}
		var slot0 sonex.Slot0
		err = poolABI.UnpackIntoInterface(&slot0, "slot0", result.ReturnData[i*2+1])
		if err != nil {
			return nil, err
		}
		poolMap[poolAddr.Hex()] = &struct {
			Liquidity *big.Int
			Slot0     *sonex.Slot0
		}{
			Liquidity: liquidity,
			Slot0:     &slot0,
		}
	}
	linq.From(paths).SelectT(func(path quoter.Path) struct {
		Path      quoter.Path
		Liquidity *big.Int
	} {
		liquidity := big.NewInt(0)
		for _, pool := range path.Pools {
			info := poolMap[pool.Address.Hex()]
			liquidity = new(big.Int).Add(liquidity, info.Liquidity)
		}
		return struct {
			Path      quoter.Path
			Liquidity *big.Int
		}{Path: path, Liquidity: liquidity}
	}).OrderByT(func(item struct {
		Path      quoter.Path
		Liquidity *big.Int
	}) int {
		return len(item.Path.Pools)
	}).ThenByDescendingT(func(item struct {
		Path      quoter.Path
		Liquidity *big.Int
	}) lib_common.CompareableBigInt {
		return lib_common.CompareableBigInt{Value: item.Liquidity}
	}).Take(int(maxNumber)).SelectT(func(item struct {
		Path      quoter.Path
		Liquidity *big.Int
	}) quoter.Path {
		return item.Path
	}).ToSlice(&paths)

	return paths, nil
}

func (cacher *TokenCacher) Get(address string) (*Token, error) {
	if address == "" || common.HexToAddress(address).Hex() == (common.Address{}).Hex() {
		eth := cacher.chainConf.ETH
		if eth == nil {
			return nil, fmt.Errorf("%s", "Miss ETH configure")
		}
		return &Token{
			Address:      (common.Address{}).Hex(),
			Name:         eth.Name,
			Symbol:       eth.Symbol,
			Decimals:     eth.Decimals,
			IsStablecoin: false,
		}, nil
	}

	cacher.tokensLock.RLock()
	token, ok := cacher.tokens[address]
	cacher.tokensLock.RUnlock()
	if !ok {
		tokenAddress := common.HexToAddress(address)
		if cacher.chainConf.Multicall == "" || !common.IsHexAddress(cacher.chainConf.Multicall) {
			return nil, fmt.Errorf("%s", "Invalid multicall contract")
		}
		nodeCaller, err := node_caller.NewNodeCaller(cacher.chainConf)
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
		var multicallResult struct {
			BlockNumber *big.Int `json:"blockNumber"`
			ReturnData  [][]byte `json:"returnData"`
		}
		callDatas := []contract_common.Multicall3Call{}
		data, err := erc20ABI.Pack("name")
		if err != nil {
			return nil, err
		}
		callDatas = append(callDatas, contract_common.Multicall3Call{
			Target:   tokenAddress,
			CallData: data,
		})
		data, err = erc20ABI.Pack("symbol")
		if err != nil {
			return nil, err
		}
		callDatas = append(callDatas, contract_common.Multicall3Call{
			Target:   tokenAddress,
			CallData: data,
		})
		data, err = erc20ABI.Pack("decimals")
		if err != nil {
			return nil, err
		}
		callDatas = append(callDatas, contract_common.Multicall3Call{
			Target:   tokenAddress,
			CallData: data,
		})
		data, err = multicallABI.Pack("aggregate", callDatas)
		if err != nil {
			return nil, err
		}
		err = nodeCaller.StaticCall(context.Background(), multicallABI, common.HexToAddress(cacher.chainConf.Multicall), "aggregate", data, &multicallResult)
		if err != nil {
			return nil, err
		}
		var name string
		err = erc20ABI.UnpackIntoInterface(&name, "name", multicallResult.ReturnData[0])
		if err != nil {
			return nil, err
		}
		var symbol string
		err = erc20ABI.UnpackIntoInterface(&symbol, "symbol", multicallResult.ReturnData[1])
		if err != nil {
			return nil, err
		}
		var decimals uint8
		err = erc20ABI.UnpackIntoInterface(&decimals, "decimals", multicallResult.ReturnData[2])
		if err != nil {
			return nil, err
		}

		token = Token{
			Address:      address,
			Name:         name,
			Symbol:       symbol,
			Decimals:     int64(decimals),
			IsStablecoin: false,
		}
		cacher.tokensLock.Lock()
		cacher.tokens[address] = token
		cacher.tokensLock.Unlock()
	}
	return &token, nil
}

func (c *TokenCacher) Set(address string, isStablecoin bool) error {
	token, err := c.Get(address)
	if err != nil {
		return err
	}
	token.IsStablecoin = isStablecoin
	c.tokensLock.Lock()
	c.tokens[address] = *token
	c.tokensLock.Unlock()
	return nil
}

func (cacher *TokenCacher) GetStablecoins() []Token {
	cacher.tokensLock.RLock()
	tokens := cacher.tokens
	cacher.tokensLock.RUnlock()
	var result []Token
	linq.From(tokens).
		WhereT(func(item linq.KeyValue) bool {
			token, _ := item.Value.(Token)
			return token.IsStablecoin
		}).
		SelectT(func(item linq.KeyValue) Token { return item.Value.(Token) }).
		ToSlice(&result)
	return result
}

func (cacher *TokenCacher) GetAll() []Token {
	cacher.tokensLock.RLock()
	tokens := cacher.tokens
	cacher.tokensLock.RUnlock()
	var result []Token
	linq.From(tokens).
		SelectT(func(item linq.KeyValue) Token { return item.Value.(Token) }).
		ToSlice(&result)
	return result
}

func (cacher *TokenCacher) SetPriceCacheExpireSeconds(duration int64) {
	if duration < 0 {
		duration = 300 // 5 minutes
	}
	cacher.priceCacheExpireSeconds = duration
}

func (cacher *TokenCacher) QuoteExactInput(ctx context.Context, pools []quoter.Pool, from common.Address, to common.Address, amountIn *big.Int, maxLength int64, maxNumber int64) ([]Price, error) {
	fromToken, err := cacher.Get(from.Hex())
	if err != nil {
		return nil, err
	}
	toToken, err := cacher.Get(to.Hex())
	if err != nil {
		return nil, err
	}

	q, err := quoter.NewQuoter(cacher.chainConf, cacher.sonexConf)
	if err != nil {
		return nil, err
	}
	paths := quoter.FindPaths(pools, common.HexToAddress(from.Hex()), common.HexToAddress(to.Hex()), maxLength)
	if len(paths) < 1 {
		return nil, fmt.Errorf("%s", "No available path")
	}
	// Filter paths
	paths, err = cacher.sortPaths(ctx, paths, maxNumber)
	if err != nil {
		return nil, err
	}

	pathBytes := [][]byte{}
	linq.From(paths).SelectT(func(item quoter.Path) []byte { return item.ToBytes() }).ToSlice(&pathBytes)
	results, err := q.BatchQuoteExactInputs(ctx, pathBytes, amountIn)
	if err != nil {
		return nil, err
	}

	// Handle results
	var prices []Price
	for _, path := range paths {
		result, ok := results[common.Bytes2Hex(path.ToBytes())]
		if !ok {
			continue
		}

		price := Price{
			From:  fromToken.Address,
			To:    toToken.Address,
			Pools: path.Pools,
			Path:  path.ToBytes(),
		}
		if result.Error != nil {
			price.Err = result.Error
		} else {
			price.Value = result.AmountOut
			price.SqrtPriceX96AfterList = result.SqrtPriceX96AfterList
			price.InitializedTicksCrossedList = result.InitializedTicksCrossedList
			price.GasEstimate = result.GasEstimate
		}

		prices = append(prices, price)
	}

	return prices, nil
}

func (cacher *TokenCacher) QuoteExactOutput(ctx context.Context, pools []quoter.Pool, from common.Address, to common.Address, amountOut *big.Int, maxLength int64, maxNumber int64) ([]Price, error) {
	fromToken, err := cacher.Get(from.Hex())
	if err != nil {
		return nil, err
	}
	toToken, err := cacher.Get(to.Hex())
	if err != nil {
		return nil, err
	}

	q, err := quoter.NewQuoter(cacher.chainConf, cacher.sonexConf)
	if err != nil {
		return nil, err
	}
	paths := quoter.FindPaths(pools, common.HexToAddress(from.Hex()), common.HexToAddress(to.Hex()), maxLength)
	if len(paths) < 1 {
		return nil, fmt.Errorf("%s", "No available path")
	}
	// Filter paths
	paths, err = cacher.sortPaths(ctx, paths, maxNumber)
	if err != nil {
		return nil, err
	}

	pathBytes := [][]byte{}
	linq.From(paths).SelectT(func(item quoter.Path) []byte { return item.ToBytes() }).ToSlice(&pathBytes)
	results, err := q.BatchQuoteExactOutputs(ctx, pathBytes, amountOut)
	if err != nil {
		return nil, err
	}

	// Handle results
	var prices []Price
	for i, path := range paths {
		prices = append(prices, Price{
			From:                        fromToken.Address,
			To:                          toToken.Address,
			Pools:                       path.Pools,
			Path:                        path.ToBytes(),
			Value:                       results[i].AmountIn,
			SqrtPriceX96AfterList:       results[i].SqrtPriceX96AfterList,
			InitializedTicksCrossedList: results[i].InitializedTicksCrossedList,
			GasEstimate:                 results[i].GasEstimate,
		})
	}

	return prices, nil
}

func (cacher *TokenCacher) GetPrices(ctx context.Context, pools []quoter.Pool, token common.Address, isCache bool, maxLength int64, maxNumber int64, wantStablecoin *common.Address) ([]Price, error) {
	if token.Hex() == (common.Address{}).Hex() {
		if !common.IsHexAddress(cacher.chainConf.WETH) {
			return nil, fmt.Errorf("%s", "Miss WETH configure")
		}
		token = common.HexToAddress(cacher.chainConf.WETH)
	}
	stablecoins := cacher.GetStablecoins()
	if len(stablecoins) == 0 {
		return nil, fmt.Errorf("%s", "No any stablecoins")
	}
	if wantStablecoin != nil {
		linq.From(stablecoins).
			WhereT(func(item Token) bool { return item.Address == wantStablecoin.Hex() }).
			ToSlice(&stablecoins)
	}
	tokenI, err := cacher.Get(token.Hex())
	if err != nil {
		return nil, err
	}

	amountIn := new(big.Int).Mul(big.NewInt(1), new(big.Int).Exp(big.NewInt(10), big.NewInt(tokenI.Decimals), nil))

	var wg sync.WaitGroup
	results := make(chan struct {
		Err    error
		Prices []Price
	}, len(stablecoins))
	sumPrices := []Price{}

	for _, stablecoin := range stablecoins {
		wg.Add(1)
		go func(ch chan struct {
			Err    error
			Prices []Price
		}) {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					ch <- struct {
						Err    error
						Prices []Price
					}{
						Err: fmt.Errorf("%s", err),
					}
				}
			}()

			var prices []Price

			if token.Hex() == common.HexToAddress(stablecoin.Address).Hex() {
				prices = []Price{{
					From:                        token.Hex(),
					To:                          stablecoin.Address,
					Pools:                       []quoter.Pool{},
					Path:                        []byte{},
					Value:                       new(big.Int).Mul(big.NewInt(1), new(big.Int).Exp(big.NewInt(10), big.NewInt(stablecoin.Decimals), nil)),
					SqrtPriceX96AfterList:       []*big.Int{},
					InitializedTicksCrossedList: []uint32{},
					GasEstimate:                 big.NewInt(0),
					PriceImpact:                 0,
				}}
			}

			if prices != nil && isCache {
				cacher.pricesLock.RLock()
				from, ok := cacher.prices[tokenI.Address]
				if ok {
					priceCache, ok := from[stablecoin.Address]
					if ok && priceCache.UpdatedAt.Add(time.Second*time.Duration(cacher.priceCacheExpireSeconds)).Unix() > time.Now().UTC().Unix() {
						prices = priceCache.Prices
					}
				}
				cacher.pricesLock.RUnlock()
			}

			if prices != nil {
				linq.From(prices).
					WhereT(func(price Price) bool {
						for _, pathPool := range price.Pools {
							contains := func(ps []quoter.Pool, p string) bool {
								for _, pool := range ps {
									if pool.Address.Hex() == p {
										return true
									}
								}
								return false
							}
							if !contains(pools, pathPool.Address.Hex()) {
								return false
							}
						}
						return true
					}).
					ToSlice(&prices)
			} else {
				prices, err = cacher.QuoteExactInput(ctx, pools, token, common.HexToAddress(stablecoin.Address), amountIn, maxLength, maxNumber)
				if err != nil {
					panic(err)
				}

				// Add cache
				updatedAt := time.Now().UTC()
				cacher.pricesLock.Lock()
				cacher.prices[token.Hex()] = map[string]PriceCache{}
				cacher.prices[token.Hex()][stablecoin.Address] = PriceCache{
					Prices:    prices,
					UpdatedAt: updatedAt,
				}
				cacher.pricesLock.Unlock()
			}

			ch <- struct {
				Err    error
				Prices []Price
			}{
				Prices: prices,
			}
		}(results)
	}

	wg.Wait()
	close(results)

	// Handle results
	for result := range results {
		if result.Err == nil {
			sumPrices = append(sumPrices, result.Prices...)
		}
	}

	return sumPrices, nil
}

func (cacher *TokenCacher) BatchGetTotalSupplies(ctx context.Context, tokens []common.Address) (map[common.Address]*big.Int, error) {
	nodeCaller, err := node_caller.NewNodeCaller(cacher.chainConf)
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
	var multicallResults []contract_common.Multicall3Result
	callDatas := []contract_common.Multicall3Call3{}
	data, err := erc20ABI.Pack("totalSupply")
	if err != nil {
		return nil, err
	}
	for _, token := range tokens {
		callDatas = append(callDatas, contract_common.Multicall3Call3{
			Target:       token,
			AllowFailure: true,
			CallData:     data,
		})
	}
	multicallData, err := multicallABI.Pack("aggregate3", callDatas)
	if err != nil {
		return nil, err
	}
	err = nodeCaller.StaticCall(ctx, multicallABI, common.HexToAddress(cacher.chainConf.Multicall), "aggregate3", multicallData, &multicallResults)
	if err != nil {
		return nil, err
	}
	result := map[common.Address]*big.Int{}
	for i, token := range tokens {
		multicallResult := multicallResults[i]
		if !multicallResult.Success {
			continue
		}
		var totalSupply *big.Int
		err = erc20ABI.UnpackIntoInterface(&totalSupply, "totalSupply", multicallResult.ReturnData)
		if err != nil {
			return nil, err
		}
		result[token] = totalSupply
	}
	return result, nil
}
