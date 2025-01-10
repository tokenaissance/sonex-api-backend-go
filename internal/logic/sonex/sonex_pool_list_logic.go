package sonex

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/handler/response"
	"sonex_api_backend_go/internal/lib/cache"
	contract_common "sonex_api_backend_go/internal/lib/contract/common"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/lib/node_caller"
	module_sonex "sonex_api_backend_go/internal/module/sonex"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
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

			prices, err := tokenPriceModel.GetCurrentPrices([]string{common.HexToAddress(pool.Token0.Address).Hex(), common.HexToAddress(pool.Token1.Address).Hex()})
			if err != nil {
				logx.Errorf("GetCurrentPrices:%v ", err)
				return
			}

			price0 := big.NewInt(0)
			price0Decimal := int64(0)
			price1 := big.NewInt(0)
			// price1Decimal := int64(0)

			// bestPrice0, _, _, err := global.GetPrices(ctx, nil, req.ChainId, common.HexToAddress(pool.Token0.Address).Hex(), true, 2)
			// if err == nil && bestPrice0 != nil {
			// 	price0 = bestPrice0.Value
			// }

			for _, v := range prices {
				if strings.EqualFold(v.TokenAddress, common.HexToAddress(pool.Token0.Address).Hex()) {
					price0, _ = new(big.Int).SetString(v.Price, 10)
					price0Decimal = int64(v.UsdDecimals)
				}

				if strings.EqualFold(v.TokenAddress, common.HexToAddress(pool.Token1.Address).Hex()) {
					price1, _ = new(big.Int).SetString(v.Price, 10)
					// price1Decimal = int64(v.UsdDecimals)
				}

			}

			logx.Infof("XXXXXXXXXXXXXtoken0 %v balance0 :%v price0 :%v token1:%v balance1:%v   price1:%v  ", pool.Token0.Address, balance0, price0, pool.Token1.Address, balance1, price1)

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

			// bestPrice1, _, _, err := global.GetPrices(ctx, nil, req.ChainId, common.HexToAddress(pool.Token1.Address).Hex(), true, 2)
			// if err == nil && bestPrice1 != nil {
			// 	price1 = bestPrice1.Value
			// }
			pool.Price1 = price1.String()
			// daiDecimals := 18
			// token0, _ := module_sonex.NewTokenValue(balance0.String(), pool.Token0.Decimals)
			// token1, _ := module_sonex.NewTokenValue(balance1.String(), pool.Token1.Decimals)
			// token0.SetTokenPrice(price0.String(), pool.Token0.Decimals)
			// token1.SetTokenPrice(price1.String(), pool.Token1.Decimals)
			// tokenPair := module_sonex.NewTokenPair(token0, token1)

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

			// balance1Float, _ := new(big.Float).SetString(balance1.String())
			// _balance1divdeimals1 := new(big.Float).Quo(balance1Float, new(big.Float).SetFloat64(math.Pow10(int(_decimals1))))

			_price0divUsdDecimals := new(big.Float).Quo(new(big.Float).SetInt(price0), new(big.Float).SetFloat64(math.Pow10(int(price0Decimal))))
			// _price1divUsdDecimals := new(big.Float).Quo(new(big.Float).SetInt(price1), new(big.Float).SetFloat64(math.Pow10(int(price1Decimal))))

			//
			// 其中一个代币的价格*余额*2
			_token0Volume := new(big.Float).Mul(_balance0divdeimals0, _price0divUsdDecimals)

			_tvlUSD := new(big.Float).Mul(_token0Volume, big.NewFloat(2))

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

			rand.Seed(time.Now().UnixNano())

			// if strings.Compare(pool.TvlUSD, "0") == 0 {
			// 	pool.Volume = "0"
			// } else {
			// 	volumeFloat, _ := yield24h.Quo(yield24h, big.NewFloat(0.015)).Float64()
			// 	pool.Volume = fmt.Sprintf("%f", volumeFloat)
			// }

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

	// sort pool tvl
	sort.SliceStable(pools, func(i, j int) bool {
		tvlUsd_i, _ := big.NewFloat(0).SetString(pools[i].TvlUSD)
		tvlUsd_j, _ := big.NewFloat(0).SetString(pools[j].TvlUSD)
		return tvlUsd_i.Cmp(tvlUsd_j) > 0
	})

	// calculate volume
	_swapAmoutEthAndUsdMap, err := CalucatedVolume24H(req, conn, chain)
	if err != nil {
		return &response.ApiResponse{
			Data: []*types.PoolsResponse_Pool{},
			Meta: &types.PageResponse{},
		}, err
	}
	// get Price from value
	for _, pool := range pools {

		poolVolume, ok := _swapAmoutEthAndUsdMap[strings.ToLower(pool.Address)]
		if ok {
			pool.Volume = poolVolume.Volume24H.String()
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

// for calcuated volume
// key : pool Address
// body: usdt_amount,eth_Amount
// 计算方法: eth_amount*eth_price  + usd
type SwapAmoutEthAndUsd struct {
	Volume24H *big.Float //
	EthAmount *big.Float //
	UsdAmount *big.Float //
}

func CalucatedVolume24H(req *types.PoolsRequest, conn sqlx.SqlConn, chain *config.ChainConfig) (map[string]*SwapAmoutEthAndUsd, error) {

	/*
		USD: amount1
		WETH: amount0
	*/

	eventSonexPoolSwapModel := repository.NewEventSonexPoolSwapModel(conn)
	ctx := context.Background()
	// 获取 event_sonex_pool_swap 里面的交易详情
	poolVolumesMap, err := eventSonexPoolSwapModel.SummarizePoolVolumesUsd(ctx, req.ChainId, true)
	if err != nil {
		logx.Errorf("[refreshCache]: SummarizePoolVolumesUsd error: %s", err)
		return nil, err
	}
	tokenPriceModel := repository.NewSonexTokenPriceRecordModel(conn)
	//  poolVolumesMapAfterProcessedItems 条目
	_swapAmoutEthAndUsdMap := make(map[string]*SwapAmoutEthAndUsd)
	for _, _v := range poolVolumesMap {

		// pre process
		_valueString := "0"
		if !strings.EqualFold(_v.Value, "0") {
			_valueString = _v.Value
		}
		_valueFloat, _ := new(big.Float).SetString(_valueString)

		if (strings.EqualFold(strings.ToLower(_v.ExtraToken0), strings.ToLower(chain.WETH)) &&
			!strings.EqualFold(strings.ToLower(_v.ExtraToken1), strings.ToLower(chain.USD))) ||
			(strings.EqualFold(strings.ToLower(_v.ExtraToken1), strings.ToLower(chain.WETH)) &&
				!strings.EqualFold(strings.ToLower(_v.ExtraToken0), strings.ToLower(chain.USD))) {
			// 1, ETH/TOKENA 交易对
			_swapAmoutEthAndUsd := &SwapAmoutEthAndUsd{
				Volume24H: big.NewFloat(0),
				EthAmount: big.NewFloat(0),
				UsdAmount: big.NewFloat(0),
			}
			if _valueFloat.Cmp(big.NewFloat(0)) > 0 {
				// 直接用ETH代币购买代币,value 不为0
				_wethValue, _ := new(big.Float).SetString(_valueString)
				_swapAmoutEthAndUsd.EthAmount = new(big.Float).Add(_swapAmoutEthAndUsd.EthAmount, _wethValue)
			} else if _valueFloat.Cmp(big.NewFloat(0)) == 0 {
				_usdValue, _ := new(big.Float).SetString(_v.Amount0)
				_swapAmoutEthAndUsd.UsdAmount = new(big.Float).Add(_swapAmoutEthAndUsd.UsdAmount, _usdValue.Abs(_usdValue))
			}

			_oldValue, ok := _swapAmoutEthAndUsdMap[strings.ToLower(_v.ExtraPool)]
			if !ok {
				_swapAmoutEthAndUsdMap[strings.ToLower(_v.ExtraPool)] = _swapAmoutEthAndUsd
			} else {
				_swapAmoutEthAndUsdMap[strings.ToLower(_v.ExtraPool)].EthAmount = new(big.Float).Add(_oldValue.EthAmount, _swapAmoutEthAndUsd.EthAmount)
				_swapAmoutEthAndUsdMap[strings.ToLower(_v.ExtraPool)].UsdAmount = new(big.Float).Add(_oldValue.UsdAmount, _swapAmoutEthAndUsd.UsdAmount)
			}

		} else if strings.EqualFold(strings.ToLower(_v.ExtraToken0), strings.ToLower(chain.USD)) &&
			!strings.EqualFold(strings.ToLower(_v.ExtraToken1), strings.ToLower(chain.WETH)) ||
			strings.EqualFold(strings.ToLower(_v.ExtraToken1), strings.ToLower(chain.USD)) &&
				!strings.EqualFold(strings.ToLower(_v.ExtraToken0), strings.ToLower(chain.WETH)) {
			_swapAmoutEthAndUsd := &SwapAmoutEthAndUsd{
				Volume24H: big.NewFloat(0),
				EthAmount: big.NewFloat(0),
				UsdAmount: big.NewFloat(0),
			}
			// 2 ,USD/TOKENA 交易对
			// 此时_value 不为 0,因此不考虑 _value = "0"  的情况,
			_usdValue, _ := new(big.Float).SetString(_v.Amount1)
			_swapAmoutEthAndUsd.UsdAmount = new(big.Float).Add(_swapAmoutEthAndUsd.UsdAmount, _usdValue.Abs(_usdValue))

			_oldValue, ok := _swapAmoutEthAndUsdMap[strings.ToLower(_v.ExtraPool)]
			if !ok {
				_swapAmoutEthAndUsdMap[strings.ToLower(_v.ExtraPool)] = _swapAmoutEthAndUsd
			} else {
				_swapAmoutEthAndUsdMap[strings.ToLower(_v.ExtraPool)].UsdAmount = new(big.Float).Add(_oldValue.UsdAmount, _swapAmoutEthAndUsd.UsdAmount)
			}

		} else if strings.EqualFold(strings.ToLower(_v.ExtraToken0), strings.ToLower(chain.USD)) &&
			strings.EqualFold(strings.ToLower(_v.ExtraToken1), strings.ToLower(chain.WETH)) ||
			strings.EqualFold(strings.ToLower(_v.ExtraToken1), strings.ToLower(chain.USD)) &&
				strings.EqualFold(strings.ToLower(_v.ExtraToken0), strings.ToLower(chain.WETH)) {
			_swapAmoutEthAndUsd := &SwapAmoutEthAndUsd{
				Volume24H: big.NewFloat(0),
				EthAmount: big.NewFloat(0),
				UsdAmount: big.NewFloat(0),
			}
			// 3, USD/ETH 交易对
			if _valueFloat.Cmp(big.NewFloat(0)) > 0 {
				// 直接用ETH代币购买代币,value 不为0
				_wethValue, _ := new(big.Float).SetString(_valueString)
				_swapAmoutEthAndUsd.EthAmount = new(big.Float).Add(_swapAmoutEthAndUsd.EthAmount, _wethValue)
			} else if _valueFloat.Cmp(big.NewFloat(0)) == 0 {
				_usdValue, _ := new(big.Float).SetString(_v.Amount1)
				_swapAmoutEthAndUsd.UsdAmount = new(big.Float).Add(_swapAmoutEthAndUsd.UsdAmount, _usdValue.Abs(_usdValue))
			}

			_oldValue, ok := _swapAmoutEthAndUsdMap[strings.ToLower(_v.ExtraPool)]
			if !ok {
				_swapAmoutEthAndUsdMap[strings.ToLower(_v.ExtraPool)] = _swapAmoutEthAndUsd
			} else {
				_swapAmoutEthAndUsdMap[strings.ToLower(_v.ExtraPool)].EthAmount = new(big.Float).Add(_oldValue.EthAmount, _swapAmoutEthAndUsd.EthAmount)
				_swapAmoutEthAndUsdMap[strings.ToLower(_v.ExtraPool)].UsdAmount = new(big.Float).Add(_oldValue.UsdAmount, _swapAmoutEthAndUsd.UsdAmount)
			}
		}
	}

	prices, _ := tokenPriceModel.GetCurrentPrices([]string{chain.WETH, chain.USD})
	_wethPrice := ""
	_usdPrice := ""
	_wethPriceForUsdDecimals := 6 // 数据库存储的价格精度，默认是6
	_usdPriceForUsdDecimals := 6  // 数据库存储的价格精度，默认是6
	for _, v := range prices {
		if strings.EqualFold(strings.ToLower(v.TokenAddress), strings.ToLower(chain.WETH)) {
			_wethPrice = v.Price
			_wethPriceForUsdDecimals = v.UsdDecimals
		}

		if strings.EqualFold(strings.ToLower(v.TokenAddress), strings.ToLower(chain.USD)) {
			_usdPrice = v.Price
			_usdPriceForUsdDecimals = v.UsdDecimals
		}
	}

	// gettoken price
	if strings.EqualFold(_wethPrice, "") {
		// weth 默认值为 $3400
		// 3400*10**6 = 3313849586
		_wethPrice = "3313849586"
	}

	if strings.EqualFold(_usdPrice, "") {
		// _usdPrice 默认值为 $1
		// 10**6 = 1000000
		// usd 为 6 位小数
		_usdPrice = "1000000"
	}

	// get the data slice with processed pools

	_wethPriceFloat, _ := new(big.Float).SetString(_wethPrice)

	//  数据库里面拿出来的ETH的价格也要除以 10**6 ，除 6 是根据 数据库里面写入的ETH价格是  6位精度的 sting
	_wethPriceFloatQuo6pow10 := new(big.Float).Quo(_wethPriceFloat, new(big.Float).SetFloat64(math.Pow10(_wethPriceForUsdDecimals)))
	// _usdPriceFloat, _ := new(big.Float).SetString(_usdPrice)

	// assign value
	for i, v := range _swapAmoutEthAndUsdMap {
		// _swapAmoutEthAndUsdMap[i].Volume24H = new(big.Float).Add(new(big.Float).Mul(_wethPriceFloat, new(big.Float).Quo(v.EthAmount, new(big.Float).SetFloat64(math.Pow10(18)))), v.UsdAmount)
		// 数据库里面拿出来ethValue（也就是数据库里面value字段），要除以10**18，根据数据库里面保存的是 18位精度  event_sonex_pool_swap
		// usdAmount 要除以10**6，根据 event_sonex_pool_swap 表里面  evenet0_token 或者event1 的token， 是 6 位精度保存的，
		// 然后相加
		_ethAmountDiv18 := new(big.Float).Quo(v.EthAmount, new(big.Float).SetFloat64(math.Pow10(18)))
		_ethValue := new(big.Float).Mul(_wethPriceFloatQuo6pow10, _ethAmountDiv18)

		_usdAmountDiv6 := new(big.Float).Quo(v.UsdAmount, new(big.Float).SetFloat64(math.Pow10(_usdPriceForUsdDecimals)))

		_swapAmoutEthAndUsdMap[i].Volume24H = new(big.Float).Add(_ethValue, _usdAmountDiv6)

		logx.Debugf(" i:%v v:%#v _wethPriceFloatQuo6pow10:%v v.EthAmount:%v _ethAmountDiv18:%v ,_ethValue:%v, v.UsdAmount:%v _usdAmountDiv6:%v  	_swapAmoutEthAndUsdMap[i].Volume24H :%v  ", i, v, _wethPriceFloatQuo6pow10, v.EthAmount, _ethAmountDiv18, _ethValue, v.UsdAmount, _usdAmountDiv6, _swapAmoutEthAndUsdMap[i].Volume24H)
	}

	return _swapAmoutEthAndUsdMap, nil
}
