package sonex

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"sonex_api_backend_go/internal/lib/cache"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/module/sonex"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type SonexPositionsLogic struct {
	logx.Logger
	ctx                    context.Context
	svcCtx                 *svc.ServiceContext
	transferModel          repository.EventSonexPositionManagerTransferModel
	collectModel           repository.EventSonexPositionManagerCollectModel
	decreaseLiquidityModel repository.EventSonexPositionManagerDecreaseLiquidityModel
	tokenPriceModel        repository.SonexTokenPriceRecordModel
}

func NewSonexPositionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexPositionsLogic {
	conn := repository.NewConn(svcCtx.Config.PostgreSQL.DateSource())
	transferModel := repository.NewEventSonexPositionManagerTransferModel(conn)
	collectModel := repository.NewEventSonexPositionManagerCollectModel(conn)
	decreaseLiquidityModel := repository.NewEventSonexPositionManagerDecreaseLiquidityModel(conn)
	tokenPriceModel := repository.NewSonexTokenPriceRecordModel(conn)
	return &SonexPositionsLogic{
		Logger:                 logx.WithContext(ctx),
		ctx:                    ctx,
		svcCtx:                 svcCtx,
		transferModel:          transferModel,
		collectModel:           collectModel,
		decreaseLiquidityModel: decreaseLiquidityModel,
		tokenPriceModel:        tokenPriceModel,
	}
}

type PositionsResponseWrapper struct {
	*types.PositionsResponse
}

type psTaskResult struct {
	collectedPairs     map[int64]*sonex.TokenPair
	tokenPricesRecords map[string]*repository.TokenPriceRecord
	err                error
}

const (
	CacheKey = "sonex_dex:uniswap:positions:%s"
)

// DelCache 删除缓存
func (l *SonexPositionsLogic) DelCache(address string) error {
	address = common.HexToAddress(address).Hex()
	cacheKey := fmt.Sprintf(CacheKey, address)
	err := cache.NewRedisCache(fmt.Sprintf("%s:%d", l.svcCtx.Config.Redis.Host, l.svcCtx.Config.Redis.Port), l.svcCtx.Config.Redis.Password, l.svcCtx.Config.Redis.DB).Del(cacheKey)
	if err != nil {
		logx.Errorf("[SonexPositionsLogic] DelCache error: %s", err.Error())
		return err
	}
	return nil
}

func (l *SonexPositionsLogic) SonexPositions(req *types.PositionsRequest) (resp *PositionsResponseWrapper, err error) {
	cacheKey := fmt.Sprintf(CacheKey, req.Address)

	// Check cache
	redisCache := cache.NewRedisCache(fmt.Sprintf("%s:%d", l.svcCtx.Config.Redis.Host, l.svcCtx.Config.Redis.Port), l.svcCtx.Config.Redis.Password, l.svcCtx.Config.Redis.DB)
	cachedResponse, err := redisCache.Get(cacheKey)
	if err == nil && len(cachedResponse) > 0 {
		var cachedResp PositionsResponseWrapper
		if err := json.Unmarshal([]byte(cachedResponse), &cachedResp); err == nil {
			ctx := context.Background()
			go l.refreshCache(ctx, req, redisCache, cacheKey, &cachedResp)
			return &cachedResp, nil
		}
	}

	return l.refreshCache(l.ctx, req, redisCache, cacheKey, nil)
}

