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
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type SonexPositionValuesLogic struct {
	logx.Logger
	ctx                    context.Context
	svcCtx                 *svc.ServiceContext
	transferModel          repository.EventSonexPositionManagerTransferModel
	collectModel           repository.EventSonexPositionManagerCollectModel
	decreaseLiquidityModel repository.EventSonexPositionManagerDecreaseLiquidityModel
	tokenPriceModel        repository.SonexTokenPriceRecordModel
}

func NewSonexPositionValuesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexPositionValuesLogic {
	conn := repository.NewConn(svcCtx.Config.PostgreSQL.DateSource())
	transferModel := repository.NewEventSonexPositionManagerTransferModel(conn)
	collectModel := repository.NewEventSonexPositionManagerCollectModel(conn)
	decreaseLiquidityModel := repository.NewEventSonexPositionManagerDecreaseLiquidityModel(conn)
	tokenPriceModel := repository.NewSonexTokenPriceRecordModel(conn)
	return &SonexPositionValuesLogic{
		Logger:                 logx.WithContext(ctx),
		ctx:                    ctx,
		svcCtx:                 svcCtx,
		transferModel:          transferModel,
		collectModel:           collectModel,
		decreaseLiquidityModel: decreaseLiquidityModel,
		tokenPriceModel:        tokenPriceModel,
	}
}

type PositionValuesResponseWrapper struct {
	*types.PositionValuesResponse
}

type pvTaskResult struct {
	collectedPairs     map[int64]*sonex.TokenPair
	positions          []*sonex.Position
	tokenPricesRecords map[string]*repository.TokenPriceRecord
	err                error
}

