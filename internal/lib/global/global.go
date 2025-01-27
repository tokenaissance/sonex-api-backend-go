package global

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"sonex_api_backend_go/internal/chain_handler"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/cache"
	"sonex_api_backend_go/internal/lib/chain_scanner"
	contract_common "sonex_api_backend_go/internal/lib/contract/common"
	"sonex_api_backend_go/internal/lib/contract/uniswap_v3"
	"sonex_api_backend_go/internal/lib/id_generater"
	"sonex_api_backend_go/internal/lib/node_caller"
	"sonex_api_backend_go/internal/lib/quoter"
	"sonex_api_backend_go/internal/module/sonex"
	"sonex_api_backend_go/internal/repository"
	"strings"
	"sync"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/robfig/cron/v3"
)

var c config.Config
var chainHandlerManager *chain_handler.ChainHandlerManager
var cronExecuter *cron.Cron
var idGenerater *id_generater.IdGenerater
var tokenCachers map[int64]*cache.TokenCacher
var tokenCachersLock sync.RWMutex

// key : fmt.Sprintf("volume24H_%v", _keyPoolAddress)
// value: volume24H ,string
// expiration : 24H
// volume24H
var Volume24HKey = "Volume24H"
var Tvl24HKey = "Tvl24HKey"
var YiedlRateKey = "YiedlRateKey"

func GetVolume24HTvlValue(cacher cache.ICache, _key string, _keyPoolAddress string) (string, error) {
	return cacher.Get(fmt.Sprintf("%v_%v", _key, _keyPoolAddress))
}

func SetVolume24HTvlValue(cacher cache.ICache, _key string, _keyPoolAddress string, _valueVolume24H string) error {
	return cacher.Set(fmt.Sprintf("%v_%v", _key, _keyPoolAddress), _valueVolume24H, 24*time.Hour)
}

// Initialize application's components
func Initialize(config config.Config) error {
	c = config

	// Initialize chain handler manager
	manager, err := chain_handler.NewChainHandlerManager(c)
	if err != nil {
		panic(err)
	}
	chainHandlerManager = manager

	// Initialize cron executer
	cronExecuter = cron.New(cron.WithSeconds(), cron.WithLocation(time.UTC))

	// Initialize token cacher
	tokenCachers = map[int64]*cache.TokenCacher{}
	for _, chain := range c.Chains {
		sonexConf := c.GetSonexConfig(chain.ChainId)
		if sonexConf == nil {
			return fmt.Errorf("%s", "Miss sonex configure")
		}

		tokenCacher := cache.NewTokenCacher(chain, *sonexConf)
		tokenCachersLock.Lock()
		tokenCachers[chain.ChainId] = tokenCacher
		tokenCachersLock.Unlock()
	}

	return nil
}

// Adds a scanning handler to specified blockchain scanner
func AddChainHandler(chainId int64, handler chain_scanner.IHandler) {
	chainHandlerManager.AddHandler(chainId, handler)
}

// Register a replay (rescan) task
func Replay(chainId int64, request chain_scanner.Replay) {
	chainHandlerManager.Replay(chainId, request)
}

// Add a scheduled job to Cron
func AddCronJob(cron string, job func()) (cron.EntryID, error) {
	return cronExecuter.AddFunc(cron, job)
}

// Starts application's components
func Start() {
	chainHandlerManager.Start()
	cronExecuter.Start()
}

// Stop application's components
func Stop() {
	chainHandlerManager.Stop()
	cronExecuter.Stop()
}

// Init Unique ID generater
func InitIdGenerater(workerId int64) error {
	generater, err := id_generater.NewIdGenerater(workerId)
	if err != nil {
		return err
	}
	idGenerater = generater
	return nil
}

// Returns a new unique ID
func NewId() (int64, error) {
	if idGenerater == nil {
		return 0, fmt.Errorf("%s", "Id generater not initialized")
	}
	return idGenerater.Generate()
}

// Get token's metadata from cache
func GetToken(chainId int64, address string) (*cache.Token, error) {
	tokenCachersLock.RLock()
	cacher, ok := tokenCachers[chainId]
	tokenCachersLock.RUnlock()
	if !ok {
		return nil, fmt.Errorf("%s", "Miss token cacher")
	}
	return cacher.Get(common.HexToAddress(address))
}

