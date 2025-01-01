package sonex

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"sync"
	"time"

	"sonex_api_backend_go/internal/lib/cache"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/module/sonex"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type SonexPositionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexPositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexPositionLogic {
	return &SonexPositionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type PositionResponseWrapper struct {
	*types.PositionResponse
}

type sonexPositionResult struct {
	transferItem  *types.Position
	collectedItem *types.AggregatedAmounts
	position      *sonex.Position
	prices        map[string]*repository.TokenPriceRecord
	err           error
}

func (l *SonexPositionLogic) SonexPosition(req *types.PositionRequest) (resp *PositionResponseWrapper, err error) {
	cacheKey := fmt.Sprintf("sonex_dex:uniswap:position:%d", req.TokenId)
	redisCache := cache.NewRedisCache(fmt.Sprintf("%s:%d", l.svcCtx.Config.Redis.Host, l.svcCtx.Config.Redis.Port), l.svcCtx.Config.Redis.Password, l.svcCtx.Config.Redis.DB)
	cachedResponse, err := redisCache.Get(cacheKey)
	if err == nil && len(cachedResponse) > 0 {
		var cachedResp PositionResponseWrapper
		if err := json.Unmarshal([]byte(cachedResponse), &cachedResp); err == nil {
			// Launch a goroutine to refresh the cache in the background
			go l.refreshCache(req, redisCache, cacheKey, &cachedResp)
			// Return the cached response immediately
			return &cachedResp, nil
		}
	}

	return l.refreshCache(req, redisCache, cacheKey, nil)
}

