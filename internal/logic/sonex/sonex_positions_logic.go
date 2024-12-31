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

type SonexPositionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexPositionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexPositionsLogic {
	return &SonexPositionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type PositionsResponseWrapper struct {
	*types.PositionsResponse
}

func (p *PositionsResponseWrapper) GetCode() interface{} {
	return p.Code
}

func (l *SonexPositionsLogic) SonexPositions(req *types.PositionsRequest) (resp *PositionsResponseWrapper, err error) {
	cacheKey := fmt.Sprintf("sonex_dex:uniswap:positions:%s", req.Address)

	// Check cache
	redisCache := cache.NewRedisCache(fmt.Sprintf("%s:%d", l.svcCtx.Config.Redis.Host, l.svcCtx.Config.Redis.Port), l.svcCtx.Config.Redis.Password, l.svcCtx.Config.Redis.DB)
	cachedResponse, err := redisCache.Get(cacheKey)
	if err == nil && len(cachedResponse) > 0 {
		var cachedResp PositionsResponseWrapper
		if err := json.Unmarshal([]byte(cachedResponse), &cachedResp); err == nil {
			return &cachedResp, nil
		}
	}

	// Fetch from database if not cached
	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	model := repository.NewEventSonexPositionManagerTransferModel(conn)
	items, err := model.GetTransferEvents(l.ctx, req.Address)
	if err != nil {
		return &PositionsResponseWrapper{
			&types.PositionsResponse{
				Code:      400,
				Positions: []*types.Position{},
				Message:   err.Error(),
			},
		}, nil
	}

	positionTokenIds := make([]*big.Int, 0)
	positionTransferItemMap := make(map[int64]*types.Position)

	for _, item := range items {
		if !item.IsClosed {
			tokenId := item.TokenID
			positionTokenIds = append(positionTokenIds, big.NewInt(tokenId))
			positionTransferItemMap[tokenId] = item
		}
	}

	caller, _ := sonex.NewStaticCaller(l.svcCtx.Config)
	positions, err := caller.BatchGetPositions(
		context.Background(),
		req.ChainId,
		common.HexToAddress(l.svcCtx.Config.GetSonexConfig(req.ChainId).PositionManager),
		positionTokenIds,
	)

	resultPositions := make([]*types.Position, 0)
	var wg sync.WaitGroup
	positionChan := make(chan *types.Position, len(positions))

	for _, position := range positions {
		if positionTransferItem, exists := positionTransferItemMap[position.TokenId.Int64()]; exists {
			wg.Add(1)
			go func(pos *sonex.Position, transferItem *types.Position) {
				defer wg.Done()
				convertedPosition := l.convertPosition(l.ctx, req, transferItem, pos)
				positionChan <- convertedPosition
			}(position, positionTransferItem)
		}
	}

	// Close channel after all goroutines complete
	go func() {
		wg.Wait()
		close(positionChan)
	}()

	totalBalance := new(big.Float).SetInt(big.NewInt(0))
	// Collect results
	for position := range positionChan {
		resultPositions = append(resultPositions, position)
		totalBalance.Add(totalBalance, big.NewFloat(position.Balance))
	}

	positionResponse := &types.PositionsResponse{
		Positions: resultPositions,
	}

	positionResponse.TotalBalances, _ = totalBalance.Float64()

	response := &PositionsResponseWrapper{
		PositionsResponse: positionResponse,
	}

	responseBytes, err := json.Marshal(response)
	if err == nil {
		redisCache.Set(cacheKey, responseBytes, 10*time.Second)
	}

	return response, nil
}

func (l *SonexPositionsLogic) convertPosition(ctx context.Context, req *types.PositionsRequest, position *types.Position, innnerPosition *sonex.Position) *types.Position {
	price0 := new(big.Float).SetFloat64(0)
	price1 := new(big.Float).SetFloat64(0)

	bestPrice0, _, err := global.GetPrices(ctx, nil, req.ChainId, innnerPosition.Token0.Hex(), true, 2, nil)
	if err == nil && bestPrice0 != nil {
		price0.SetInt(bestPrice0.Value)
	}

	factor0 := new(big.Float).SetFloat64(math.Pow10(int(position.Pool.Token0.Decimals)))
	factor1 := new(big.Float).SetFloat64(math.Pow10(int(position.Pool.Token1.Decimals)))

	poolPrice0 := new(big.Float).Quo(price0, factor0)
	position.Price0, _ = poolPrice0.Float64()
	position.Pool.Token0Price = poolPrice0.String()

	bestPrice1, _, err := global.GetPrices(ctx, nil, req.ChainId, innnerPosition.Token1.Hex(), true, 2, nil)
	if err == nil && bestPrice1 != nil {
		price1.SetInt(bestPrice1.Value)
	}

	// 转换 price 为 *big.Float 类型，并根据 factor0 和 factor1 做调整
	poolPrice1 := new(big.Float).Quo(price1, factor1)
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
	amount0Balance := new(big.Float).Mul(amount0Float, poolPrice0)
	amount1Balance := new(big.Float).Mul(amount1Float, poolPrice1)

	position.Amount0, _ = amount0.Float64()
	position.Amount1, _ = amount1.Float64()

	// 计算总余额
	totalBalance := new(big.Float).Add(amount0Balance, amount1Balance)
	position.Balance, _ = totalBalance.Float64()

	// 计算未结费用 unclaimedFee0 和 unclaimedFee1
	position.UnclaimedFee0, _ = new(big.Float).SetInt(unclaimedFee0).Float64()
	position.UnclaimedFee1, _ = new(big.Float).SetInt(unclaimedFee1).Float64()

	// 计算 claimable fees
	unclaimedFee0Float := new(big.Float).SetInt(unclaimedFee0)
	unclaimedFee1Float := new(big.Float).SetInt(unclaimedFee1)
	unclaimedFee0Float.Quo(unclaimedFee0Float, factor0)
	unclaimedFee1Float.Quo(unclaimedFee1Float, factor1)

	claimableFee0 := new(big.Float).Mul(unclaimedFee0Float, poolPrice0)
	claimableFee1 := new(big.Float).Mul(unclaimedFee1Float, poolPrice1)

	// 计算总的 claimableFees
	claimableFees := new(big.Float).Add(claimableFee0, claimableFee1)
	position.ClaimableFees, _ = claimableFees.Float64()

	return position
}