// Cache token's metadata
func SetToken(chainId int64, address string) error {
	tokenCachersLock.RLock()
	cacher, ok := tokenCachers[chainId]
	tokenCachersLock.RUnlock()
	if !ok {
		return fmt.Errorf("%s", "Miss token cacher")
	}
	err := cacher.Set(common.HexToAddress(address))
	if err != nil {
		return err
	}
	return nil
}

// Calculates `amount / (10 ** decimals)`, assumes that it's USD token on `chainId`
func PriceToUSD(chainId int64, price *big.Int) (float64, error) {
	chainConf := c.GetChainConfig(chainId)
	if chainConf == nil {
		return 0, fmt.Errorf("%s", "Miss chain configure")
	}
	usdToken, err := GetToken(chainId, common.HexToAddress(chainConf.USD).Hex())
	if err != nil {
		return 0, err
	}
	usdAccuracy := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(usdToken.Decimals), nil))
	val, _ := new(big.Float).Quo(new(big.Float).SetInt(price), usdAccuracy).SetPrec(math.MaxUint32).Float64()
	return val, nil
}

// Get prices in USD of `token`,
// and sort the prices by defined strategy.
func GetPrices(ctx context.Context, pools []quoter.Pool, chainId int64, token string, isCache bool, maxLength int64) (*cache.Price, []cache.Price, []cache.Price, error) {
	tokenCachersLock.RLock()
	cacher, ok := tokenCachers[chainId]
	tokenCachersLock.RUnlock()
	if !ok {
		return nil, nil, nil, fmt.Errorf("%s", "Miss token cacher")
	}

	prices, err := cacher.GetPrices(ctx, pools, common.HexToAddress(token), isCache, maxLength, 5)
	if err != nil {
		return nil, nil, nil, err
	}
	successedPrices := []cache.Price{}
	failedPrices := []cache.Price{}
	for _, price := range prices {
		if price.Err != nil {
			failedPrices = append(failedPrices, price)
		} else {
			successedPrices = append(successedPrices, price)
		}
	}

	// Best price
	var bestPrice *cache.Price
	if len(successedPrices) > 0 {
		// Sort by path length
		linq.From(successedPrices).OrderByT(func(item cache.Price) int { return len(item.Pools) }).ToSlice(&prices)
		pathLength := len(successedPrices[0].Pools)
		linq.From(successedPrices).WhereT(func(item cache.Price) bool {
			return len(item.Pools) == pathLength
		}).ToSlice(&successedPrices)
		poolMap := map[string]*big.Int{}
		linq.From(successedPrices).ForEachT(func(price cache.Price) {
			linq.From(price.Pools).ForEachT(func(pool quoter.Pool) {
				poolMap[pool.Address.Hex()] = common.Big0
			})
		})
		usedPools := []common.Address{}
		for key := range poolMap {
			usedPools = append(usedPools, common.HexToAddress(key))
		}

		chain := c.GetChainConfig(chainId)
		if chain == nil {
			return nil, nil, nil, fmt.Errorf("%s", "Miss chain configure")
		}
		if chain.Multicall == "" || !common.IsHexAddress(chain.Multicall) {
			return nil, nil, nil, fmt.Errorf("%s", "Invalid multicall contract")
		}

		callDatas := []contract_common.Multicall3Call{}
		poolABI, err := abi.JSON(strings.NewReader(uniswap_v3.UniswapV3PoolABI))
		if err != nil {
			return nil, nil, nil, err
		}
		data, err := poolABI.Pack("liquidity")
		if err != nil {
			return nil, nil, nil, err
		}
		for _, pool := range usedPools {
			callDatas = append(callDatas, contract_common.Multicall3Call{
				Target:   pool,
				CallData: data,
			})
		}
		multicallABI, err := abi.JSON(strings.NewReader(contract_common.MulticallABI))
		if err != nil {
			return nil, nil, nil, err
		}
		data, err = multicallABI.Pack("aggregate", callDatas)
		if err != nil {
			return nil, nil, nil, err
		}
		nodeCaller, err := node_caller.NewNodeCaller(*chain)
		if err != nil {
			return nil, nil, nil, err
		}
		var result struct {
			BlockNumber *big.Int `json:"blockNumber"`
			ReturnData  [][]byte `json:"returnData"`
		}
		err = nodeCaller.StaticCall(ctx, multicallABI, common.HexToAddress(chain.Multicall), "aggregate", data, &result)
		if err != nil {
			return nil, nil, nil, err
		}
		for i, value := range result.ReturnData {
			poolMap[usedPools[i].Hex()] = new(big.Int).SetBytes(value)
		}
		maxLiquidity := common.Big0
		for _, price := range successedPrices {
			liquidity := common.Big0
			for _, pool := range price.Pools {
				liquidity = new(big.Int).Add(liquidity, poolMap[pool.Address.Hex()])
			}
			if maxLiquidity.Cmp(liquidity) <= 0 {
				maxLiquidity = liquidity
				bestPrice = &price
			}
		}
	}

	return bestPrice, successedPrices, failedPrices, nil
}

