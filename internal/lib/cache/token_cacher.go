package cache

import (
	"context"
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/config"
	contract_common "sonex_api_backend_go/internal/lib/contract/common"
	"sonex_api_backend_go/internal/lib/quoter"
	"sync"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const DEFAULT_PRICE_CACHE_EXPIRE_SECONDS int64 = 300

type TokenCacher struct {
	chain                   config.ChainConfig
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
}

type PriceCache struct {
	Prices    []Price
	UpdatedAt time.Time
}

func NewTokenCacher(chain config.ChainConfig) *TokenCacher {
	return &TokenCacher{
		chain:                   chain,
		tokens:                  map[string]Token{},
		prices:                  map[string]map[string]PriceCache{},
		priceCacheExpireSeconds: DEFAULT_PRICE_CACHE_EXPIRE_SECONDS,
	}
}

func (cacher *TokenCacher) Get(address string) (Token, error) {
	cacher.tokensLock.RLock()
	token, ok := cacher.tokens[address]
	cacher.tokensLock.RUnlock()
	if !ok {
		client, err := ethclient.Dial(cacher.chain.GetNode())
		if err != nil {
			return Token{}, err
		}
		erc20, err := contract_common.NewErc20(common.HexToAddress(address), client)
		if err != nil {
			return Token{}, err
		}
		name, err := erc20.Name(nil)
		if err != nil {
			return Token{}, err
		}
		symbol, err := erc20.Symbol(nil)
		if err != nil {
			return Token{}, err
		}
		decimals, err := erc20.Decimals(nil)
		if err != nil {
			return Token{}, err
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
	return token, nil
}

func (c *TokenCacher) Set(address string, isStablecoin bool) error {
	token, err := c.Get(address)
	if err != nil {
		return err
	}
	token.IsStablecoin = isStablecoin
	c.tokensLock.Lock()
	c.tokens[address] = token
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

func (cacher *TokenCacher) QuoteExactInput(pools []quoter.Pool, from common.Address, to common.Address, amountIn *big.Int, maxLength int64) ([]Price, error) {
	fromToken, err := cacher.Get(from.Hex())
	if err != nil {
		return nil, err
	}
	toToken, err := cacher.Get(to.Hex())
	if err != nil {
		return nil, err
	}

	q, err := quoter.NewQuoter(cacher.chain)
	if err != nil {
		return nil, err
	}
	paths := quoter.FindPaths(pools, common.HexToAddress(from.Hex()), common.HexToAddress(to.Hex()), maxLength)
	if len(paths) < 1 {
		return nil, fmt.Errorf("%s", "No available path")
	}

	var wg sync.WaitGroup
	results := make(chan struct {
		Path  quoter.Path
		Err   error
		Price quoter.QuoteExactInputResult
	}, len(paths))

	for _, path := range paths {
		wg.Add(1)

		go func(ch chan struct {
			Path  quoter.Path
			Err   error
			Price quoter.QuoteExactInputResult
		}) {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					ch <- struct {
						Path  quoter.Path
						Err   error
						Price quoter.QuoteExactInputResult
					}{
						Path: path,
						Err:  fmt.Errorf("%s", err),
					}
				}
			}()

			result, err := q.QuoteExactInput(context.Background(), path.ToBytes(), amountIn)
			if err != nil {
				panic(err)
			}

			ch <- struct {
				Path  quoter.Path
				Err   error
				Price quoter.QuoteExactInputResult
			}{
				Path:  path,
				Price: *result,
			}
		}(results)
	}

	wg.Wait()
	close(results)

	// Handle results
	var prices []Price
	linq.From(results).
		WhereT(func(item struct {
			Path  quoter.Path
			Err   error
			Price quoter.QuoteExactInputResult
		}) bool {
			return item.Err == nil
		}).
		SelectT(func(item struct {
			Path  quoter.Path
			Err   error
			Price quoter.QuoteExactInputResult
		}) Price {
			return Price{
				From:                        fromToken.Address,
				To:                          toToken.Address,
				Pools:                       item.Path.Pools,
				Path:                        item.Path.ToBytes(),
				Value:                       item.Price.AmountOut,
				SqrtPriceX96AfterList:       item.Price.SqrtPriceX96AfterList,
				InitializedTicksCrossedList: item.Price.InitializedTicksCrossedList,
				GasEstimate:                 item.Price.GasEstimate,
			}
		}).
		ToSlice(&prices)

	return prices, nil
}

func (cacher *TokenCacher) QuoteExactOutput(pools []quoter.Pool, from common.Address, to common.Address, amountOut *big.Int, maxLength int64) ([]Price, error) {
	fromToken, err := cacher.Get(from.Hex())
	if err != nil {
		return nil, err
	}
	toToken, err := cacher.Get(to.Hex())
	if err != nil {
		return nil, err
	}

	q, err := quoter.NewQuoter(cacher.chain)
	if err != nil {
		return nil, err
	}
	paths := quoter.FindPaths(pools, common.HexToAddress(from.Hex()), common.HexToAddress(to.Hex()), maxLength)
	if len(paths) < 1 {
		return nil, fmt.Errorf("%s", "No available path")
	}

	var wg sync.WaitGroup
	results := make(chan struct {
		Path  quoter.Path
		Err   error
		Price quoter.QuoteExactOutputResult
	}, len(paths))

	for _, path := range paths {
		wg.Add(1)

		go func(ch chan struct {
			Path  quoter.Path
			Err   error
			Price quoter.QuoteExactOutputResult
		}) {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					ch <- struct {
						Path  quoter.Path
						Err   error
						Price quoter.QuoteExactOutputResult
					}{
						Path: path,
						Err:  fmt.Errorf("%s", err),
					}
				}
			}()

			result, err := q.QuoteExactOutput(context.Background(), path.ToBytes(), amountOut)
			if err != nil {
				panic(err)
			}

			ch <- struct {
				Path  quoter.Path
				Err   error
				Price quoter.QuoteExactOutputResult
			}{
				Path:  path,
				Price: *result,
			}
		}(results)
	}

	wg.Wait()
	close(results)

	// Handle results
	var prices []Price
	linq.From(results).
		WhereT(func(item struct {
			Path  quoter.Path
			Err   error
			Price quoter.QuoteExactOutputResult
		}) bool {
			return item.Err == nil
		}).
		SelectT(func(item struct {
			Path  quoter.Path
			Err   error
			Price quoter.QuoteExactOutputResult
		}) Price {
			return Price{
				From:                        fromToken.Address,
				To:                          toToken.Address,
				Pools:                       item.Path.Pools,
				Path:                        item.Path.ToBytes(),
				Value:                       item.Price.AmountIn,
				SqrtPriceX96AfterList:       item.Price.SqrtPriceX96AfterList,
				InitializedTicksCrossedList: item.Price.InitializedTicksCrossedList,
				GasEstimate:                 item.Price.GasEstimate,
			}
		}).
		ToSlice(&prices)

	return prices, nil
}

func (cacher *TokenCacher) GetPrices(pools []quoter.Pool, token common.Address, isCache bool, maxLength int64, wantStablecoin *common.Address) ([]Price, error) {
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

			if isCache {
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
				prices, err = cacher.QuoteExactInput(pools, token, common.HexToAddress(stablecoin.Address), amountIn, maxLength)
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
