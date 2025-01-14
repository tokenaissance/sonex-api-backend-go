package sonex

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"sonex_api_backend_go/internal/lib/cache"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/module/sonex"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"
	"strconv"
	"sync"
	"time"

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
	transferItem   *types.Position
	collectedPairs map[int64]*sonex.TokenPair
	position       *sonex.Position
	prices         map[string]*repository.TokenPriceRecord
	totalInvested  *sonex.TokenPair
	totalWithdrawn *sonex.TokenPair
	totalCollected *sonex.TokenPair
	err            error
}

func (l *SonexPositionLogic) SonexPosition(req *types.PositionRequest) (resp *PositionResponseWrapper, err error) {
	cacheKey := fmt.Sprintf("sonex_dex:uniswap:position:%d", req.TokenId)
	redisCache := cache.NewRedisCache(fmt.Sprintf("%s:%d", l.svcCtx.Config.Redis.Host, l.svcCtx.Config.Redis.Port), l.svcCtx.Config.Redis.Password, l.svcCtx.Config.Redis.DB)
	cachedResponse, err := redisCache.Get(cacheKey)
	if err == nil && len(cachedResponse) > 0 {
		var cachedResp PositionResponseWrapper
		if err := json.Unmarshal([]byte(cachedResponse), &cachedResp); err == nil {
			// Launch a goroutine to refresh the cache in the background
			ctx := context.Background()
			go l.refreshCache(ctx, req, redisCache, cacheKey, &cachedResp)
			// Return the cached response immediately
			return &cachedResp, nil
		}
	}

	return l.refreshCache(l.ctx, req, redisCache, cacheKey, nil)
}