// Wrap [*cache.TokenCacher.QuoteExactInput] and sort the result by defined strategy
func QuoteExactInput(ctx context.Context, pools []quoter.Pool, chainId int64, from string, to string, amountIn *big.Int, maxLength int64) ([]cache.Price, []cache.Price, error) {
	chainConf := c.GetChainConfig(chainId)
	if chainConf == nil {
		return nil, nil, fmt.Errorf("%s", "Miss chain configure")
	}

	fromAddress := common.HexToAddress(from)
	toAddress := common.HexToAddress(to)

	tokenCachersLock.RLock()
	cacher, ok := tokenCachers[chainId]
	tokenCachersLock.RUnlock()
	if !ok {
		return nil, nil, fmt.Errorf("%s", "Miss token cacher")
	}

	if pools == nil {
		conn := repository.NewConn(c.PostgreSQL.DateSource())
		eventSonexFactoryPoolCreatedModel := repository.NewEventSonexFactoryPoolCreatedModel(conn)
		records, err := eventSonexFactoryPoolCreatedModel.QueryByChainId(ctx, chainId)
		if err != nil {
			return nil, nil, err
		}
		linq.From(records).
			SelectT(func(item repository.EventSonexFactoryPoolCreated) quoter.Pool {
				return quoter.Pool{
					Address: common.HexToAddress(item.EventPool),
					Token0:  common.HexToAddress(item.EventToken0),
					Token1:  common.HexToAddress(item.EventToken1),
					Fee:     item.EventFee,
				}
			}).
			ToSlice(&pools)
	}

	prices, err := cacher.QuoteExactInput(ctx, pools, fromAddress, toAddress, amountIn, maxLength, 200)
	if err != nil {
		return nil, nil, err
	}
	// ETH or USD
	if len(prices) == 1 && len(prices[0].Pools) == 0 {
		return prices, []cache.Price{}, nil
	}

	successedPrices := []cache.Price{}
	failedPrices := []cache.Price{}
	for _, price := range prices {
		if price.Err != nil {
			failedPrices = append(failedPrices, price)
		} else {
			successedPrices = append(successedPrices, price)
		}
	}

	if len(successedPrices) > 0 {
		poolMap := map[string]*struct {
			Liquidity *big.Int
			Slot0     *sonex.Slot0
		}{}
		linq.From(successedPrices).ForEachT(func(price cache.Price) {
			linq.From(price.Pools).ForEachT(func(pool quoter.Pool) {
				poolMap[pool.Address.Hex()] = nil
			})
		})
		usedPools := []common.Address{}
		for key := range poolMap {
			usedPools = append(usedPools, common.HexToAddress(key))
		}
		chain := c.GetChainConfig(chainId)
		if chain == nil {
			return nil, nil, fmt.Errorf("%s", "Miss chain configure")
		}
		if chain.Multicall == "" || !common.IsHexAddress(chain.Multicall) {
			return nil, nil, fmt.Errorf("%s", "Invalid multicall contract")
		}

		callDatas := []contract_common.Multicall3Call{}
		poolABI, err := abi.JSON(strings.NewReader(uniswap_v3.UniswapV3PoolABI))
		if err != nil {
			return nil, nil, err
		}
		liquidityData, err := poolABI.Pack("liquidity")
		if err != nil {
			return nil, nil, err
		}
		slot0Data, err := poolABI.Pack("slot0")
		if err != nil {
			return nil, nil, err
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
			return nil, nil, err
		}
		data, err := multicallABI.Pack("aggregate", callDatas)
		if err != nil {
			return nil, nil, err
		}
		nodeCaller, err := node_caller.NewNodeCaller(*chain)
		if err != nil {
			return nil, nil, err
		}
		var result struct {
			BlockNumber *big.Int `json:"blockNumber"`
			ReturnData  [][]byte `json:"returnData"`
		}
		err = nodeCaller.StaticCall(ctx, multicallABI, common.HexToAddress(chain.Multicall), "aggregate", data, &result)
		if err != nil {
			return nil, nil, err
		}
		for i, poolAddr := range usedPools {
			var liquidity *big.Int
			err = poolABI.UnpackIntoInterface(&liquidity, "liquidity", result.ReturnData[i*2])
			if err != nil {
				return nil, nil, err
			}
			var slot0 sonex.Slot0
			err = poolABI.UnpackIntoInterface(&slot0, "slot0", result.ReturnData[i*2+1])
			if err != nil {
				return nil, nil, err
			}
			poolMap[poolAddr.Hex()] = &struct {
				Liquidity *big.Int
				Slot0     *sonex.Slot0
			}{
				Liquidity: liquidity,
				Slot0:     &slot0,
			}
		}

		// Filter price if has zero liquidity pool
		linq.From(successedPrices).WhereT(func(item cache.Price) bool {
			for _, pool := range item.Pools {
				info := poolMap[pool.Address.Hex()]
				if info.Liquidity.Cmp(big.NewInt(0)) <= 0 {
					return false
				}
			}
			return true
		}).ToSlice(&successedPrices)

		// Calculate price impact and fee rate
		var wg sync.WaitGroup
		results := make(chan struct {
			Index       int64
			PriceImpact float64
			FeeRate     float64
			Err         error
		}, len(successedPrices))
		for i, price := range successedPrices {
			wg.Add(1)
			go func(ch chan struct {
				Index       int64
				PriceImpact float64
				FeeRate     float64
				Err         error
			}) {
				defer wg.Done()
				defer func() {
					if err := recover(); err != nil {
						ch <- struct {
							Index       int64
							PriceImpact float64
							FeeRate     float64
							Err         error
						}{
							Index: int64(i),
							Err:   fmt.Errorf("%s(panic)", err),
						}
					}
				}()

				totalPriceImpact := big.NewFloat(0)
				totalAmount := big.NewFloat(0)
				stepAmountIn := amountIn
				var totalFeeRate float64 = 1
				for j, pool := range price.Pools {
					info := poolMap[pool.Address.Hex()]
					_, from, to, err := price.Path.GetPool(int64(j))
					if err != nil {
						ch <- struct {
							Index       int64
							PriceImpact float64
							FeeRate     float64
							Err         error
						}{
							Index: int64(i),
							Err:   fmt.Errorf("%s", err),
						}
						return
					}
					stepPrices, err := cacher.QuoteExactInput(ctx, []quoter.Pool{pool}, *from, *to, stepAmountIn, 0, 1)
					if err != nil {
						ch <- struct {
							Index       int64
							PriceImpact float64
							FeeRate     float64
							Err         error
						}{
							Index: int64(i),
							Err:   fmt.Errorf("%s", err),
						}
						return
					}
					if len(stepPrices) < 1 {
						ch <- struct {
							Index       int64
							PriceImpact float64
							FeeRate     float64
							Err         error
						}{
							Index: int64(i),
							Err:   fmt.Errorf("%s", "Can not quote step price"),
						}
						return
					}
					fromToken, err := GetToken(chainId, from.Hex())
					if err != nil {
						ch <- struct {
							Index       int64
							PriceImpact float64
							FeeRate     float64
							Err         error
						}{
							Index: int64(i),
							Err:   err,
						}
						return
					}
					fromAccuracy := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(fromToken.Decimals), nil))
					stepAmountInF := new(big.Float).Quo(new(big.Float).SetInt(stepAmountIn), fromAccuracy)
					priceImpact := sonex.CalculatePriceImpact(info.Slot0.SqrtPriceX96, price.SqrtPriceX96AfterList[j])
					totalPriceImpact = new(big.Float).Add(totalPriceImpact, new(big.Float).Mul(priceImpact, stepAmountInF))
					totalAmount = new(big.Float).Add(totalAmount, stepAmountInF)
					totalFeeRate = totalFeeRate * (1 - float64(pool.Fee)/1000000)
					stepAmountIn = stepPrices[0].Value
				}
				avgPriceImpact, _ := new(big.Float).Quo(totalPriceImpact, totalAmount).Float64()
				ch <- struct {
					Index       int64
					PriceImpact float64
					FeeRate     float64
					Err         error
				}{
					Index:       int64(i),
					PriceImpact: avgPriceImpact,
					FeeRate:     1 - totalFeeRate,
				}
			}(results)
		}
		wg.Wait()
		close(results)
		for result := range results {
			if result.Err != nil {
				return nil, nil, result.Err
			}
			successedPrices[result.Index].PriceImpact = result.PriceImpact
			successedPrices[result.Index].FeeRate = result.FeeRate
		}
	}

	return successedPrices, failedPrices, nil
}

