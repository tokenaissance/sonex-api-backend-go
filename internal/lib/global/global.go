package global

import (
	"context"
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/cache"
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
var cronExecuter *cron.Cron
var idGenerater *id_generater.IdGenerater
var tokenCachers map[int64]*cache.TokenCacher
var tokenCachersLock sync.RWMutex

func init() {
	cronExecuter = cron.New(cron.WithSeconds(), cron.WithLocation(time.UTC))
	tokenCachers = map[int64]*cache.TokenCacher{}

	cronExecuter.Start()
}

func Initialize(config config.Config) error {
	c = config

	// Initialize token cacher
	for _, chain := range c.Chains {
		sonexConf := c.GetSonexConfig(chain.ChainId)
		if sonexConf == nil {
			return fmt.Errorf("%s", "Miss sonex configure")
		}

		tokenCacher := cache.NewTokenCacher(chain, *sonexConf)
		tokenCachersLock.Lock()
		tokenCachers[chain.ChainId] = tokenCacher
		tokenCachersLock.Unlock()
		for _, token := range chain.Tokens {
			err := SetToken(chain.ChainId, token.Address, token.IsStablecoin)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func AddCronJob(cron string, job func()) (cron.EntryID, error) {
	return cronExecuter.AddFunc(cron, job)
}

func InitIdGenerater(workerId int64) error {
	generater, err := id_generater.NewIdGenerater(workerId)
	if err != nil {
		return err
	}
	idGenerater = generater
	return nil
}

func NewId() (int64, error) {
	return idGenerater.Generate()
}

func GetToken(chainId int64, address string) (*cache.Token, error) {
	tokenCachersLock.RLock()
	cacher, ok := tokenCachers[chainId]
	tokenCachersLock.RUnlock()
	if !ok {
		return nil, fmt.Errorf("%s", "Miss token cacher")
	}
	return cacher.Get(address)
}

func SetToken(chainId int64, address string, isStablecoin bool) error {
	tokenCachersLock.RLock()
	cacher, ok := tokenCachers[chainId]
	tokenCachersLock.RUnlock()
	if !ok {
		return fmt.Errorf("%s", "Miss token cacher")
	}
	err := cacher.Set(address, isStablecoin)
	if err != nil {
		return err
	}
	return nil
}

func GetStablecoin(chainId int64) []cache.Token {
	tokenCachersLock.RLock()
	cacher, ok := tokenCachers[chainId]
	tokenCachersLock.RUnlock()
	if !ok {
		return nil
	}
	return cacher.GetStablecoins()
}

func GetPrices(ctx context.Context, pools []quoter.Pool, chainId int64, token string, isCache bool, maxLength int64, wantStablecoin *string) (*cache.Price, []cache.Price, error) {
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

	var wantStablecoinAddress *common.Address
	if wantStablecoin != nil {
		temp := common.HexToAddress(*wantStablecoin)
		wantStablecoinAddress = &temp
	}
	prices, err := cacher.GetPrices(pools, common.HexToAddress(token), isCache, maxLength, wantStablecoinAddress)
	if err != nil {
		return nil, nil, err
	}

	// Best price
	var bestPrice *cache.Price
	if len(prices) > 0 {
		// Sort by path length
		linq.From(prices).OrderByT(func(item cache.Price) int { return len(item.Pools) }).ToSlice(&prices)
		pathLength := len(prices[0].Pools)
		linq.From(prices).WhereT(func(item cache.Price) bool {
			return len(item.Pools) == pathLength
		}).ToSlice(&prices)
		poolMap := map[string]*big.Int{}
		linq.From(prices).ForEachT(func(price cache.Price) {
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
		data, err := poolABI.Pack("liquidity")
		if err != nil {
			return nil, nil, err
		}
		for _, pool := range usedPools {
			callDatas = append(callDatas, contract_common.Multicall3Call{
				Target:   pool,
				CallData: data,
			})
		}
		multicallABI, err := abi.JSON(strings.NewReader(contract_common.MulticallABI))
		if err != nil {
			return nil, nil, err
		}
		data, err = multicallABI.Pack("aggregate", callDatas)
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
		for i, value := range result.ReturnData {
			poolMap[usedPools[i].Hex()] = new(big.Int).SetBytes(value)
		}
		maxLiquidity := common.Big0
		for _, price := range prices {
			liquidity := common.Big0
			for _, pool := range price.Pools {
				liquidity = new(big.Int).Add(liquidity, poolMap[pool.Address.Hex()])
			}
			if maxLiquidity.Cmp(liquidity) < 0 {
				maxLiquidity = liquidity
				bestPrice = &price
			}
		}
	}

	return bestPrice, prices, nil
}

func QuoteExactInput(ctx context.Context, pools []quoter.Pool, chainId int64, from string, to string, amountIn *big.Int, maxLength int64) ([]cache.Price, error) {
	tokenCachersLock.RLock()
	cacher, ok := tokenCachers[chainId]
	tokenCachersLock.RUnlock()
	if !ok {
		return nil, fmt.Errorf("%s", "Miss token cacher")
	}

	if pools == nil {
		conn := repository.NewConn(c.PostgreSQL.DateSource())
		eventSonexFactoryPoolCreatedModel := repository.NewEventSonexFactoryPoolCreatedModel(conn)
		records, err := eventSonexFactoryPoolCreatedModel.QueryByChainId(ctx, chainId)
		if err != nil {
			return nil, err
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

	prices, err := cacher.QuoteExactInput(pools, common.HexToAddress(from), common.HexToAddress(to), amountIn, maxLength)
	if err != nil {
		return nil, err
	}

	if len(prices) > 0 {
		poolMap := map[string]*struct {
			Liquidity *big.Int
			Slot0     *sonex.Slot0
		}{}
		linq.From(prices).ForEachT(func(price cache.Price) {
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
			return nil, fmt.Errorf("%s", "Miss chain configure")
		}
		if chain.Multicall == "" || !common.IsHexAddress(chain.Multicall) {
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
		nodeCaller, err := node_caller.NewNodeCaller(*chain)
		if err != nil {
			return nil, err
		}
		var result struct {
			BlockNumber *big.Int `json:"blockNumber"`
			ReturnData  [][]byte `json:"returnData"`
		}
		err = nodeCaller.StaticCall(ctx, multicallABI, common.HexToAddress(chain.Multicall), "aggregate", data, &result)
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

		// Filter price if has zero liquidity pool
		linq.From(prices).WhereT(func(item cache.Price) bool {
			for _, pool := range item.Pools {
				info := poolMap[pool.Address.Hex()]
				if info.Liquidity.Cmp(big.NewInt(0)) <= 0 {
					return false
				}
			}
			return true
		}).ToSlice(&prices)

		// Calculate price impact
		Q192 := new(big.Float).SetInt(new(big.Int).Lsh(big.NewInt(1), 192))
		for i, price := range prices {
			totalPriceImpact := big.NewFloat(0)
			totalLiquidity := big.NewFloat(0)
			for j, pool := range price.Pools {
				info := poolMap[pool.Address.Hex()]
				currentPrice := new(big.Float).Quo(new(big.Float).Mul(new(big.Float).SetInt(info.Slot0.SqrtPriceX96), new(big.Float).SetInt(info.Slot0.SqrtPriceX96)), Q192)
				newPrice := new(big.Float).Quo(new(big.Float).Mul(new(big.Float).SetInt(price.SqrtPriceX96AfterList[j]), new(big.Float).SetInt(price.SqrtPriceX96AfterList[j])), Q192)
				liquidity := new(big.Float).SetInt(info.Liquidity)
				priceImpact := new(big.Float).Quo(new(big.Float).Sub(newPrice, currentPrice), currentPrice)
				totalPriceImpact = new(big.Float).Add(totalPriceImpact, new(big.Float).Mul(priceImpact, liquidity))
				totalLiquidity = new(big.Float).Add(totalLiquidity, liquidity)
			}
			avgPriceImpact, _ := new(big.Float).Quo(totalPriceImpact, totalLiquidity).Float64()
			prices[i].PriceImpact = avgPriceImpact
		}
	}

	return prices, nil
}

func QuoteExactOutput(ctx context.Context, pools []quoter.Pool, chainId int64, from string, to string, amountOut *big.Int, maxLength int64) ([]cache.Price, error) {
	tokenCachersLock.RLock()
	cacher, ok := tokenCachers[chainId]
	tokenCachersLock.RUnlock()
	if !ok {
		return nil, fmt.Errorf("%s", "Miss token cacher")
	}

	if pools == nil {
		conn := repository.NewConn(c.PostgreSQL.DateSource())
		eventSonexFactoryPoolCreatedModel := repository.NewEventSonexFactoryPoolCreatedModel(conn)
		records, err := eventSonexFactoryPoolCreatedModel.QueryByChainId(ctx, chainId)
		if err != nil {
			return nil, err
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

	prices, err := cacher.QuoteExactOutput(pools, common.HexToAddress(from), common.HexToAddress(to), amountOut, maxLength)
	if err != nil {
		return nil, err
	}

	if len(prices) > 0 {
		poolMap := map[string]*struct {
			Liquidity *big.Int
			Slot0     *sonex.Slot0
		}{}
		linq.From(prices).ForEachT(func(price cache.Price) {
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
			return nil, fmt.Errorf("%s", "Miss chain configure")
		}
		if chain.Multicall == "" || !common.IsHexAddress(chain.Multicall) {
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
		nodeCaller, err := node_caller.NewNodeCaller(*chain)
		if err != nil {
			return nil, err
		}
		var result struct {
			BlockNumber *big.Int `json:"blockNumber"`
			ReturnData  [][]byte `json:"returnData"`
		}
		err = nodeCaller.StaticCall(ctx, multicallABI, common.HexToAddress(chain.Multicall), "aggregate", data, &result)
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

		// Filter price if has zero liquidity pool
		linq.From(prices).WhereT(func(item cache.Price) bool {
			for _, pool := range item.Pools {
				info := poolMap[pool.Address.Hex()]
				if info.Liquidity.Cmp(big.NewInt(0)) <= 0 {
					return false
				}
			}
			return true
		}).ToSlice(&prices)

		// Calculate price impact
		Q192 := new(big.Float).SetInt(new(big.Int).Lsh(big.NewInt(1), 192))
		for i, price := range prices {
			totalPriceImpact := big.NewFloat(0)
			totalLiquidity := big.NewFloat(0)
			for j, pool := range price.Pools {
				info := poolMap[pool.Address.Hex()]
				currentPrice := new(big.Float).Quo(new(big.Float).Mul(new(big.Float).SetInt(info.Slot0.SqrtPriceX96), new(big.Float).SetInt(info.Slot0.SqrtPriceX96)), Q192)
				newPrice := new(big.Float).Quo(new(big.Float).Mul(new(big.Float).SetInt(price.SqrtPriceX96AfterList[j]), new(big.Float).SetInt(price.SqrtPriceX96AfterList[j])), Q192)
				liquidity := new(big.Float).SetInt(info.Liquidity)
				priceImpact := new(big.Float).Quo(new(big.Float).Sub(newPrice, currentPrice), currentPrice)
				totalPriceImpact = new(big.Float).Add(totalPriceImpact, new(big.Float).Mul(priceImpact, liquidity))
				totalLiquidity = new(big.Float).Add(totalLiquidity, liquidity)
			}
			avgPriceImpact, _ := new(big.Float).Quo(totalPriceImpact, totalLiquidity).Float64()
			prices[i].PriceImpact = avgPriceImpact
		}
	}

	return prices, nil
}
