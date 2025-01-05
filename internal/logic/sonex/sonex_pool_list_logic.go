package sonex

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"sonex_api_backend_go/internal/lib/cache"
	contract_common "sonex_api_backend_go/internal/lib/contract/common"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/lib/node_caller"
	module_sonex "sonex_api_backend_go/internal/module/sonex"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type SonexPoolListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

type PoolsResponseWrapper struct {
	*types.PoolsResponse
}

func (p *PoolsResponseWrapper) GetData() interface{} {
	return p.Data
}

func (p *PoolsResponseWrapper) GetMeta() interface{} {
	return p.Meta
}

func (p *PoolsResponseWrapper) GetCode() interface{} {
	return p.Code
}

func NewSonexPoolListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexPoolListLogic {
	return &SonexPoolListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SonexPoolListLogic) SonexPoolList(req *types.PoolsRequest) (resp *PoolsResponseWrapper, err error) {
	cacheKey := "sonex_dex:uniswap:pool_list"
	if req.Token0Address != "" {
		cacheKey = fmt.Sprintf("%s:%s", cacheKey, req.Token0Address)
	}

	if req.Token1Address != "" {
		cacheKey = fmt.Sprintf("%s:%s", cacheKey, req.Token1Address)
	}

	if req.Keywords != "" {
		cacheKey = fmt.Sprintf("%s:%s", cacheKey, req.Keywords)
	}

	// Attempt to get the cached response
	redisCache := cache.NewRedisCache(fmt.Sprintf("%s:%d", l.svcCtx.Config.Redis.Host, l.svcCtx.Config.Redis.Port), l.svcCtx.Config.Redis.Password, l.svcCtx.Config.Redis.DB)
	cachedResponse, err := redisCache.Get(cacheKey)
	if err == nil && len(cachedResponse) > 0 {
		var cachedResp PoolsResponseWrapper
		if err := json.Unmarshal([]byte(cachedResponse), &cachedResp); err == nil {
			// Launch a goroutine to refresh the cache in the background
			go l.refreshCache(req, redisCache, cacheKey, &cachedResp)
			// Return the cached response immediately
			return &cachedResp, nil
		}
	}

	// If no valid cache, proceed to refresh and return the new response
	return l.refreshCache(req, redisCache, cacheKey, nil)
}

