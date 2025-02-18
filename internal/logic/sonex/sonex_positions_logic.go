package sonex

import (
	"context"
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
	increaseLiquidityModel repository.EventSonexPositionManagerIncreaseLiquidityModel
	tokenPriceModel        repository.SonexTokenPriceRecordModel
}

func NewSonexPositionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexPositionsLogic {
	conn := repository.NewConn(svcCtx.Config.PostgreSQL.DateSource())
	return &SonexPositionsLogic{
		Logger:                 logx.WithContext(ctx),
		ctx:                    ctx,
		svcCtx:                 svcCtx,
		transferModel:          repository.NewEventSonexPositionManagerTransferModel(conn),
		collectModel:           repository.NewEventSonexPositionManagerCollectModel(conn),
		decreaseLiquidityModel: repository.NewEventSonexPositionManagerDecreaseLiquidityModel(conn),
		increaseLiquidityModel: repository.NewEventSonexPositionManagerIncreaseLiquidityModel(conn),
		tokenPriceModel:        repository.NewSonexTokenPriceRecordModel(conn),
	}
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
		logx.Errorf("[SonexPositions]:DelCache error: %s", err.Error())
		return err
	}
	return nil
}

// SonexPositions get live position list
func (l *SonexPositionsLogic) SonexPositions(ctx context.Context, chainId int64, address string) (positions []*types.Position, err error) {
	items, err := l.transferModel.GetTransferEvents(ctx, chainId, address, true)
	if err != nil {
		logx.Errorf("[SonexPositions]: GetTransferEvents error: %s", err.Error())
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
	tokenIds, err := pv.GetAddressTokenIds(ctx, address)
	if err != nil {
		logx.Errorf("[SonexPositions]: GetAddressTokenIds error: %s", err.Error())
	}

	// merge positionTokenIds and tokenIds, and remove duplicate
	positionTokenIds = func(tokenIds []int64) []int64 {
		// Create a map to track unique tokenIds
		uniqueIds := make(map[int64]bool)
		for _, id := range tokenIds {
			uniqueIds[id] = true
		}

		// Convert unique ids to slice
		filtered := make([]int64, 0, len(uniqueIds))
		for id := range uniqueIds {
			filtered = append(filtered, id)
		}
		return filtered
	}(append(positionTokenIds, tokenIds...))

	caller, _ := sonex.NewStaticCaller(l.svcCtx.Config)
	chainConf := l.svcCtx.Config.GetSonexConfig(chainId)
	if chainConf == nil {
		logx.Errorf("[SonexPositions]: Miss chain configure")
		return nil, fmt.Errorf("%s", "Miss chain configure")
	}

	positionM, err := caller.BatchGetPositionMap(
		ctx,
		chainId,
		common.HexToAddress(chainConf.PositionManager),
		func(tokenIds []int64) []*big.Int {
			// Convert unique ids to bitInt slice
			converted := make([]*big.Int, 0, len(tokenIds))
			for _, id := range tokenIds {
				converted = append(converted, big.NewInt(id))
			}
			return converted
		}(positionTokenIds),
	)
	if err != nil {
		logx.Errorf("[SonexPositions]: BatchGetPositionMap error: %s", err.Error())
		return nil, err
	}

	// positionTokenIds 中部分 tokenId 可能被销毁，但是 event 监听未及时处理
	// 所以 BatchGetPositions 返回的有效 position 数量可能小于 positionTokenIds 数量
	validTokenIds := make([]int64, 0)

	for tokenId, p := range positionM {
		// validTokenIds 部分 tokenId 来源是缓存，缓存的 tokenId 可能已经被销毁或者通过 farm 等功能转移
		// 所以需要过滤掉 owner 变更的 tokenId
		if p.Owner.Hex() != common.HexToAddress(address).Hex() {
			continue
		}
		validTokenIds = append(validTokenIds, tokenId)

		token0, err := global.GetToken(chainId, p.Token0.Hex())
		if err != nil {
			logx.Errorf("[SonexPositions]: GetToken error: %s", err.Error())
			return nil, err
		}
		token1, err := global.GetToken(chainId, p.Token1.Hex())
		if err != nil {
			logx.Errorf("[SonexPositions]: GetToken error: %s", err.Error())
			return nil, err
		}

		token0Decimals = append(token0Decimals, int(token0.Decimals))
		token1Decimals = append(token1Decimals, int(token1.Decimals))
		token0Addresses = append(token0Addresses, token0.Address)
		token1Addresses = append(token1Addresses, token1.Address)

		// merge cached position and position from db
		if _, ok := positionTransferItemMap[tokenId]; !ok {
			positionTransferItemMap[tokenId] = &types.Position{
				TokenID: tokenId,
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
				// position id get from cache few seconds after position created,
				// so it's not accurate but approximate.
				CreatedAt: time.Now().UTC().Format(time.RFC3339),
				UpdatedAt: time.Now().UTC().Format(time.RFC3339),
			}
		}
	}

	var wg sync.WaitGroup
	results := make(chan psTaskResult, 2)

	// Task 1: Get collected fees
	wg.Add(1)
	go func() {
		defer wg.Done()
		pairs, err := l.getCollectedPairs(ctx, chainId, validTokenIds, token0Decimals, token1Decimals)
		results <- psTaskResult{collectedPairs: pairs, err: err}
	}()

	// Task 2: Get token prices
	wg.Add(1)
	go func() {
		defer wg.Done()
		pricesRecords, err := l.tokenPriceModel.GetCurrentPrices(ctx, append(token0Addresses, token1Addresses...))
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
			logx.Errorf("[SonexPositions]: waitgroup error: %s", result.err.Error())
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

	positions = make([]*types.Position, 0)

	// Collect results
	for tokenId, position := range positionM {
		if positionTransferItem, exists := positionTransferItemMap[tokenId]; exists {
			convertedPosition, err := l.fillPosition(ctx, chainId, positionTransferItem, position, tokenPricesRecords)
			if err != nil {
				logx.Errorf("[SonexPositions]:refreshCache fillPosition error: %s", err.Error())
				continue
			}
			positions = append(positions, convertedPosition)
		}
	}

	// sort by position in DESC
	sort.SliceStable(positions, func(i, j int) bool {
		return positions[i].Position > positions[j].Position
	})

	return positions, nil
}

// 获取 collectedPairs，过滤掉 deceaseLiquidity events
func (l *SonexPositionsLogic) getCollectedPairs(ctx context.Context, chainId int64, positionTokenIds []int64, token0Decimals, token1Decimals []int) (map[int64]*sonex.TokenPair, error) {
	collectEventsMap, err := l.collectModel.GetCollectEvents(ctx, chainId, positionTokenIds)
	if err != nil {
		logx.Errorf("[SonexPositions]:getCollectedPairs GetCollectEvents error: %s", err.Error())
		return nil, err
	}
	decreaseLiquidityEventsMap, err := l.decreaseLiquidityModel.GetDecreaseLiquidityEvents(ctx, chainId, positionTokenIds)
	if err != nil {
		logx.Errorf("[SonexPositions]:getCollectedPairs GetDecreaseLiquidityEvents error: %s", err.Error())
		return nil, err
	}
	collectedPairs := make(map[int64]*sonex.TokenPair)
	for i, tokenId := range positionTokenIds {
		// 1. initial token pair
		tokenValue0, err := sonex.NewTokenValue("0", token0Decimals[i])
		if err != nil {
			logx.Errorf("[SonexPositions]:getCollectedPairs NewTokenValue error: %s", err.Error())
			return nil, err
		}
		tokenValue1, err := sonex.NewTokenValue("0", token1Decimals[i])
		if err != nil {
			logx.Errorf("[SonexPositions]:getCollectedPairs NewTokenValue error: %s", err.Error())
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
	chainId int64,
	position *types.Position,
	innnerPosition *sonex.Position,
	pricesRecords map[string]*repository.TokenPriceRecord) (*types.Position, error) {
	price0 := new(big.Float).SetFloat64(0)
	price1 := new(big.Float).SetFloat64(0)

	chainConf := l.svcCtx.Config.GetChainConfig(chainId)
	if chainConf == nil {
		return nil, fmt.Errorf("Miss chain configure")
	}
	usdToken, err := global.GetToken(chainId, common.HexToAddress(chainConf.USD).Hex())
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
		logx.Errorf("[SonexPositions]:fillPosition Token0 price not found from db for address: %s", innnerPosition.Token0.Hex())

		bestPrice0, _, _, err := global.GetPrices(ctx, nil, chainId, innnerPosition.Token0.Hex(), true, 2)
		if err == nil && bestPrice0 != nil {
			price0.SetInt(bestPrice0.Value)
			price0String = bestPrice0.Value.String()
		}
		if err != nil {
			logx.Errorf("[SonexPositions]:fillPosition GetPrices error: %s", err.Error())
		}
	}

	if price1Record, exists1 := pricesRecords[innnerPosition.Token1.Hex()]; exists1 {
		price1, _ = price1.SetString(price1Record.Price)
		price1String = price1Record.Price
	} else {
		logx.Errorf("[SonexPositions]:fillPosition Token1 price not found from db for address: %s", innnerPosition.Token1.Hex())

		bestPrice1, _, _, err := global.GetPrices(ctx, nil, chainId, innnerPosition.Token1.Hex(), true, 2)
		if err == nil && bestPrice1 != nil {
			price1.SetInt(bestPrice1.Value)
			price1String = bestPrice1.Value.String()
		}
		if err != nil {
			logx.Errorf("[SonexPositions]:fillPosition GetPrices error: %s", err.Error())
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
	token0Value, err := sonex.NewTokenValue(amount0String, position.Pool.Token0.Decimals, sonex.WithPrice(price0String, usdDecimals))
	if err != nil {
		logx.Errorf("[SonexPositions]:fillPosition NewTokenValue error: %s, amount0: %s", err.Error(), amount0String)
		return nil, err
	}
	token1Value, err := sonex.NewTokenValue(amount1String, position.Pool.Token1.Decimals, sonex.WithPrice(price1String, usdDecimals))
	if err != nil {
		logx.Errorf("[SonexPositions]:fillPosition NewTokenValue error: %s, amount1: %s", err.Error(), amount1String)
		return nil, err
	}
	tokenPair := sonex.NewTokenPair(token0Value, token1Value)
	position.Position, _ = tokenPair.ValueFloat64()

	// 计算未结费用 unclaimedFee0 和 unclaimedFee1
	position.UnclaimedFee0, _ = unclaimedFee0.Float64()
	position.UnclaimedFee1, _ = unclaimedFee1.Float64()

	// 3. claimableFees
	claimableToken0Value, err := sonex.NewTokenValue(unclaimedFee0.String(), position.Pool.Token0.Decimals, sonex.WithPrice(price0String, usdDecimals))
	if err != nil {
		return nil, err
	}
	claimableToken1Value, err := sonex.NewTokenValue(unclaimedFee1.String(), position.Pool.Token1.Decimals, sonex.WithPrice(price1String, usdDecimals))
	if err != nil {
		return nil, err
	}
	claimableTokenPair := sonex.NewTokenPair(claimableToken0Value, claimableToken1Value)
	position.ClaimableFees, _ = claimableTokenPair.ValueFloat64()

	position.UnclaimedFee0, _ = claimableToken0Value.ValueFloat64()
	position.UnclaimedFee1, _ = claimableToken1Value.ValueFloat64()
	// 4. PNL(Fee Revenue) = claimableFees + collectedFees
	pnl := position.ClaimableFees
	collectedPair := position.TotalCollected
	token0Value, err = sonex.NewTokenValue(collectedPair.Token0.Amount, position.Pool.Token0.Decimals, sonex.WithPrice(price0String, usdDecimals))
	if err != nil {
		logx.Errorf("[SonexPositions]:fillPosition NewTokenValue error: %s, collectedPair: %v", err.Error(), collectedPair)
		return nil, err
	}
	token1Value, err = sonex.NewTokenValue(collectedPair.Token1.Amount, position.Pool.Token1.Decimals, sonex.WithPrice(price1String, usdDecimals))
	if err != nil {
		logx.Errorf("[SonexPositions]:fillPosition NewTokenValue error: %s, collectedPair: %v", err.Error(), collectedPair)
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

type LiquidityOperationEventType int64

const (
	LiquidityOperationEventTypeIncrease LiquidityOperationEventType = iota
	LiquidityOperationEventTypeDecrease
)

// EventSonexPositionManagerIncreaseLiquidity 增加流动性事件
type EventSonexPositionManagerIncreaseLiquidity struct {
	repository.EventSonexPositionManagerIncreaseLiquidity
	EventValue *big.Float // token0 + token1 usd value
}

// EventSonexPositionManagerDecreaseLiquidity 减少流动性事件
type EventSonexPositionManagerDecreaseLiquidity struct {
	repository.EventSonexPositionManagerDecreaseLiquidity
	EventValue *big.Float // token0 + token1 usd value
}

// LiquidityOperationEvent 流动性操作事件
type LiquidityOperationEvent struct {
	TokenId                     int64
	LiquidityOperationEventType LiquidityOperationEventType
	IncreaseEvent               *EventSonexPositionManagerIncreaseLiquidity
	DecreaseEvent               *EventSonexPositionManagerDecreaseLiquidity
}

// GetLiquidityOperationEvents 获取 liquidity operation events, by user address, 没有计算出 token0 + token1 usd value
func (l *SonexPositionsLogic) GetLiquidityOperationEvents(ctx context.Context, chainId int64, address string) (m map[int64][]*LiquidityOperationEvent, err error) {
	type increaseLiquidityEventsResult struct {
		events map[int64][]*repository.EventSonexPositionManagerIncreaseLiquidity
		err    error
	}
	type decreaseLiquidityEventsResult struct {
		events map[int64][]*repository.EventSonexPositionManagerDecreaseLiquidity
		err    error
	}

	increaseLiquidityEventsChan := make(chan increaseLiquidityEventsResult, 1)
	decreaseLiquidityEventsChan := make(chan decreaseLiquidityEventsResult, 1)

	go func() {
		increaseLiquidityEvents, err := l.increaseLiquidityModel.GetIncreaseLiquidityEventsByAddress(ctx, chainId, address)
		increaseLiquidityEventsChan <- increaseLiquidityEventsResult{events: increaseLiquidityEvents, err: err}
	}()

	go func() {
		decreaseLiquidityEvents, err := l.decreaseLiquidityModel.GetDecreaseLiquidityEventsByAddress(ctx, chainId, address)
		decreaseLiquidityEventsChan <- decreaseLiquidityEventsResult{events: decreaseLiquidityEvents, err: err}
	}()

	increaseLiquidityEvents := <-increaseLiquidityEventsChan
	decreaseLiquidityEvents := <-decreaseLiquidityEventsChan

	if increaseLiquidityEvents.err != nil {
		logx.Errorf("[SonexPositions]:GetLiquidityOperationEvents error: %s", increaseLiquidityEvents.err.Error())
		return nil, increaseLiquidityEvents.err
	}
	if decreaseLiquidityEvents.err != nil {
		logx.Errorf("[SonexPositions]:GetLiquidityOperationEvents error: %s", decreaseLiquidityEvents.err.Error())
		return nil, decreaseLiquidityEvents.err
	}

	m = make(map[int64][]*LiquidityOperationEvent)

	for tokenId, events := range increaseLiquidityEvents.events {
		for _, event := range events {
			m[tokenId] = append(m[tokenId], &LiquidityOperationEvent{
				TokenId:                     tokenId,
				LiquidityOperationEventType: LiquidityOperationEventTypeIncrease,
				IncreaseEvent: &EventSonexPositionManagerIncreaseLiquidity{
					EventSonexPositionManagerIncreaseLiquidity: *event,
				},
			})
		}
	}

	for tokenId, events := range decreaseLiquidityEvents.events {
		for _, event := range events {
			m[tokenId] = append(m[tokenId], &LiquidityOperationEvent{
				TokenId:                     tokenId,
				LiquidityOperationEventType: LiquidityOperationEventTypeDecrease,
				DecreaseEvent: &EventSonexPositionManagerDecreaseLiquidity{
					EventSonexPositionManagerDecreaseLiquidity: *event,
				},
			})
		}
	}

	return m, nil
}

// GetLiquidityOperationEventsWithTokenValue 获取 liquidity operation events, by user address, 计算出 token0 + token1 usd value
func (l *SonexPositionsLogic) GetLiquidityOperationEventsWithTokenValue(ctx context.Context, chainId int64, address string) (m map[int64][]*LiquidityOperationEvent, err error) {
	sonexConf := l.svcCtx.Config.GetSonexConfig(chainId)
	if sonexConf == nil {
		logx.Errorf("[SonexPositions]:GetLiquidityOperationEventsWithTokenValue Miss sonex configure")
		return nil, fmt.Errorf("Miss sonex configure")
	}

	events, err := l.GetLiquidityOperationEvents(ctx, chainId, address)
	if err != nil {
		return nil, err
	}

	tokenIds := make([]*big.Int, 0)
	for tokenId := range events {
		tokenIds = append(tokenIds, big.NewInt(tokenId))
	}

	staticCaller, err := sonex.NewStaticCaller(l.svcCtx.Config)
	if err != nil {
		logx.Errorf("[SonexPositions]:GetLiquidityOperationEventsWithTokenValue NewStaticCaller error: %s", err.Error())
		return nil, err
	}
	pm, err := staticCaller.BatchGetPositionMap(l.ctx, chainId, common.HexToAddress(sonexConf.PositionManager), tokenIds)
	if err != nil {
		logx.Errorf("[SonexPositions]:GetLiquidityOperationEventsWithTokenValue BatchGetPositionMap error: %s", err.Error())
		return nil, err
	}

	tokenAddresses := make([]string, 0)
	for _, position := range pm {
		tokenAddresses = append(tokenAddresses, position.Token0.Hex())
		tokenAddresses = append(tokenAddresses, position.Token1.Hex())
	}

	tokenPricesRecords, err := l.getCurrentPrices(ctx, tokenAddresses)
	if err != nil {
		return nil, err
	}

	for tokenId, list := range events {
		for _, event := range list {
			// increase liquidity event
			if event.LiquidityOperationEventType == LiquidityOperationEventTypeIncrease {
				event.IncreaseEvent.EventValue = new(big.Float)
			}
			// decrease liquidity event
			if event.LiquidityOperationEventType == LiquidityOperationEventTypeDecrease {
				event.DecreaseEvent.EventValue = new(big.Float)
			}

			// todo: 被 burn 掉的 token 无法参与计算，后期需要优化，从 db 中直接 join 出 token 信息，不直接从链上取
			position, ok := pm[tokenId]
			if !ok {
				continue
			}

			token0Price, ok0 := tokenPricesRecords[position.Token0.Hex()]
			if !ok0 {
				continue
			}
			token1Price, ok1 := tokenPricesRecords[position.Token1.Hex()]
			if !ok1 {
				continue
			}
			token0, err := global.GetToken(chainId, position.Token0.Hex())
			if err != nil {
				continue
			}
			token1, err := global.GetToken(chainId, position.Token1.Hex())
			if err != nil {
				continue
			}

			token0Value, err := sonex.NewTokenValue(event.IncreaseEvent.EventAmount0, int(token0.Decimals), sonex.WithTokenPrice(token0Price))
			if err != nil {
				continue
			}
			token1Value, err := sonex.NewTokenValue(event.IncreaseEvent.EventAmount1, int(token1.Decimals), sonex.WithTokenPrice(token1Price))
			if err != nil {
				continue
			}
			tokenPair := sonex.NewTokenPair(token0Value, token1Value)

			// increase liquidity event
			if event.LiquidityOperationEventType == LiquidityOperationEventTypeIncrease {
				event.IncreaseEvent.EventValue = tokenPair.Value()
			}
			// decrease liquidity event
			if event.LiquidityOperationEventType == LiquidityOperationEventTypeDecrease {
				event.DecreaseEvent.EventValue = tokenPair.Value()
			}
		}
	}

	return events, nil
}

// getCurrentPrices 获取当前价格
func (l *SonexPositionsLogic) getCurrentPrices(ctx context.Context, addressList []string) (prices map[string]*sonex.TokenPrice, err error) {
	records, err := l.tokenPriceModel.GetCurrentPrices(ctx, addressList)
	if err != nil {
		logx.Errorf("[SonexPositions]:getCurrentPrices Get current prices error: %s", err)
		return nil, err
	}
	prices = make(map[string]*sonex.TokenPrice, 0)
	for _, record := range records {
		tp, err := sonex.NewTokenPrice(record.Price, record.UsdDecimals)
		if err != nil {
			logx.Errorf("[SonexPositions]:getCurrentPrices Invalid token price: %s", err)
			return nil, err
		}
		prices[record.TokenAddress] = tp
	}

	// todo db 获取不到的话，从链上获取价格
	return prices, nil
}