func (l *SonexPositionLogic) refreshCache(ctx context.Context, req *types.PositionRequest, cacher cache.ICache, cacheKey string, cachedResp *PositionResponseWrapper) (resp *PositionResponseWrapper, err error) {
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

	// Continue with the rest of the logic using transferItem, collectedItem, and positions
	response := &PositionResponseWrapper{
		PositionResponse: &types.PositionResponse{},
	}

	defer func() {
		if lockAcquired {
			cacher.Del(lockKey)
		}
	}()

	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	tokenPriceModel := repository.NewSonexTokenPriceRecordModel(conn)
	collectModel := repository.NewEventSonexPositionManagerCollectModel(conn)
	increaseLiquidityModel := repository.NewEventSonexPositionManagerIncreaseLiquidityModel(conn)
	decreaseLiquidityModel := repository.NewEventSonexPositionManagerDecreaseLiquidityModel(conn)

	// Task 1: Get transfer event
	transferModel := repository.NewEventSonexPositionManagerTransferModel(conn)
	itemM, err := transferModel.GetTransferEvent(ctx, req.ChainId, req.TokenId, true)
	if err != nil {
		logx.Errorf("[SonexPosition] GetTransferEvent error: %s", err.Error())
		return nil, err
	}

	// if no transfer event and position NFT is burnt, return nil
	position, ok := itemM[req.TokenId]
	// if no transfer event, construct empty position
	if !ok {
		position = &types.Position{TokenID: req.TokenId, Pool: types.Token_PoolDetail{}}
		response.Position = position
	} else if position.IsBurnt {
		response.Position = &types.Position{IsBurnt: true}
		responseBytes, err := json.Marshal(response)
		if err == nil {
			err = cacher.Set(cacheKey, responseBytes, 0)
			if err != nil {
				l.Error("Failed to cache response", err)
			}
		}
		return response, nil
	}

	// Task 2: Get positions
	caller, _ := sonex.NewStaticCaller(l.svcCtx.Config)
	positions, err := caller.BatchGetPositions(
		ctx,
		req.ChainId,
		common.HexToAddress(l.svcCtx.Config.GetSonexConfig(req.ChainId).PositionManager),
		[]*big.Int{big.NewInt(req.TokenId)},
	)
	if err != nil {
		logx.Errorf("[SonexPosition] BatchGetPositions error: %s", err.Error())
		return nil, err
	}
	if len(positions) == 0 {
		logx.Errorf("[SonexPosition] BatchGetPositions error: %s", err.Error())
		return nil, err
	}
	innerPosition := positions[0]

	// cache tokenid
	err = l.cacheAddressTokenIds(ctx, innerPosition.Owner.Hex(), req.TokenId)
	if err != nil {
		logx.Errorf("[SonexPosition] cacheAddressTokenIds error: %s", err.Error())
	}

	// del cache
	if !ok {
		ps := NewSonexPositionsLogic(ctx, l.svcCtx)
		err = ps.DelCache(innerPosition.Owner.Hex())
		if err != nil {
			logx.Errorf("[SonexPosition] NewSonexPositionsLogic DelCache error: %s", err.Error())
		}
	}

	token0, err := global.GetToken(req.ChainId, innerPosition.Token0.Hex())
	if err != nil {
		logx.Errorf("[SonexPosition] GetToken error: %s", err.Error())
		return nil, err
	}
	token1, err := global.GetToken(req.ChainId, innerPosition.Token1.Hex())
	if err != nil {
		logx.Errorf("[SonexPosition] GetToken error: %s", err.Error())
		return nil, err
	}

	position.Origin = innerPosition.Owner.Hex()
	position.Pool.FeeTier = int(innerPosition.Pool.Fee.Int64())
	position.Pool.Address = innerPosition.Pool.Address.Hex()
	position.Pool.Token0 = types.Token_TokenInfo{
		Address:  token0.Address,
		Decimals: int(token0.Decimals),
		Symbol:   token0.Symbol,
		Name:     token0.Name,
	}
	position.Pool.Token1 = types.Token_TokenInfo{
		Address:  token1.Address,
		Decimals: int(token1.Decimals),
		Symbol:   token1.Symbol,
		Name:     token1.Name,
	}

	// Collect results
	var collectedPairs = make(map[int64]*sonex.TokenPair)
	var prices = make(map[string]*repository.TokenPriceRecord)

	var wg sync.WaitGroup
	results := make(chan sonexPositionResult, 2)
	// Task 1: Get collected amounts
	wg.Add(1)
	go func() {
		defer wg.Done()
		collectedPairs, err := collectModel.AggregateAmounts(ctx, req.ChainId, req.TokenId,
			position.Pool.Token0.Decimals,
			position.Pool.Token1.Decimals)
		results <- sonexPositionResult{collectedPairs: collectedPairs, err: err}
	}()
	// Task 2: Get total collected amounts
	wg.Add(1)
	go func() {
		defer wg.Done()
		totalCollected, err := collectModel.TotalAmounts(ctx, req.ChainId, req.TokenId,
			position.Pool.Token0.Decimals,
			position.Pool.Token1.Decimals)
		results <- sonexPositionResult{totalCollected: totalCollected, err: err}
	}()
	// Task 3: Get total invested amounts
	wg.Add(1)
	go func() {
		defer wg.Done()
		totalInvested, err := increaseLiquidityModel.TotalAmounts(ctx, req.ChainId, req.TokenId,
			position.Pool.Token0.Decimals,
			position.Pool.Token1.Decimals)
		results <- sonexPositionResult{totalInvested: totalInvested, err: err}
	}()
	// Task 4: Get total withdrawn amounts
	wg.Add(1)
	go func() {
		defer wg.Done()
		totalWithdrawn, err := decreaseLiquidityModel.TotalAmounts(ctx, req.ChainId, req.TokenId,
			position.Pool.Token0.Decimals,
			position.Pool.Token1.Decimals)
		results <- sonexPositionResult{totalWithdrawn: totalWithdrawn, err: err}
	}()
	// Task 5: Get current prices
	wg.Add(1)
	go func() {
		defer wg.Done()
		prices, err := tokenPriceModel.GetCurrentPrices([]string{innerPosition.Token0.Hex(), innerPosition.Token1.Hex()})
		if err != nil {
			logx.Errorf("[SonexPosition] GetCurrentPrices error: %s", err.Error())
			results <- sonexPositionResult{prices: nil, err: err}
			return
		}
		results <- sonexPositionResult{prices: prices, err: nil}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		if result.err != nil {
			logx.Errorf("[SonexPosition] waitgroup error: %s", result.err.Error())
			return nil, result.err
		}
		if result.prices != nil {
			prices = result.prices
		}
		if result.collectedPairs != nil {
			collectedPairs = result.collectedPairs
		}
		if result.totalInvested != nil {
			position.TotalInvested = *result.totalInvested.ToTokenValuePair()
		}
		if result.totalWithdrawn != nil {
			position.TotalWithdrawn = *result.totalWithdrawn.ToTokenValuePair()
		}
		if result.totalCollected != nil {
			position.TotalCollected = *result.totalCollected.ToTokenValuePair()
		}
	}

	position, err = l.fillPosition(ctx, req, position, innerPosition, prices, collectedPairs)
	if err != nil {
		logx.Errorf("[SonexPosition] fillPosition error: %s", err.Error())
		return nil, err
	}

	response.Position = position

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
	index       int
	price       *big.Float
	priceString string
	err         error
}