func (l *SonexPositionsLogic) refreshCache(ctx context.Context, req *types.PositionsRequest, cacher cache.ICache, cacheKey string, cachedResp *PositionsResponseWrapper) (resp *PositionsResponseWrapper, err error) {
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

	items, err := l.transferModel.GetTransferEvents(ctx, req.ChainId, req.Address, true)
	if err != nil {
		logx.Errorf("[SonexPositionsLogic] GetTransferEvents error: %s", err.Error())
		return nil, err
	}

	positionTokenIds := make([]int64, 0)
	positionTransferItemMap := make(map[int64]*types.Position)
	token0Decimals, token1Decimals := make([]int, 0), make([]int, 0)
	token0Addresses, token1Addresses := make([]string, 0), make([]string, 0)
	for _, item := range items {
		if !item.IsBurnt {
			tokenId := item.TokenID
			positionTokenIds = append(positionTokenIds, tokenId)
			positionTransferItemMap[tokenId] = item
		}
	}

	pv := NewSonexPositionLogic(l.ctx, l.svcCtx)
	tokenIds, err := pv.GetAddressTokenIds(ctx, req.Address)
	if err != nil {
		logx.Errorf("[SonexPositionsLogic] GetAddressTokenIds error: %s", err.Error())
	}

	caller, _ := sonex.NewStaticCaller(l.svcCtx.Config)
	positions, err := caller.BatchGetPositions(
		ctx,
		req.ChainId,
		common.HexToAddress(l.svcCtx.Config.GetSonexConfig(req.ChainId).PositionManager),
		func(tokenIds []int64) []*big.Int {
			// Create a map to track unique tokenIds
			uniqueIds := make(map[int64]bool)
			for _, id := range tokenIds {
				uniqueIds[id] = true
			}

			// Convert unique ids to slice
			converted := make([]*big.Int, 0, len(uniqueIds))
			for id := range uniqueIds {
				converted = append(converted, big.NewInt(id))
			}
			return converted
		}(append(positionTokenIds, tokenIds...)),
	)

	for _, p := range positions {
		token0, err := global.GetToken(req.ChainId, p.Token0.Hex())
		if err != nil {
			logx.Errorf("[SonexPositionsLogic] GetToken error: %s", err.Error())
			return nil, err
		}
		token1, err := global.GetToken(req.ChainId, p.Token1.Hex())
		if err != nil {
			logx.Errorf("[SonexPositionsLogic] GetToken error: %s", err.Error())
			return nil, err
		}

		token0Decimals = append(token0Decimals, int(token0.Decimals))
		token1Decimals = append(token1Decimals, int(token1.Decimals))
		token0Addresses = append(token0Addresses, token0.Address)
		token1Addresses = append(token1Addresses, token1.Address)

		if _, ok := positionTransferItemMap[p.TokenId.Int64()]; !ok {
			positionTransferItemMap[p.TokenId.Int64()] = &types.Position{
				TokenID: p.TokenId.Int64(),
				Origin:  p.Owner.Hex(),
				Pool: types.Token_PoolDetail{
					Address: p.Pool.Address.Hex(),
					FeeTier: int(p.Pool.Fee.Int64()),
					Token0: types.Token_TokenInfo{
						Address:  token0.Address,
						Decimals: int(token0.Decimals),
						Symbol:   token0.Symbol,
						Name:     token0.Name,
					},
					Token1: types.Token_TokenInfo{
						Address:  token1.Address,
						Decimals: int(token1.Decimals),
						Symbol:   token1.Symbol,
						Name:     token1.Name,
					},
					Tick: p.Pool.Slot0.Tick.String(),
				},
				TickUpper:   p.TickUpper.String(),
				TickLower:   p.TickLower.String(),
				PositionAge: 0,
			}
		}
	}

	var wg sync.WaitGroup
	results := make(chan psTaskResult, 2)

	// Task 1: Get collected fees
	wg.Add(1)
	go func() {
		defer wg.Done()
		pairs, err := l.getCollectedPairs(ctx, req.ChainId, positionTokenIds, token0Decimals, token1Decimals)
		results <- psTaskResult{collectedPairs: pairs, err: err}
	}()

	// Task 2: Get token prices
	wg.Add(1)
	go func() {
		defer wg.Done()
		pricesRecords, err := l.tokenPriceModel.GetCurrentPrices(append(token0Addresses, token1Addresses...))
		results <- psTaskResult{tokenPricesRecords: pricesRecords, err: err}
	}()

	// Close results channel after all goroutines complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var tokenPricesRecords = make(map[string]*repository.TokenPriceRecord)

	for result := range results {
		if result.err != nil {
			logx.Errorf("[SonexPositions] waitgroup error: %s", result.err.Error())
			return nil, result.err
		}
		if result.collectedPairs != nil {
			for tokenId, pair := range result.collectedPairs {
				position, exists := positionTransferItemMap[tokenId]
				if exists {
					position.TotalCollected = *pair.ToTokenValuePair()
				}
			}
		}
		if result.tokenPricesRecords != nil {
			tokenPricesRecords = result.tokenPricesRecords
		}
	}

	response := &PositionsResponseWrapper{
		PositionsResponse: &types.PositionsResponse{
			Positions: make([]*types.Position, 0),
		},
	}

	// Collect results
	for _, position := range positions {
		if positionTransferItem, exists := positionTransferItemMap[position.TokenId.Int64()]; exists {
			convertedPosition, err := l.fillPosition(ctx, req, positionTransferItem, position, tokenPricesRecords)
			if err != nil {
				logx.Errorf("[SonexPositions]:refreshCache fillPosition error: %s", err.Error())
				continue
			}
			response.Positions = append(response.Positions, convertedPosition)
		}
	}

	// sort by position in DESC
	sort.Slice(response.Positions, func(i, j int) bool {
		return response.Positions[i].Position > response.Positions[j].Position
	})

	responseBytes, err := json.Marshal(response)
	if err == nil {
		cacher.Set(cacheKey, responseBytes, 0) // TODO: cache for xxx seconds??
	}

	return response, nil
}