func (l *SonexPositionLogic) refreshCache(req *types.PositionRequest, cacher cache.ICache, cacheKey string, cachedResp *PositionResponseWrapper) (resp *PositionResponseWrapper, err error) {
	lockKey := fmt.Sprintf("lock:%s", cacheKey)

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
	collectModel := repository.NewEventSonexPositionManagerCollectModel(conn)

	var wg1 sync.WaitGroup
	results1 := make(chan sonexPositionResult, 2)

	// Task 1: Get transfer event
	wg1.Add(1)
	go func() {
		defer wg1.Done()
		transferModel := repository.NewEventSonexPositionManagerTransferModel(conn)
		item, err := transferModel.GetTransferEvent(l.ctx, req.TokenId)
		results1 <- sonexPositionResult{transferItem: item, err: err}
	}()

	// Task 2: Get positions
	wg1.Add(1)
	go func() {
		defer wg1.Done()
		caller, _ := sonex.NewStaticCaller(l.svcCtx.Config)
		positions, err := caller.BatchGetPositions(
			l.ctx,
			req.ChainId,
			common.HexToAddress(l.svcCtx.Config.GetSonexConfig(req.ChainId).PositionManager),
			[]*big.Int{big.NewInt(req.TokenId)},
		)
		position := &sonex.Position{}
		if len(positions) > 0 {
			position = positions[0]
		}
		results1 <- sonexPositionResult{position: position, err: err}
	}()

	// Close results channel after all goroutines complete
	go func() {
		wg1.Wait()
		close(results1)
	}()

	// Collect results
	var position = &types.Position{}
	var collectedItem = &types.AggregatedAmounts{}
	var innerPosition = &sonex.Position{}
	var prices = make(map[string]*repository.TokenPriceRecord)

	for result := range results1 {
		if result.err != nil {
			logx.Errorf("[SonexPosition] waitgroup error: %s", result.err.Error())
			return nil, result.err
		}
		if result.transferItem != nil {
			position = result.transferItem
		}
		if result.position != nil {
			innerPosition = result.position
		}
	}

	var wg2 sync.WaitGroup
	results2 := make(chan sonexPositionResult, 2)

	// Task 1: Get collected amounts
	wg1.Add(1)
	go func() {
		defer wg1.Done()
		collectedItem, err := collectModel.AggregateAmounts(l.ctx, req.ChainId, req.TokenId,
			position.Pool.Token0.Decimals,
			position.Pool.Token1.Decimals)
		results2 <- sonexPositionResult{collectedItem: collectedItem, err: err}
	}()

	// Task 4: Get current prices
	wg1.Add(1)
	go func() {
		defer wg1.Done()
		prices, err := tokenPriceModel.GetCurrentPrices([]string{innerPosition.Token0.Hex(), innerPosition.Token1.Hex()})
		if err != nil {
			logx.Errorf("[SonexPosition] GetCurrentPrices error: %s", err.Error())
			results2 <- sonexPositionResult{prices: nil, err: err}
			return
		}
		results2 <- sonexPositionResult{prices: prices, err: nil}
	}()

	go func() {
		wg2.Wait()
		close(results2)
	}()

	for result := range results2 {
		if result.err != nil {
			logx.Errorf("[SonexPosition] waitgroup error: %s", result.err.Error())
			return nil, result.err
		}
		if result.prices != nil {
			prices = result.prices
		}
	}

	position = l.fillPosition(l.ctx, req, position, innerPosition, prices, collectedItem)

	// Continue with the rest of the logic using transferItem, collectedItem, and positions
	response := &PositionResponseWrapper{
		PositionResponse: &types.PositionResponse{
			Position: position,
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

type priceResult struct {
	index int
	price *big.Float
	err   error
}

func (l *SonexPositionLogic) fillPosition(ctx context.Context,
	req *types.PositionRequest,
	position *types.Position,
	innnerPosition *sonex.Position,
	pricesRecords map[string]*repository.TokenPriceRecord,
	collectedAmounts *types.AggregatedAmounts) *types.Position {
	position.TokenID = innnerPosition.TokenId.Int64()
	position.TickUpper = innnerPosition.TickUpper.String()
	position.TickLower = innnerPosition.TickLower.String()
	factorDai := new(big.Float).SetFloat64(math.Pow10(18))
	prices := []*big.Float{
		new(big.Float).SetFloat64(0), // price0
		new(big.Float).SetFloat64(0), // price1
	}
	tokens := []common.Address{
		innnerPosition.Token0,
		innnerPosition.Token1,
	}

	var wg sync.WaitGroup
	resultChan := make(chan priceResult, len(tokens))

	// Fetch prices in parallel
	for i, token := range tokens {
		wg.Add(1)
		go func(index int, tokenAddr common.Address) {
			defer wg.Done()

			if priceRecord, exists := pricesRecords[tokenAddr.Hex()]; exists {
				if p, ok := new(big.Float).SetString(priceRecord.Price); ok {
					resultChan <- priceResult{index: index, price: p}
					return
				}
			}

			logx.Errorf("[convertPosition] Token price not found from db for address: %s", tokenAddr.Hex())
			bestPrice, _, err := global.GetPrices(ctx, nil, req.ChainId, tokenAddr.Hex(), true, 2, nil)
			if err == nil && bestPrice != nil {
				p := new(big.Float).SetInt(bestPrice.Value)
				resultChan <- priceResult{index: index, price: p}
				return
			}
			resultChan <- priceResult{
				index: index,
				err:   fmt.Errorf("failed to get price for token%d", index),
			}
		}(i, token)
	}

	// Close channel after all goroutines complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	for result := range resultChan {
		if result.err == nil {
			prices[result.index] = result.price
		}
	}

	price0, price1 := prices[0], prices[1]

	factor0 := new(big.Float).SetFloat64(math.Pow10(int(position.Pool.Token0.Decimals)))
	factor1 := new(big.Float).SetFloat64(math.Pow10(int(position.Pool.Token1.Decimals)))

	poolPrice0 := new(big.Float).Quo(price0, factorDai)
	position.Price0, _ = poolPrice0.Float64()
	position.Pool.Token0Price = poolPrice0.String()

	// 转换 price 为 *big.Float 类型，并根据 factor0 和 factor1 做调整
	poolPrice1 := new(big.Float).Quo(price1, factorDai)
	position.Price1, _ = poolPrice1.Float64()
	position.Pool.Token1Price = poolPrice1.String()

	// 获取 amount0, amount1 和 unclaimedFee0, unclaimedFee1
	amount0, amount1 := innnerPosition.Amounts()
	unclaimedFee0, unclaimedFee1 := innnerPosition.UnclaimedFees()

	// 转换 amount0, amount1 为 *big.Float 类型，并根据 factor0 和 factor1 做调整
	amount0Float := new(big.Float).SetInt(amount0)
	amount1Float := new(big.Float).SetInt(amount1)

	amount0Float.Quo(amount0Float, factor0)
	amount1Float.Quo(amount1Float, factor1)

	// 计算 amount0 和 amount1 的余额
	positionValue0 := new(big.Float).Mul(amount0Float, poolPrice0)
	positinoValue1 := new(big.Float).Mul(amount1Float, poolPrice1)

	// 1. 计算 amount0 和 amount1
	position.Amount0, _ = amount0Float.Float64()
	position.Amount1, _ = amount1Float.Float64()

	// 2. 计算总 position
	totalPosition := new(big.Float).Add(positionValue0, positinoValue1)
	position.Position, _ = totalPosition.Float64()

	// 计算未结费用 unclaimedFee0 和 unclaimedFee1
	position.UnclaimedFee0, _ = new(big.Float).SetInt(unclaimedFee0).Float64()
	position.UnclaimedFee1, _ = new(big.Float).SetInt(unclaimedFee1).Float64()

	// 计算 claimable fees
	unclaimedFee0Float := new(big.Float).SetInt(unclaimedFee0)
	unclaimedFee1Float := new(big.Float).SetInt(unclaimedFee1)
	unclaimedFee0Float = unclaimedFee0Float.Quo(unclaimedFee0Float, factor0)
	unclaimedFee1Float = unclaimedFee1Float.Quo(unclaimedFee1Float, factor1)

	claimableFee0 := new(big.Float).Mul(unclaimedFee0Float, poolPrice0)
	claimableFee1 := new(big.Float).Mul(unclaimedFee1Float, poolPrice1)

	// 3. claimableFees
	claimableFees := new(big.Float).Add(claimableFee0, claimableFee1)
	position.ClaimableFees, _ = claimableFees.Float64()
	position.UnclaimedFee0, _ = unclaimedFee0Float.Float64()
	position.UnclaimedFee1, _ = unclaimedFee1Float.Float64()

	// 4. PNL = claimableFees + collectedFees
	position.PNL = position.ClaimableFees + collectedAmounts.TotalAmount0 + collectedAmounts.TotalAmount1

	return position
}