// Wrap [*cache.TokenCacher.QuoteExactOutput] and sort the result by defined strategy
func QuoteExactOutput(ctx context.Context, pools []quoter.Pool, chainId int64, from string, to string, amountOut *big.Int, maxLength int64) ([]cache.Price, []cache.Price, error) {
	chainConf := c.GetChainConfig(chainId)
	if chainConf == nil {
		return nil, nil, fmt.Errorf("%s", "Miss chain configure")
	}

	fromAddress := common.HexToAddress(from)
	toAddress := common.HexToAddress(to)

	tokenCachersLock.RLock()
	cacher, ok := tokenCachers[chainId]
	tokenCachersLock.RUnlock()
	if !ok {
		return nil, nil, fmt.Errorf("%s", "Miss token cacher")
	}

	if pools == nil {
		conn := repository.NewConn(c.PostgreSQL.DateSource())
		eventSonexFactoryPoolCreatedModel := repository.NewEventSonexFactoryPoolCreatedModel(conn)
		records, err := eventSonexFactoryPoolCreatedModel.QueryByChainId(ctx, chainId)
		if err != nil {
			return nil, nil, err
		}
		linq.From(records).
			SelectT(func(item repository.EventSonexFactoryPoolCreated) quoter.Pool {
				return quoter.Pool{
					Address: common.HexToAddress(item.EventPool),
					Token0:  common.HexToAddress(item.EventToken0),
					Token1:  common.HexToAddress(item.EventToken1),
					Fee:     item.EventFee,
				}
			}).
			ToSlice(&pools)
	}

	prices, err := cacher.QuoteExactOutput(ctx, pools, fromAddress, toAddress, amountOut, maxLength, 200)
	if err != nil {
		return nil, nil, err
	}
	// ETH or USD
	if len(prices) == 1 && len(prices[0].Pools) == 0 {
		return prices, []cache.Price{}, nil
	}

	successedPrices := []cache.Price{}
	failedPrices := []cache.Price{}
	for _, price := range prices {
		if price.Err != nil {
			failedPrices = append(failedPrices, price)
		} else {
			successedPrices = append(successedPrices, price)
		}
	}

	if len(successedPrices) > 0 {
		poolMap := map[string]*struct {
			Liquidity *big.Int
			Slot0     *sonex.Slot0
		}{}
		linq.From(successedPrices).ForEachT(func(price cache.Price) {
			linq.From(price.Pools).ForEachT(func(pool quoter.Pool) {
				poolMap[pool.Address.Hex()] = nil
			})
		})
		usedPools := []common.Address{}
		for key := range poolMap {
			usedPools = append(usedPools, common.HexToAddress(key))
		}
		chain := c.GetChainConfig(chainId)
		if chain == nil {
			return nil, nil, fmt.Errorf("%s", "Miss chain configure")
		}
		if chain.Multicall == "" || !common.IsHexAddress(chain.Multicall) {
			return nil, nil, fmt.Errorf("%s", "Invalid multicall contract")
		}

		callDatas := []contract_common.Multicall3Call{}
		poolABI, err := abi.JSON(strings.NewReader(uniswap_v3.UniswapV3PoolABI))
		if err != nil {
			return nil, nil, err
		}
		liquidityData, err := poolABI.Pack("liquidity")
		if err != nil {
			return nil, nil, err
		}
		slot0Data, err := poolABI.Pack("slot0")
		if err != nil {
			return nil, nil, err
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
			return nil, nil, err
		}
		data, err := multicallABI.Pack("aggregate", callDatas)
		if err != nil {
			return nil, nil, err
		}
		nodeCaller, err := node_caller.NewNodeCaller(*chain)
		if err != nil {
			return nil, nil, err
		}
		var result struct {
			BlockNumber *big.Int `json:"blockNumber"`
			ReturnData  [][]byte `json:"returnData"`
		}
		err = nodeCaller.StaticCall(ctx, multicallABI, common.HexToAddress(chain.Multicall), "aggregate", data, &result)
		if err != nil {
			return nil, nil, err
		}
		for i, poolAddr := range usedPools {
			var liquidity *big.Int
			err = poolABI.UnpackIntoInterface(&liquidity, "liquidity", result.ReturnData[i*2])
			if err != nil {
				return nil, nil, err
			}
			var slot0 sonex.Slot0
			err = poolABI.UnpackIntoInterface(&slot0, "slot0", result.ReturnData[i*2+1])
			if err != nil {
				return nil, nil, err
			}
			poolMap[poolAddr.Hex()] = &struct {
				Liquidity *big.Int
				Slot0     *sonex.Slot0
			}{
				Liquidity: liquidity,
				Slot0:     &slot0,
			}
		}

		// Filter price if has zero liquidity pool
		linq.From(successedPrices).WhereT(func(item cache.Price) bool {
			for _, pool := range item.Pools {
				info := poolMap[pool.Address.Hex()]
				if info.Liquidity.Cmp(big.NewInt(0)) <= 0 {
					return false
				}
			}
			return true
		}).ToSlice(&successedPrices)

		// Calculate price impact and fee rate
		var wg sync.WaitGroup
		results := make(chan struct {
			Index       int64
			PriceImpact float64
			FeeRate     float64
			Err         error
		}, len(successedPrices))
		for i, price := range successedPrices {
			wg.Add(1)
			go func(ch chan struct {
				Index       int64
				PriceImpact float64
				FeeRate     float64
				Err         error
			}) {
				defer wg.Done()
				defer func() {
					if err := recover(); err != nil {
						ch <- struct {
							Index       int64
							PriceImpact float64
							FeeRate     float64
							Err         error
						}{
							Index: int64(i),
							Err:   fmt.Errorf("%s(panic)", err),
						}
					}
				}()

				totalPriceImpact := big.NewFloat(0)
				totalAmount := big.NewFloat(0)
				stepAmountOut := amountOut
				var totalFeeRate float64 = 1
				for j, pool := range price.Pools {
					info := poolMap[pool.Address.Hex()]
					_, from, to, err := price.Path.GetPool(int64(j))
					if err != nil {
						ch <- struct {
							Index       int64
							PriceImpact float64
							FeeRate     float64
							Err         error
						}{
							Index: int64(i),
							Err:   fmt.Errorf("%s", err),
						}
						return
					}
					stepPrices, err := cacher.QuoteExactOutput(ctx, []quoter.Pool{pool}, *from, *to, stepAmountOut, 0, 1)
					if err != nil {
						ch <- struct {
							Index       int64
							PriceImpact float64
							FeeRate     float64
							Err         error
						}{
							Index: int64(i),
							Err:   fmt.Errorf("%s", err),
						}
						return
					}
					if len(stepPrices) < 1 {
						ch <- struct {
							Index       int64
							PriceImpact float64
							FeeRate     float64
							Err         error
						}{
							Index: int64(i),
							Err:   fmt.Errorf("%s", "Can not quote step price"),
						}
						return
					}
					toToken, err := GetToken(chainId, to.Hex())
					if err != nil {
						ch <- struct {
							Index       int64
							PriceImpact float64
							FeeRate     float64
							Err         error
						}{
							Index: int64(i),
							Err:   err,
						}
						return
					}
					toAccuracy := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(toToken.Decimals), nil))
					stepAmountOutF := new(big.Float).Quo(new(big.Float).SetInt(stepAmountOut), toAccuracy)
					priceImpact := sonex.CalculatePriceImpact(info.Slot0.SqrtPriceX96, price.SqrtPriceX96AfterList[j])
					totalPriceImpact = new(big.Float).Add(totalPriceImpact, new(big.Float).Mul(priceImpact, stepAmountOutF))
					totalAmount = new(big.Float).Add(totalAmount, stepAmountOutF)
					totalFeeRate = totalFeeRate * (1 - float64(pool.Fee)/1000000)
					stepAmountOut = stepPrices[0].Value
				}
				avgPriceImpact, _ := new(big.Float).Quo(totalPriceImpact, totalAmount).Float64()
				ch <- struct {
					Index       int64
					PriceImpact float64
					FeeRate     float64
					Err         error
				}{
					Index:       int64(i),
					PriceImpact: avgPriceImpact,
					FeeRate:     1 - totalFeeRate,
				}
			}(results)
		}
		wg.Wait()
		close(results)
		for result := range results {
			if result.Err != nil {
				return nil, nil, result.Err
			}
			successedPrices[result.Index].PriceImpact = result.PriceImpact
			successedPrices[result.Index].FeeRate = result.FeeRate
		}
	}

	return successedPrices, failedPrices, nil
}