func (l *SonexPositionLogic) fillPosition(
	ctx context.Context,
	req *types.PositionRequest,
	position *types.Position,
	innnerPosition *sonex.Position,
	pricesRecords map[string]*repository.TokenPriceRecord,
	collectedPairs map[int64]*sonex.TokenPair,
) (*types.Position, error) {
	position.TokenID = innnerPosition.TokenId.Int64()
	position.TickUpper = innnerPosition.TickUpper.String()
	position.TickLower = innnerPosition.TickLower.String()

	chainConf := l.svcCtx.Config.GetChainConfig(req.ChainId)
	if chainConf == nil {
		return nil, errors.New("Miss chain configure")
	}
	usdToken, err := global.GetToken(req.ChainId, common.HexToAddress(chainConf.USD).Hex())
	if err != nil {
		return nil, err
	}

	usdDecimals := int(usdToken.Decimals)
	factorUsd := new(big.Float).SetFloat64(math.Pow10(usdDecimals))
	prices := []*big.Float{
		new(big.Float).SetFloat64(0), // price0
		new(big.Float).SetFloat64(0), // price1
	}
	tokens := []common.Address{
		innnerPosition.Token0,
		innnerPosition.Token1,
	}
	priceStrings := []string{
		"",
		"",
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
					resultChan <- priceResult{index: index, price: p, priceString: priceRecord.Price}
					return
				}
			}

			logx.Errorf("[convertPosition] Token price not found from db for address: %s", tokenAddr.Hex())
			bestPrice, _, _, err := global.GetPrices(ctx, nil, req.ChainId, tokenAddr.Hex(), true, 2)
			if err == nil && bestPrice != nil {
				p := new(big.Float).SetInt(bestPrice.Value)
				resultChan <- priceResult{index: index, price: p, priceString: bestPrice.Value.String()}
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
			priceStrings[result.index] = result.priceString
		}
	}

	price0, price1 := prices[0], prices[1]
	price0String, price1String := priceStrings[0], priceStrings[1]

	poolPrice0 := new(big.Float).Quo(price0, factorUsd)
	position.Pool.Token0Price = poolPrice0.String()

	// 转换 price 为 *big.Float 类型，并根据 factor0 和 factor1 做调整
	poolPrice1 := new(big.Float).Quo(price1, factorUsd)
	position.Pool.Token1Price = poolPrice1.String()

	// 获取 amount0, amount1 和 unclaimedFee0, unclaimedFee1
	amount0, amount1 := innnerPosition.Amounts()
	unclaimedFee0, unclaimedFee1 := innnerPosition.UnclaimedFees()

	// 2. 计算总 position
	amount0String := amount0.String()
	amount1String := amount1.String()
	token0Value, _ := sonex.NewTokenValue(amount0String, position.Pool.Token0.Decimals, sonex.WithPrice(price0String, usdDecimals))
	token1Value, _ := sonex.NewTokenValue(amount1String, position.Pool.Token1.Decimals, sonex.WithPrice(price1String, usdDecimals))
	tokenPair := sonex.NewTokenPair(token0Value, token1Value)
	position.Position, _ = tokenPair.ValueFloat64()

	// 计算未结费用 unclaimedFee0 和 unclaimedFee1
	position.UnclaimedFee0, _ = unclaimedFee0.Float64()
	position.UnclaimedFee1, _ = unclaimedFee1.Float64()

	// 3. claimableFees
	claimableToken0Value, _ := sonex.NewTokenValue(unclaimedFee0.String(), position.Pool.Token0.Decimals, sonex.WithPrice(price0String, usdDecimals))
	claimableToken1Value, _ := sonex.NewTokenValue(unclaimedFee1.String(), position.Pool.Token1.Decimals, sonex.WithPrice(price1String, usdDecimals))
	claimableTokenPair := sonex.NewTokenPair(claimableToken0Value, claimableToken1Value)
	position.ClaimableFees, _ = claimableTokenPair.ValueFloat64()

	position.UnclaimedFee0, _ = claimableToken0Value.ValueFloat64()
	position.UnclaimedFee1, _ = claimableToken1Value.ValueFloat64()

	pnl := position.ClaimableFees
	if collectedPair, ok := collectedPairs[req.TokenId]; ok {
		// set collected token price
		collectedPair.Token0().SetTokenPrice(price0String, usdDecimals)
		collectedPair.Token1().SetTokenPrice(price1String, usdDecimals)
		collectedFees, _ := collectedPair.ValueFloat64()
		pnl = pnl + collectedFees
	}
	position.PNL = pnl
	// 5. fill tokenValue0, tokenValue1
	position.TokenValue0 = types.Token_TokenValue{
		Amount:   amount0.String(),
		Decimals: position.Pool.Token0.Decimals,
		Token_TokenPrice: &types.Token_TokenPrice{
			Price:    price0String,
			Decimals: usdDecimals,
		},
	}
	position.TokenValue1 = types.Token_TokenValue{
		Amount:   amount1.String(),
		Decimals: position.Pool.Token1.Decimals,
		Token_TokenPrice: &types.Token_TokenPrice{
			Price:    price1String,
			Decimals: usdDecimals,
		},
	}

	// 6. fill tick
	position.Pool.Tick = innnerPosition.Pool.Slot0.Tick.String()

	return position, nil
}