// 获取 collectedPairs，过滤掉 deceaseLiquidity events
func (l *SonexPositionsLogic) getCollectedPairs(ctx context.Context, chainId int64, positionTokenIds []int64, token0Decimals, token1Decimals []int) (map[int64]*sonex.TokenPair, error) {
	collectEventsMap, err := l.collectModel.GetCollectEvents(ctx, chainId, positionTokenIds)
	if err != nil {
		logx.Errorf("[SonexPosition] GetCollectEvents error: %s", err.Error())
		return nil, err
	}
	decreaseLiquidityEventsMap, err := l.decreaseLiquidityModel.GetDecreaseLiquidityEvents(ctx, chainId, positionTokenIds)
	if err != nil {
		logx.Errorf("[SonexPosition] GetDecreaseLiquidityEvents error: %s", err.Error())
		return nil, err
	}
	collectedPairs := make(map[int64]*sonex.TokenPair)
	for i, tokenId := range positionTokenIds {
		// 1. initial token pair
		tokenValue0, err := sonex.NewTokenValue("0", token0Decimals[i])
		if err != nil {
			logx.Errorf("[SonexPosition] NewTokenValue error: %s", err.Error())
			return nil, err
		}
		tokenValue1, err := sonex.NewTokenValue("0", token1Decimals[i])
		if err != nil {
			logx.Errorf("[SonexPosition] NewTokenValue error: %s", err.Error())
			return nil, err
		}
		collectedPairs[tokenId] = sonex.NewTokenPair(tokenValue0, tokenValue1)

		// 2. collect events
		collectEvents, ok := collectEventsMap[tokenId]
		if !ok {
			continue
		}

		// 3. decrease liquidity events
		decreEvents, _ := decreaseLiquidityEventsMap[tokenId]

		// 4. count collected pairs and filter decrease liquidity events
		for _, collectEvent := range collectEvents {
			isDecreaseLiquidityEvent := false
			for _, decreEvent := range decreEvents {
				if collectEvent.TxHash == decreEvent.TxHash {
					// this collect event is decrease liquidity event, do not count it.
					isDecreaseLiquidityEvent = true
					break
				}
			}
			if isDecreaseLiquidityEvent {
				continue
			}
			tokenValue0 = tokenValue0.AmountAdd(collectEvent.EventAmount0)
			tokenValue1 = tokenValue1.AmountAdd(collectEvent.EventAmount1)
		}
		collectedPairs[tokenId] = sonex.NewTokenPair(tokenValue0, tokenValue1)
	}

	return collectedPairs, nil
}

func (l *SonexPositionsLogic) fillPosition(ctx context.Context,
	req *types.PositionsRequest,
	position *types.Position,
	innnerPosition *sonex.Position,
	pricesRecords map[string]*repository.TokenPriceRecord) (*types.Position, error) {
	price0 := new(big.Float).SetFloat64(0)
	price1 := new(big.Float).SetFloat64(0)

	chainConf := l.svcCtx.Config.GetChainConfig(req.ChainId)
	if chainConf == nil {
		return nil, fmt.Errorf("Miss chain configure")
	}
	usdToken, err := global.GetToken(req.ChainId, common.HexToAddress(chainConf.USD).Hex())
	if err != nil {
		return nil, err
	}

	usdDecimals := int(usdToken.Decimals)
	factorUsd := new(big.Float).SetFloat64(math.Pow10(usdDecimals))

	price0String := ""
	price1String := ""
	// 优先获取从 db 中的价格，如果获取失败，则从链上获取
	if price0Record, exists0 := pricesRecords[innnerPosition.Token0.Hex()]; exists0 {
		price0, _ = price0.SetString(price0Record.Price)
		price0String = price0Record.Price
	} else {
		logx.Errorf("[convertPosition] Token0 price not found from db for address: %s", innnerPosition.Token0.Hex())

		bestPrice0, _, _, err := global.GetPrices(ctx, nil, req.ChainId, innnerPosition.Token0.Hex(), true, 2)
		if err == nil && bestPrice0 != nil {
			price0.SetInt(bestPrice0.Value)
			price0String = bestPrice0.Value.String()
		}
	}

	if price1Record, exists1 := pricesRecords[innnerPosition.Token1.Hex()]; exists1 {
		price1, _ = price1.SetString(price1Record.Price)
		price1String = price1Record.Price
	} else {
		logx.Errorf("[convertPosition] Token1 price not found from db for address: %s", innnerPosition.Token1.Hex())

		bestPrice1, _, _, err := global.GetPrices(ctx, nil, req.ChainId, innnerPosition.Token1.Hex(), true, 2)
		if err == nil && bestPrice1 != nil {
			price1.SetInt(bestPrice1.Value)
			price1String = bestPrice1.Value.String()
		}
	}

	poolPrice0 := new(big.Float).Quo(price0, factorUsd)
	position.Pool.Token0Price = poolPrice0.String()

	// 转换 price 为 *big.Float 类型
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
	// 4. PNL(Fee Revenue) = claimableFees + collectedFees
	pnl := position.ClaimableFees
	collectedPair := position.TotalCollected
	token0Value, err = sonex.NewTokenValue(collectedPair.Token0.Amount, position.Pool.Token0.Decimals, sonex.WithPrice(price0String, usdDecimals))
	if err != nil {
		logx.Errorf("[SonexPosition] NewTokenValue error: %s", err.Error())
		return nil, err
	}
	token1Value, err = sonex.NewTokenValue(collectedPair.Token1.Amount, position.Pool.Token1.Decimals, sonex.WithPrice(price1String, usdDecimals))
	if err != nil {
		logx.Errorf("[SonexPosition] NewTokenValue error: %s", err.Error())
		return nil, err
	}
	tokenPair = sonex.NewTokenPair(token0Value, token1Value)
	collectedFees, _ := tokenPair.ValueFloat64()
	position.PNL = pnl + collectedFees

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