type SummarizePositionsResult_Detail struct {
	TokenId    *big.Int
	PoolAddres common.Address
	BalanceUSD *big.Int
	EarnedUSD  *big.Int
}

type SummarizePositionsResult struct {
	TotalBalanceUSD *big.Int
	TotalEarnedUSD  *big.Int
	Details         []SummarizePositionsResult_Detail
}

// Returns summary of `positions`
func SummarizePositions(conf config.Config, ctx context.Context, chainId int64, positions []*sonex.Position) (*big.Int, *big.Int, []SummarizePositionsResult_Detail, error) {
	totalBalanceUSD := big.NewInt(0)
	totalEarnedUSD := big.NewInt(0)
	details := []SummarizePositionsResult_Detail{}
	var wg sync.WaitGroup
	results := make(chan struct {
		BalanceUSD *big.Int
		EarnedUSD  *big.Int
		Err        error
	}, len(positions))

	for _, position := range positions {
		wg.Add(1)

		go func(ch chan struct {
			BalanceUSD *big.Int
			EarnedUSD  *big.Int
			Err        error
		}) {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					ch <- struct {
						BalanceUSD *big.Int
						EarnedUSD  *big.Int
						Err        error
					}{
						Err: fmt.Errorf("%s", err),
					}
				}
			}()

			balanceUSD := big.NewInt(0)
			earnedUSD := big.NewInt(0)

			amount0, amount1 := position.Amounts()
			fee0, fee1 := position.UnclaimedFees()

			token0CachePrice, _, _, err := GetPrices(ctx, nil, chainId, position.Token0.Hex(), true, 2)
			if err != nil {
				ch <- struct {
					BalanceUSD *big.Int
					EarnedUSD  *big.Int
					Err        error
				}{Err: err}
				return
			}
			token0, err := GetToken(chainId, position.Token0.Hex())
			if err != nil {
				ch <- struct {
					BalanceUSD *big.Int
					EarnedUSD  *big.Int
					Err        error
				}{Err: err}
				return
			}
			tokne0Decimals := new(big.Int).Exp(big.NewInt(10), big.NewInt(token0.Decimals), nil)
			token0Price := big.NewInt(0)
			if token0CachePrice != nil {
				token0Price = token0CachePrice.Value
			}
			balanceUSD = new(big.Int).Add(balanceUSD, new(big.Int).Div(new(big.Int).Mul(amount0, token0Price), tokne0Decimals))
			earnedUSD = new(big.Int).Add(earnedUSD, new(big.Int).Div(new(big.Int).Mul(fee0, token0Price), tokne0Decimals))

			token1CachePrice, _, _, err := GetPrices(ctx, nil, chainId, position.Token1.Hex(), true, 2)
			if err != nil {
				ch <- struct {
					BalanceUSD *big.Int
					EarnedUSD  *big.Int
					Err        error
				}{Err: err}
				return
			}
			token1, err := GetToken(chainId, position.Token1.Hex())
			if err != nil {
				ch <- struct {
					BalanceUSD *big.Int
					EarnedUSD  *big.Int
					Err        error
				}{Err: err}
				return
			}
			tokne1Decimals := new(big.Int).Exp(big.NewInt(10), big.NewInt(token1.Decimals), nil)
			token1Price := big.NewInt(0)
			if token1CachePrice != nil {
				token1Price = token1CachePrice.Value
			}
			balanceUSD = new(big.Int).Add(balanceUSD, new(big.Int).Quo(new(big.Int).Mul(amount1, token1Price), tokne1Decimals))
			earnedUSD = new(big.Int).Add(earnedUSD, new(big.Int).Quo(new(big.Int).Mul(fee1, token1Price), tokne1Decimals))
			details = append(details, SummarizePositionsResult_Detail{
				TokenId:    position.TokenId,
				PoolAddres: position.Pool.Address,
				BalanceUSD: balanceUSD,
				EarnedUSD:  earnedUSD,
			})

			ch <- struct {
				BalanceUSD *big.Int
				EarnedUSD  *big.Int
				Err        error
			}{BalanceUSD: balanceUSD, EarnedUSD: earnedUSD}
		}(results)
	}

	wg.Wait()
	close(results)

	for result := range results {
		if result.Err != nil {
			return nil, nil, nil, result.Err
		}
		totalBalanceUSD = new(big.Int).Add(totalBalanceUSD, result.BalanceUSD)
		totalEarnedUSD = new(big.Int).Add(totalEarnedUSD, result.EarnedUSD)
	}

	return totalBalanceUSD, totalEarnedUSD, details, nil
}