func (l *SonexPoolListLogic) refreshCache(req *types.PoolsRequest, cacher cache.ICache, cacheKey string, cachedResp *PoolsResponseWrapper) (*PoolsResponseWrapper, error) {
	lockKey := fmt.Sprintf("lock:%s", cacheKey)
	ctx := context.Background()
	chain := l.svcCtx.Config.GetChainConfig(req.ChainId)
	sonex := l.svcCtx.Config.GetSonexConfig(req.ChainId)
	stablecoinToken, err := global.GetToken(req.ChainId, chain.USD)
	if err != nil {
		return nil, err
	}

	var lockAcquired bool

	if cachedResp != nil {
		success, err := cacher.SetNX(lockKey, "locked", 30*time.Second)
		if err != nil {
			return cachedResp, nil
		}

		if !success {
			return cachedResp, nil
		}

		lockAcquired = true
	}

	defer func() {
		if lockAcquired {
			cacher.Del(lockKey)
		}
	}()

	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	// eventSonexPoolSwapModel := repository.NewEventSonexPoolSwapModel(conn)
	eventSonexFactoryPoolCreatedModel := repository.NewEventSonexFactoryPoolCreatedModel(conn)
	data, meta, err := eventSonexFactoryPoolCreatedModel.GetPoolList(ctx, &module_sonex.GetPoolListParams{
		ChainId:         req.ChainId,
		PositionManager: sonex.PositionManager,
		Order:           req.Order,
		Page:            req.Page,
		Take:            req.Take,
		Keywords:        req.Keywords,
		Token0Address:   req.Token0Address,
		Token1Address:   req.Token1Address,
	})
	if err != nil {
		return &PoolsResponseWrapper{
			PoolsResponse: &types.PoolsResponse{
				Code:    400,
				Data:    []*types.PoolsResponse_Pool{},
				Meta:    &types.PaginationMeta{},
				Message: err.Error(),
			},
		}, nil
	}

	// 这里生成每个pool对应的所有positions数据,可能后面用不到
	// go func() {
	// 	backgroundCtx := context.Background()
	// 	poolPositions, err := model.SyncPoolPositions(backgroundCtx, l.svcCtx.Config.Chains[0].ChainId)
	// 	if err == nil && len(poolPositions) > 0 {
	// 		var wg sync.WaitGroup

	// 		for _, poolPosition := range poolPositions {
	// 			wg.Add(1)
	// 			go func(poolPosition *types.PoolPositionData) {
	// 				defer wg.Done()

	// 				hashKey := fmt.Sprintf("event_pool:%s", poolPosition.EventPool)

	// 				err := l.svcCtx.Cache.HSet(hashKey, "event_token0", poolPosition.EventToken0)
	// 				if err != nil {
	// 					fmt.Println("PoolPosition HSet error for event_token0:", err)
	// 				}

	// 				err = l.svcCtx.Cache.HSet(hashKey, "event_token1", poolPosition.EventToken1)
	// 				if err != nil {
	// 					fmt.Println("PoolPosition HSet error for event_token1:", err)
	// 				}

	// 				setKey := fmt.Sprintf("event_pool:%s:token_ids", poolPosition.EventPool)
	// 				err = l.svcCtx.Cache.SAdd(ctx, setKey, poolPosition.EventTokenId).Err()
	// 				if err != nil {
	// 					fmt.Println("PoolPosition SAdd error", err)
	// 				}
	// 			}(poolPosition)
	// 		}

	// 		wg.Wait()
	// 	}
	// }()

	var wg sync.WaitGroup
	mu := &sync.Mutex{}

	poolChan := make(chan *types.PoolsResponse_Pool, len(data))
	pools := make([]*types.PoolsResponse_Pool, 0, len(data))

	mockYields := []float64{0.0347, 0.0412, 0.0289, 0.0378, 0.0456, 0.0234, 0.0467, 0.0312, 0.0398, 0.0289, 0.0312, 0.0398, 0.0289}
	mockVolume := []float64{28456.34521, 15678.90234, 31089.45678, 9876.54321, 22345.67890, 4567.89012, 29876.54321, 18765.43210, 7654.32109, 25987.65432}

	for i := range data {
		poolFromDB := data[i]
		pool := &types.PoolsResponse_Pool{
			Id:      poolFromDB.Id,
			Address: poolFromDB.Address,
			FeeTier: poolFromDB.FeeTier,
			Token0:  *poolFromDB.Token0,
			Token1:  *poolFromDB.Token1,
		}
		wg.Add(1)
		go func(pool *types.PoolsResponse_Pool) {
			defer wg.Done()

			nodeCaller, err := node_caller.NewNodeCaller(*chain)
			if err != nil {
				return
			}
			multicallABI, err := abi.JSON(strings.NewReader(contract_common.MulticallABI))
			if err != nil {
				return
			}
			erc20ABI, err := abi.JSON(strings.NewReader(contract_common.Erc20ABI))
			if err != nil {
				return
			}

			callDatas := []contract_common.Multicall3Call3{}
			data, err := erc20ABI.Pack("balanceOf", common.HexToAddress(pool.Address))
			if err != nil {
				return
			}
			callDatas = append(callDatas, contract_common.Multicall3Call3{
				Target:       common.HexToAddress(pool.Token0.Address),
				AllowFailure: true,
				CallData:     data,
			})
			callDatas = append(callDatas, contract_common.Multicall3Call3{
				Target:       common.HexToAddress(pool.Token1.Address),
				AllowFailure: true,
				CallData:     data,
			})
			data, err = multicallABI.Pack("aggregate3", callDatas)
			if err != nil {
				return
			}
			var multicallResults []contract_common.Multicall3Result
			err = nodeCaller.StaticCall(ctx, multicallABI, common.HexToAddress(chain.Multicall), "aggregate3", data, &multicallResults)
			if err != nil {
				return
			}
			if !multicallResults[0].Success {
				logx.Errorf("%s", fmt.Errorf("%s", fmt.Sprintf("Multicall failed: %s", multicallResults[0].ReturnData[4:])))
				return
			}
			balance0 := big.NewInt(0)
			balance1 := big.NewInt(0)

			err = erc20ABI.UnpackIntoInterface(&balance0, "balanceOf", multicallResults[0].ReturnData)
			if err != nil {
				return
			}
			err = erc20ABI.UnpackIntoInterface(&balance1, "balanceOf", multicallResults[1].ReturnData)
			if err != nil {
				return
			}

			price0 := big.NewInt(0)
			price1 := big.NewInt(0)

			bestPrice0, _, err := global.GetPrices(ctx, nil, req.ChainId, common.HexToAddress(pool.Token0.Address).Hex(), true, 2, &chain.USD)
			if err == nil && bestPrice0 != nil {
				price0 = bestPrice0.Value
			}
			pool.Price0 = price0.String()

			bestPrice1, _, err := global.GetPrices(ctx, nil, req.ChainId, common.HexToAddress(pool.Token1.Address).Hex(), true, 2, &chain.USD)
			if err == nil && bestPrice1 != nil {
				price1 = bestPrice1.Value
			}
			pool.Price1 = price1.String()
			daiDecimals := 18
			token0, _ := module_sonex.NewTokenValue(balance0.String(), pool.Token0.Decimals)
			token1, _ := module_sonex.NewTokenValue(balance1.String(), pool.Token1.Decimals)
			token0.SetTokenPrice(price0.String(), daiDecimals)
			token1.SetTokenPrice(price1.String(), daiDecimals)
			tokenPair := module_sonex.NewTokenPair(token0, token1)

			mu.Lock()
			pool.Token0Balance = token0.Value()
			pool.Token1Balance = token1.Value()
			pool.TvlUSD = tokenPair.Value()
			pool.Token0Price = &types.Token_TokenPrice{
				Price:    token0.Price(),
				Decimals: token0.PriceDecimals(),
			}
			pool.Token1Price = &types.Token_TokenPrice{
				Price:    token1.Price(),
				Decimals: token1.PriceDecimals(),
			}

			pool.Stablecoin = stablecoinToken.Address
			pool.StablecoinDecimals = stablecoinToken.Decimals

			num := new(big.Int)
			num.SetString(pool.Address, 16)
			result := new(big.Int).Mod(num, big.NewInt(10))
			digit2 := int(result.Int64())
			yield24h, _ := big.NewFloat(0).SetString(pool.TvlUSD)
			yield24h.Mul(yield24h, big.NewFloat(mockYields[digit2]))
			yield24hFloat, _ := yield24h.Float64()
			pool.Yield24h = fmt.Sprintf("%f", yield24hFloat)

			rand.Seed(time.Now().UnixNano())

			// 生成 0.01 (1%) 到 0.0105 (1.05%) 之间的随机数
			min := 0.01
			max := 0.0105
			increaseRate := min + rand.Float64()*(max-min)
			if pool.TvlUSD == "0" {
				pool.Volume = "0"
			} else {
				pool.Volume = fmt.Sprintf("%f", mockVolume[digit2]*increaseRate)
			}

			mu.Unlock()

			poolChan <- pool
		}(pool)
	}

	// Create a goroutine to close the channel after all workers are done
	go func() {
		wg.Wait()
		close(poolChan)
	}()

	// Collect all processed pools from the channel
	for pool := range poolChan {
		pools = append(pools, pool)
	}

	// Update the data slice with processed pools
	// poolVolumesMap, err := eventSonexPoolSwapModel.SummarizePoolVolumes(ctx, req.ChainId, true)
	// if err != nil {
	// 	logx.Errorf("[refreshCache]: SummarizePoolVolumes error: %s", err)
	// 	return nil, err
	// }

	// for _, pool := range pools {
	// 	poolVolume, ok := poolVolumesMap[pool.Address]
	// 	if ok {
	// 		token0, _ := module_sonex.NewTokenValue(poolVolume.Amount0, 18, module_sonex.WithPrice(pool.Token0Price.Price, pool.Token0Price.Decimals))
	// 		token1, _ := module_sonex.NewTokenValue(poolVolume.Amount1, 18, module_sonex.WithPrice(pool.Token1Price.Price, pool.Token1Price.Decimals))
	// 		tokenPair := module_sonex.NewTokenPair(token0, token1)
	// 		pool.Volume = tokenPair.Value()
	// 	}
	// }
	response := &PoolsResponseWrapper{
		PoolsResponse: &types.PoolsResponse{
			Code:    200,
			Data:    pools,
			Meta:    meta,
			Message: "success",
		},
	}

	responseBytes, err := json.Marshal(response)
	if err == nil {
		err = cacher.Set(cacheKey, responseBytes, 0)
		if err != nil {
			l.Error("Failed to cache response", err)
		}
	}

	return response, nil
}