func (l *SonexPositionLogic) cacheAddressTokenIds(ctx context.Context, address string, tokenid int64) error {
	address = common.HexToAddress(address).Hex()
	cacheKey := fmt.Sprintf("sonex_dex:uniswap:position:tokenids:%s", address)
	cacher := cache.NewRedisCache(fmt.Sprintf("%s:%d", l.svcCtx.Config.Redis.Host, l.svcCtx.Config.Redis.Port), l.svcCtx.Config.Redis.Password, l.svcCtx.Config.Redis.DB)

	// add tokenid to Set
	err := cacher.SAdd(ctx, cacheKey, tokenid).Err()
	if err != nil {
		return err
	}

	// set expire time
	err = cacher.Expire(cacheKey, 30*time.Second)
	if err != nil {
		return err
	}

	return nil
}

func (l *SonexPositionLogic) GetAddressTokenIds(ctx context.Context, address string) ([]int64, error) {
	address = common.HexToAddress(address).Hex()
	cacheKey := fmt.Sprintf("sonex_dex:uniswap:position:tokenids:%s", address)
	cacher := cache.NewRedisCache(fmt.Sprintf("%s:%d", l.svcCtx.Config.Redis.Host, l.svcCtx.Config.Redis.Port), l.svcCtx.Config.Redis.Password, l.svcCtx.Config.Redis.DB)

	result, err := cacher.SMembers(ctx, cacheKey).Result()
	if err != nil {
		return nil, err
	}

	// 将字符串转换为 int64
	tokenIds := make([]int64, 0, len(result))
	for _, str := range result {
		tokenId, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse tokenId: %w", err)
		}
		tokenIds = append(tokenIds, tokenId)
	}

	return tokenIds, nil
}
