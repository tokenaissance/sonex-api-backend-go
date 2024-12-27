package sonex

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"sync"
	"time"

	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SonexSwapLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexSwapLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexSwapLogic {
	return &SonexSwapLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type SwapEventResponseWrapper struct {
	*types.SwapEventResponse
}

func (p *SwapEventResponseWrapper) GetData() interface{} {
	return p.Data
}

func (p *SwapEventResponseWrapper) GetMeta() interface{} {
	return p.Meta
}

func (p *SwapEventResponseWrapper) GetCode() interface{} {
	return p.Code
}

func (l *SonexSwapLogic) SonexSwap(req *types.SwapEventRequest) (resp *SwapEventResponseWrapper, err error) {
	cacheKey := "sonex_dex:uniswap:swap"

	cachedResponse, err := l.svcCtx.Cache.Get(cacheKey)
	if err == nil && len(cachedResponse) > 0 {
		var cachedResp SwapEventResponseWrapper
		if err := json.Unmarshal([]byte(cachedResponse), &cachedResp); err == nil {
			go l.refreshCache(req, cacheKey, &cachedResp)
			return &cachedResp, nil
		}
	}

	return l.refreshCache(req, cacheKey, nil)
}

func (l *SonexSwapLogic) refreshCache(req *types.SwapEventRequest, cacheKey string, cachedResp *SwapEventResponseWrapper) (resp *SwapEventResponseWrapper, err error) {
	lockKey := fmt.Sprintf("lock:%s", cacheKey)
	ctx := context.Background()

	var lockAcquired bool

	if cachedResp != nil {
		success, err := l.svcCtx.Cache.SetNX(lockKey, "locked", 30*time.Second)
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
			l.svcCtx.Cache.Del(lockKey)
		}
	}()

	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	model := repository.NewEventSonexPoolSwapModel(conn)
	data, meta, err := model.GetSwapEvents(ctx, req)
	if err != nil {
		return &SwapEventResponseWrapper{
			&types.SwapEventResponse{
				Code:    400,
				Data:    []*types.SwapEvent{},
				Meta:    &types.PaginationMeta{},
				Message: err.Error(),
			},
		}, nil
	}

	var wg sync.WaitGroup
	mu := &sync.Mutex{}

	for _, item := range data {
		wg.Add(1)
		go func(item *types.SwapEvent) {
			defer wg.Done()

			bestPrice, _, err := global.GetPrices(ctx, nil, req.ChainId, item.Token1.Address, true, 2, nil)
			if err == nil && bestPrice != nil {
				priceFloat := new(big.Float).SetInt(bestPrice.Value)
				amount1Positive := new(big.Float).SetFloat64(math.Abs(item.Amount1))
				decimalFactor := new(big.Float).SetFloat64(math.Pow10(int(item.Token1.Decimals)))
				usdValue := new(big.Float).Quo(new(big.Float).Mul(priceFloat, amount1Positive), decimalFactor)

				mu.Lock()
				item.Usd, _ = usdValue.Float64()
				item.Amount1, _ = amount1Positive.Float64()
				mu.Unlock()
			}
		}(item)
	}

	wg.Wait()

	response := &SwapEventResponseWrapper{
		&types.SwapEventResponse{
			Code:    200,
			Data:    data,
			Meta:    meta,
			Message: "success",
		},
	}

	responseBytes, err := json.Marshal(response)
	if err == nil {
		err = l.svcCtx.Cache.Set(cacheKey, responseBytes, 20*time.Second)
		if err != nil {
			l.Error("Failed to cache response", err)
		}
	}

	return response, nil
}
