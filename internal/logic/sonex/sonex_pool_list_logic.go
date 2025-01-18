package sonex

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"sonex_api_backend_go/internal/handler/response"
	"sonex_api_backend_go/internal/lib/cache"
	"sonex_api_backend_go/internal/lib/global"
	module_sonex "sonex_api_backend_go/internal/module/sonex"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type SonexPoolListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexPoolListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexPoolListLogic {
	return &SonexPoolListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SonexPoolListLogic) SonexPoolList(req *types.PoolsRequest) (resp *response.ApiResponse, err error) {
	if req.ChainId <= 0 {
		return nil, fmt.Errorf("invalid parameter")
	}
	if req.PageSize == nil || *req.PageSize < 1 {
		var defaultSize int64 = 20
		req.PageSize = &defaultSize
	}
	if req.PageIndex == nil || *req.PageIndex < 1 {
		var defaultIndex int64 = 1
		req.PageIndex = &defaultIndex
	}

	cacheKey := "sonex_dex:uniswap:pool_list"
	if req.PageIndex != nil {
		cacheKey = fmt.Sprintf("%s:%d", cacheKey, *req.PageIndex)
	}

	if req.PageSize != nil {
		cacheKey = fmt.Sprintf("%s:%d", cacheKey, *req.PageSize)
	}

	if req.Order != "" {
		cacheKey = fmt.Sprintf("%s:%s", cacheKey, req.Order)
	}
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
		var cachedResp response.ApiResponse
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

func (l *SonexPoolListLogic) refreshCache(req *types.PoolsRequest, cacher cache.ICache, cacheKey string, cachedResp *response.ApiResponse) (*response.ApiResponse, error) {
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
	tokenPriceModel := repository.NewSonexTokenPriceRecordModel(conn)
	tokenModel := repository.NewTokenModel(conn)

	eventSonexFactoryPoolCreatedModel := repository.NewEventSonexFactoryPoolCreatedModel(conn)
	data, meta, err := eventSonexFactoryPoolCreatedModel.GetPoolList(ctx, &module_sonex.GetPoolListParams{
		ChainId:         req.ChainId,
		PositionManager: sonex.PositionManager,
		Order:           req.Order,
		Page:            *req.PageIndex,
		Take:            *req.PageSize,
		Keywords:        req.Keywords,
		Token0Address:   req.Token0Address,
		Token1Address:   req.Token1Address,
	})
	if err != nil {
		return &response.ApiResponse{
			Data: []*types.PoolsResponse_Pool{},
			Meta: &types.PageResponse{},
		}, err
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

	mockYields := []float64{0.00000347, 0.00000412, 0.00000289, 0.00000378, 0.00000456, 0.00000234, 0.00000467, 0.00000312, 0.00000398, 0.00000289, 0.00000312, 0.00000398, 0.00000289}

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

			balance0, balance1 := big.NewInt(0), big.NewInt(0)
			Multicall(
				ctx,
				*chain,
				Call{
					Contract: common.HexToAddress(pool.Token0.Address),
					Method:   "balanceOf",
					Args:     []interface{}{common.HexToAddress(pool.Address)},
					Result:   &balance0,
				},
				Call{
					Contract: common.HexToAddress(pool.Token1.Address),
					Method:   "balanceOf",
					Args:     []interface{}{common.HexToAddress(pool.Address)},
					Result:   &balance1,
				},
			)
			prices, err := tokenPriceModel.GetCurrentPrices([]string{common.HexToAddress(pool.Token0.Address).Hex(), common.HexToAddress(pool.Token1.Address).Hex()})
			if err != nil {
				logx.Errorf("GetCurrentPrices:%v ", err)
				return
			}

			price0 := big.NewInt(0)
			price0Decimal := int64(0)
			price1 := big.NewInt(0)
			price1Decimal := int64(0)

			for _, v := range prices {
				if strings.EqualFold(v.TokenAddress, common.HexToAddress(pool.Token0.Address).Hex()) {
					price0, _ = new(big.Int).SetString(v.Price, 10)
					price0Decimal = int64(v.UsdDecimals)
				}

				if strings.EqualFold(v.TokenAddress, common.HexToAddress(pool.Token1.Address).Hex()) {
					price1, _ = new(big.Int).SetString(v.Price, 10)
					price1Decimal = int64(v.UsdDecimals)
				}
			}

			_tokenInfo, err := tokenModel.QueryByChainIdAddress(context.Background(), chain.ChainId, []string{common.HexToAddress(pool.Token0.Address).Hex(), common.HexToAddress(pool.Token1.Address).Hex()})
			if err != nil {
				logx.Errorf("QueryByChainIdAddress:%v ", err)
				return
			}

			_decimals0 := int64(0)

			_decimals1 := int64(0)

			for _, v := range _tokenInfo {
				if strings.EqualFold(v.Address, common.HexToAddress(pool.Token0.Address).Hex()) {
					_decimals0 = v.Decimals

				}

				if strings.EqualFold(v.Address, common.HexToAddress(pool.Token1.Address).Hex()) {
					_decimals1 = v.Decimals
				}
			}
			pool.Price0 = price0.String()

			pool.Price1 = price1.String()
			mu.Lock()
			pool.Token0Balance = balance0.String()
			pool.Token1Balance = balance1.String()

			if balance0.Cmp(big.NewInt(0)) < 0 {
				balance0 = big.NewInt(1)
			}

			if balance1.Cmp(big.NewInt(0)) < 0 {
				balance1 = big.NewInt(1)
			}

			// for tvl calucate
			balance0Float, _ := new(big.Float).SetString(balance0.String())
			_balance0divdeimals0 := new(big.Float).Quo(balance0Float, new(big.Float).SetFloat64(math.Pow10(int(_decimals0))))

			balance1Float, _ := new(big.Float).SetString(balance1.String())
			_balance1divdeimals1 := new(big.Float).Quo(balance1Float, new(big.Float).SetFloat64(math.Pow10(int(_decimals1))))

			_price0divUsdDecimals := new(big.Float).Quo(new(big.Float).SetInt(price0), new(big.Float).SetFloat64(math.Pow10(int(price0Decimal))))
			_price1divUsdDecimals := new(big.Float).Quo(new(big.Float).SetInt(price1), new(big.Float).SetFloat64(math.Pow10(int(price1Decimal))))

			_token0Volume := new(big.Float).Mul(_balance0divdeimals0, _price0divUsdDecimals)
			_token1Volume := new(big.Float).Mul(_balance1divdeimals1, _price1divUsdDecimals)

			_tvlUSD := new(big.Float).Add(_token0Volume, _token1Volume)

			pool.Liquidity = _tvlUSD.String()
			pool.TvlUSD = _tvlUSD.String()
			pool.Token0Price = &types.Token_TokenPrice{
				Price:    price0.String(),
				Decimals: int(_decimals0),
			}
			pool.Token1Price = &types.Token_TokenPrice{
				Price:    price1.String(),
				Decimals: int(_decimals1),
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

	// sort pool tvl
	sort.SliceStable(pools, func(i, j int) bool {
		tvlUsd_i, _ := big.NewFloat(0).SetString(pools[i].TvlUSD)
		tvlUsd_j, _ := big.NewFloat(0).SetString(pools[j].TvlUSD)
		return tvlUsd_i.Cmp(tvlUsd_j) > 0
	})

	// calculate volume
	_swapAmoutEthAndUsdMap, _, err := CalucatedVolume24H(conn, req.ChainId)
	if err != nil {
		return &response.ApiResponse{
			Data: []*types.PoolsResponse_Pool{},
			Meta: &types.PageResponse{},
		}, err
	}
	// get Price from value
	for _, pool := range pools {
		for _, v := range _swapAmoutEthAndUsdMap {
			if strings.EqualFold(pool.Address, v.ExtraPool) {
				pool.Volume = v.Value
			}
		}
	}
	response := &response.ApiResponse{
		Data: pools,
		Meta: &types.PageResponse{
			PageSize:  int64(meta.Take),
			PageIndex: int64(meta.Page),
			Total:     int64(meta.ItemCount),
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

// 如果 _specialPoolAddress 不为空，就把 _specialPoolAddress对应的volume24H 返回
func CalucatedVolume24H(conn sqlx.SqlConn, _chainId int64, _specialPoolAddress ...string) ([]*repository.PoolVolumeUsd, string, error) {

	eventSonexPoolSwapModel := repository.NewEventSonexPoolSwapModel(conn)
	ctx := context.Background()
	// 获取 event_sonex_pool_swap 里面的交易详情
	poolVolumesMap, err := eventSonexPoolSwapModel.SummarizePoolVolumesUsd(ctx, _chainId, true)
	if err != nil {
		logx.Errorf("[refreshCache]: SummarizePoolVolumesUsd error: %s", err)
		return nil, "", err
	}
	//  poolVolumesMapAfterProcessedItems 条目

	_specialPoolAddressVolume24H := "0"
	for _, v := range poolVolumesMap {
		if len(_specialPoolAddress) == 1 {
			if strings.EqualFold(v.ExtraPool, _specialPoolAddress[0]) {
				_specialPoolAddressVolume24H = v.Value
			}
		}
	}

	return poolVolumesMap, _specialPoolAddressVolume24H, nil
}