// Query current positions of an LP of `userAddress`, returns summary.
func SummarizeCurrentPositions(conf config.Config, ctx context.Context, chainId int64, userAddress common.Address, stablecoinAddress common.Address) (*SummarizePositionsResult, []*sonex.Position, error) {
	conn := repository.NewConn(conf.PostgreSQL.DateSource())
	model := repository.NewEventSonexPositionManagerTransferModel(conn)
	items, err := model.GetTransferEvents(ctx, chainId, userAddress.Hex(), true)
	if err != nil {
		return nil, nil, err
	}

	var tokenIds []*big.Int
	for _, item := range items {
		if item.IsBurnt {
			continue
		}
		tokenIds = append(tokenIds, big.NewInt(item.TokenID))
	}

	totalBalanceUSDResult := big.NewInt(0)
	totalEarnedUSDResult := big.NewInt(0)
	detailsResult := []SummarizePositionsResult_Detail{}

	positions := []*sonex.Position{}
	if len(tokenIds) > 0 {
		sonexConf := conf.GetSonexConfig(chainId)
		if sonexConf == nil {
			return nil, nil, fmt.Errorf("%s", "Miss sonex configure")
		}
		caller, _ := sonex.NewStaticCaller(conf)
		positions, err = caller.BatchGetPositions(
			context.Background(),
			sonexConf.ChainId,
			common.HexToAddress(sonexConf.PositionManager),
			tokenIds,
		)
		if err != nil {
			return nil, nil, err
		}

		totalBalanceUSD, totalEarnedUSD, details, err := SummarizePositions(conf, ctx, chainId, positions)
		if err != nil {
			return nil, nil, err
		}
		totalBalanceUSDResult = totalBalanceUSD
		totalEarnedUSDResult = totalEarnedUSD
		detailsResult = details
	}

	return &SummarizePositionsResult{
		TotalBalanceUSD: totalBalanceUSDResult,
		TotalEarnedUSD:  totalEarnedUSDResult,
		Details:         detailsResult,
	}, positions, nil
}

