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
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexPositionValuesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexPositionValuesLogic {
	return &SonexPositionValuesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type PositionValuesResponseWrapper struct {
	*types.PositionValuesResponse
}

type taskResult struct {
	collectedAmounts   map[int64]*types.AggregatedAmounts
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
	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	transferModel := repository.NewEventSonexPositionManagerTransferModel(conn)
	collectModel := repository.NewEventSonexPositionManagerCollectModel(conn)
	tokenPriceModel := repository.NewSonexTokenPriceRecordModel(conn)
	items, err := transferModel.GetTransferEvents(l.ctx, address, true)
	if err != nil {
		logx.Errorf("[SonexPositionsLogic] GetTransferEvents error: %s", err.Error())
		return nil, err
	}

	positionTokenIds := make([]int64, 0)
	positionTransferItemMap := make(map[int64]*types.Position)
	token0Decimals := make([]int, 0)
	token0Addresses := make([]string, 0)
	token1Decimals := make([]int, 0)
	token1Addresses := make([]string, 0)
	for _, item := range items {
		if !item.IsClosed {
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
	results := make(chan taskResult, 3)

	// Task 1: Get collected fees
	wg.Add(1)
	go func() {
		defer wg.Done()
		amounts, err := collectModel.BatchAggregateAmounts(l.ctx, chainId, positionTokenIds, token0Decimals, token1Decimals)
		results <- taskResult{collectedAmounts: amounts, err: err}
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
		results <- taskResult{positions: pos, err: err}
	}()

	// Task 3: Get token prices
	wg.Add(1)
	go func() {
		defer wg.Done()
		pricesRecords, err := tokenPriceModel.GetCurrentPrices(append(token0Addresses, token1Addresses...))
		results <- taskResult{tokenPricesRecords: pricesRecords, err: err}
	}()

	// Close results channel after all goroutines complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var collectedAmounts = make(map[int64]*types.AggregatedAmounts)
	var positions = make([]*sonex.Position, 0)
	var tokenPricesRecords = make(map[string]*repository.TokenPriceRecord)

	for result := range results {
		if result.err != nil {
			logx.Errorf("[SonexPositions] waitgroup error: %s", result.err.Error())
			return nil, result.err
		}
		if result.collectedAmounts != nil {
			collectedAmounts = result.collectedAmounts
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
			convertedPosition := l.fillPosition(l.ctx, chainId, positionTransferItem, position, tokenPricesRecords, collectedAmounts)
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

func (l *SonexPositionValuesLogic) fillPosition(ctx context.Context,
	chainId int64,
	position *types.Position,
	innnerPosition *sonex.Position,
	pricesRecords map[string]*repository.TokenPriceRecord,
	collectedAmounts map[int64]*types.AggregatedAmounts) *types.Position {
	price0 := new(big.Float).SetFloat64(0)
	price1 := new(big.Float).SetFloat64(0)
	factorDai := new(big.Float).SetFloat64(math.Pow10(18))
	// 优先获取从 db 中的价格，如果获取失败，则从链上获取
	if price0Record, exists0 := pricesRecords[innnerPosition.Token0.Hex()]; exists0 {
		price0, _ = price0.SetString(price0Record.Price)
	} else {
		logx.Errorf("[convertPosition] Token0 price not found from db for address: %s", innnerPosition.Token0.Hex())

		bestPrice0, _, err := global.GetPrices(ctx, nil, chainId, innnerPosition.Token0.Hex(), true, 2, &l.svcCtx.Config.GetChainConfig(chainId).USD)
		if err == nil && bestPrice0 != nil {
			price0.SetInt(bestPrice0.Value)
		}
	}

	if price1Record, exists1 := pricesRecords[innnerPosition.Token1.Hex()]; exists1 {
		price1, _ = price1.SetString(price1Record.Price)
	} else {
		logx.Errorf("[convertPosition] Token1 price not found from db for address: %s", innnerPosition.Token1.Hex())

		bestPrice1, _, err := global.GetPrices(ctx, nil, chainId, innnerPosition.Token1.Hex(), true, 2, &l.svcCtx.Config.GetChainConfig(chainId).USD)
		if err == nil && bestPrice1 != nil {
			price1.SetInt(bestPrice1.Value)
		}
	}

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

	amount0Float = amount0Float.Quo(amount0Float, factor0)
	amount1Float = amount1Float.Quo(amount1Float, factor1)

	// 计算 amount0 和 amount1 的余额
	positionValue0 := new(big.Float).Mul(amount0Float, poolPrice0)
	positionValue1 := new(big.Float).Mul(amount1Float, poolPrice1)

	// 1. 计算 amount0 和 amount1
	position.Amount0, _ = amount0Float.Float64()
	position.Amount1, _ = amount1Float.Float64()

	// 2. 计算总 position
	totalPosition := new(big.Float).Add(positionValue0, positionValue1)
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
	collectedFees, exists := collectedAmounts[position.TokenID]
	if exists {
		position.PNL = position.ClaimableFees + collectedFees.TotalAmount0 + collectedFees.TotalAmount1
	} else {
		position.PNL = position.ClaimableFees
	}

	return position
}