func (l *SonexPositionValuesLogic) SonexPositionValuesMethod(chainId int64, address string) (resp *sonex.PositionValues, err error) {
	cacheKey := fmt.Sprintf("sonex_dex:uniswap:position_values:%s", address)

	// Check cache
	redisCache := cache.NewRedisCache(fmt.Sprintf("%s:%d", l.svcCtx.Config.Redis.Host, l.svcCtx.Config.Redis.Port), l.svcCtx.Config.Redis.Password, l.svcCtx.Config.Redis.DB)
	cachedResponse, err := redisCache.Get(cacheKey)
	if err == nil && len(cachedResponse) > 0 {
		var cachedResp sonex.PositionValues
		if err := json.Unmarshal([]byte(cachedResponse), &cachedResp); err == nil {
			return &cachedResp, nil
		}
	}

	// Fetch from database if not cached
	items, err := l.transferModel.GetTransferEvents(l.ctx, chainId, address, true)
	if err != nil {
		logx.Errorf("[SonexPositionValuesMethod] GetTransferEvents error: %s", err.Error())
		return nil, err
	}

	positionTokenIds := make([]int64, 0)
	positionTransferItemMap := make(map[int64]*types.Position)
	token0Decimals := make([]int, 0)
	token0Addresses := make([]string, 0)
	token1Decimals := make([]int, 0)
	token1Addresses := make([]string, 0)
	for _, item := range items {
		if !item.IsBurnt {
			tokenId := item.TokenID
			positionTokenIds = append(positionTokenIds, tokenId)
			positionTransferItemMap[tokenId] = item
			token0Decimals = append(token0Decimals, item.Pool.Token0.Decimals)
			token1Decimals = append(token1Decimals, item.Pool.Token1.Decimals)
			token0Addresses = append(token0Addresses, item.Pool.Token0.Address)
			token1Addresses = append(token1Addresses, item.Pool.Token1.Address)
		}
	}

	var wg sync.WaitGroup
	results := make(chan pvTaskResult, 3)

	// Task 1: Get collected fees
	wg.Add(1)
	go func() {
		defer wg.Done()
		pairs, err := l.getCollectedPairs(l.ctx, chainId, positionTokenIds, token0Decimals, token1Decimals)
		results <- pvTaskResult{collectedPairs: pairs, err: err}
	}()

	// Task 2: Get positions
	wg.Add(1)
	go func() {
		defer wg.Done()
		caller, _ := sonex.NewStaticCaller(l.svcCtx.Config)
		pos, err := caller.BatchGetPositions(
			context.Background(),
			chainId,
			common.HexToAddress(l.svcCtx.Config.GetSonexConfig(chainId).PositionManager),
			func(tokenIds []int64) []*big.Int {
				converted := make([]*big.Int, len(tokenIds))
				for i, tokenId := range tokenIds {
					converted[i] = big.NewInt(tokenId)
				}
				return converted
			}(positionTokenIds),
		)
		results <- pvTaskResult{positions: pos, err: err}
	}()

	// Task 3: Get token prices
	wg.Add(1)
	go func() {
		defer wg.Done()
		pricesRecords, err := l.tokenPriceModel.GetCurrentPrices(append(token0Addresses, token1Addresses...))
		results <- pvTaskResult{tokenPricesRecords: pricesRecords, err: err}
	}()

	// Close results channel after all goroutines complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var positions = make([]*sonex.Position, 0)
	var tokenPricesRecords = make(map[string]*repository.TokenPriceRecord)

	for result := range results {
		if result.err != nil {
			logx.Errorf("[SonexPositionValuesMethod] waitgroup error: %s", result.err.Error())
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
		if result.positions != nil {
			positions = result.positions
		}
		if result.tokenPricesRecords != nil {
			tokenPricesRecords = result.tokenPricesRecords
		}
	}

	positionValues := &sonex.PositionValues{
		PositionValues: make([]*types.PositionValueRecord, 0),
	}

	// Collect results
	for _, position := range positions {
		if positionTransferItem, exists := positionTransferItemMap[position.TokenId.Int64()]; exists {
			convertedPosition, err := l.fillPosition(l.ctx, chainId, positionTransferItem, position, tokenPricesRecords)
			if err != nil {
				logx.Errorf("[SonexPositionValuesMethod] fillPosition error: %s", err.Error())
				return nil, err
			}
			positionValues.TotalPosition += convertedPosition.Position
			positionValues.Pnl += convertedPosition.PNL
			positionValues.TotalClaimable += convertedPosition.ClaimableFees
		}
	}

	positionValuesBytes, err := json.Marshal(positionValues)
	if err == nil {
		redisCache.Set(cacheKey, positionValuesBytes, 10*time.Second)
	}

	return positionValues, nil
}

func (l *SonexPositionValuesLogic) SonexPositionValues(req *types.PositionValuesRequest) (resp *PositionValuesResponseWrapper, err error) {
	positionValues, err := l.SonexPositionValuesMethod(req.ChainId, req.Address)
	if err != nil {
		logx.Errorf("[SonexPositionValues] SonexPositionValuesMethod error: %s", err.Error())
		return nil, err
	}

	response := &PositionValuesResponseWrapper{
		PositionValuesResponse: &types.PositionValuesResponse{
			TotalPosition:  positionValues.TotalPosition,
			Pnl:            positionValues.Pnl,
			PositionValues: positionValues.PositionValues,
		},
	}

	return response, nil
}

// 获取 collectedPairs，过滤掉 deceaseLiquidity events
func (l *SonexPositionValuesLogic) getCollectedPairs(ctx context.Context, chainId int64, positionTokenIds []int64, token0Decimals, token1Decimals []int) (map[int64]*sonex.TokenPair, error) {
	collectEventsMap, err := l.collectModel.GetCollectEvents(ctx, chainId, positionTokenIds)
	if err != nil {
		logx.Errorf("[SonexPositionValuesMethod]:getCollectedPairs GetCollectEvents error: %s", err.Error())
		return nil, err
	}
	decreaseLiquidityEventsMap, err := l.decreaseLiquidityModel.GetDecreaseLiquidityEvents(ctx, chainId, positionTokenIds)
	if err != nil {
		logx.Errorf("[SonexPositionValuesMethod]:getCollectedPairs GetDecreaseLiquidityEvents error: %s", err.Error())
		return nil, err
	}
	collectedPairs := make(map[int64]*sonex.TokenPair)
	for i, tokenId := range positionTokenIds {
		// 1. initial token pair
		tokenValue0, err := sonex.NewTokenValue("0", token0Decimals[i])
		if err != nil {
			logx.Errorf("[SonexPositionValuesMethod]:getCollectedPairs NewTokenValue error: %s", err.Error())
			return nil, err
		}
		tokenValue1, err := sonex.NewTokenValue("0", token1Decimals[i])
		if err != nil {
			logx.Errorf("[SonexPositionValuesMethod]:getCollectedPairs NewTokenValue error: %s", err.Error())
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

func (l *SonexPositionValuesLogic) fillPosition(ctx context.Context,
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
		logx.Errorf("[SonexPositionValuesMethod]:fillPosition GetToken error: %s", err.Error())
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
		logx.Errorf("[SonexPositionValuesMethod]:fillPosition Token0 price not found from db for address: %s", innnerPosition.Token0.Hex())

		bestPrice0, _, _, err := global.GetPrices(ctx, nil, chainId, innnerPosition.Token0.Hex(), true, 2)
		if err == nil && bestPrice0 != nil {
			price0.SetInt(bestPrice0.Value)
			price0String = bestPrice0.Value.String()
		}
		if err != nil {
			logx.Errorf("[SonexPositionValuesMethod]:fillPosition GetPrices error: %s", err.Error())
		}
	}

	if price1Record, exists1 := pricesRecords[innnerPosition.Token1.Hex()]; exists1 {
		price1, _ = price1.SetString(price1Record.Price)
		price1String = price1Record.Price
	} else {
		logx.Errorf("[SonexPositionValuesMethod]:fillPosition Token1 price not found from db for address: %s", innnerPosition.Token1.Hex())

		bestPrice1, _, _, err := global.GetPrices(ctx, nil, chainId, innnerPosition.Token1.Hex(), true, 2)
		if err == nil && bestPrice1 != nil {
			price1.SetInt(bestPrice1.Value)
			price1String = bestPrice1.Value.String()
		}
		if err != nil {
			logx.Errorf("[SonexPositionValuesMethod]:fillPosition GetPrices error: %s", err.Error())
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
		logx.Errorf("[SonexPositionValuesMethod]:fillPosition NewTokenValue error: %s", err.Error())
		return nil, err
	}
	token1Value, err := sonex.NewTokenValue(amount1String, position.Pool.Token1.Decimals, sonex.WithPrice(price1String, usdDecimals))
	if err != nil {
		logx.Errorf("[SonexPositionValuesMethod]:fillPosition NewTokenValue error: %s", err.Error())
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
		logx.Errorf("[SonexPositionValuesMethod]:fillPosition NewTokenValue error: %s", err.Error())
		return nil, err
	}
	claimableToken1Value, err := sonex.NewTokenValue(unclaimedFee1.String(), position.Pool.Token1.Decimals, sonex.WithPrice(price1String, usdDecimals))
	if err != nil {
		logx.Errorf("[SonexPositionValuesMethod]:fillPosition NewTokenValue error: %s", err.Error())
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
		logx.Errorf("[SonexPositionValuesMethod]:fillPosition NewTokenValue error: %s", err.Error())
		return nil, err
	}
	token1Value, err = sonex.NewTokenValue(collectedPair.Token1.Amount, position.Pool.Token1.Decimals, sonex.WithPrice(price1String, usdDecimals))
	if err != nil {
		logx.Errorf("[SonexPositionValuesMethod]:fillPosition NewTokenValue error: %s", err.Error())
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