// Query closed positions of an LP of `userAddress`, returns summary.
func SummarizeClosedPositions(conf config.Config, ctx context.Context, chainId int64, userAddress common.Address, stablecoinAddress common.Address) (*SummarizePositionsResult, error) {
	conn := repository.NewConn(conf.PostgreSQL.DateSource())
	model := repository.NewEventSonexPositionManagerTransferModel(conn)
	items, err := model.GetTransferEvents(ctx, chainId, userAddress.Hex(), false)
	if err != nil {
		return nil, err
	}

	var tokenIds []*big.Int
	for _, item := range items {
		if !item.IsBurnt {
			continue
		}
		tokenIds = append(tokenIds, big.NewInt(item.TokenID))
	}

	totalBalanceUSDResult := big.NewInt(0)
	totalEarnedUSDResult := big.NewInt(0)
	detailsResult := []SummarizePositionsResult_Detail{}

	if len(tokenIds) > 0 {
		sonexConf := conf.GetSonexConfig(chainId)
		if sonexConf == nil {
			return nil, fmt.Errorf("%s", "Miss sonex configure")
		}
		caller, _ := sonex.NewStaticCaller(conf)
		positions, err := caller.BatchGetPositions(
			context.Background(),
			sonexConf.ChainId,
			common.HexToAddress(sonexConf.PositionManager),
			tokenIds,
		)
		if err != nil {
			return nil, err
		}

		totalBalanceUSD, totalEarnedUSD, details, err := SummarizePositions(conf, ctx, chainId, positions)
		if err != nil {
			return nil, err
		}
		totalBalanceUSDResult = totalBalanceUSD
		totalEarnedUSDResult = totalEarnedUSD
		detailsResult = details
	}

	return &SummarizePositionsResult{
		TotalBalanceUSD: totalBalanceUSDResult,
		TotalEarnedUSD:  totalEarnedUSDResult,
		Details:         detailsResult,
	}, nil
}

// Batch query tokens' total supplies
func BatchGetTotalSupplies(ctx context.Context, chainId int64, tokens []common.Address) (map[common.Address]*big.Int, error) {
	tokenCachersLock.RLock()
	cacher, ok := tokenCachers[chainId]
	tokenCachersLock.RUnlock()
	if !ok {
		return nil, fmt.Errorf("%s", "Miss token cacher")
	}
	return cacher.BatchGetTotalSupplies(ctx, tokens)
}
