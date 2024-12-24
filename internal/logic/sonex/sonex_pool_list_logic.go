package sonex

import (
	"context"
	"encoding/json"
	"math"
	"strconv"
	"sync"
	"time"

	"sonex_api_backend_go/internal/lib/global"
	lp "sonex_api_backend_go/internal/lib/pool"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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

func (l *SonexPoolListLogic) SonexPoolList(req *types.PoolsRequest) (resp *types.PoolsResponse, err error) {
	const cacheKey = "sonex_dex:uniswap:pool_list"

	cachedResponse, err := l.svcCtx.Cache.Get(cacheKey)
	if err == nil && len(cachedResponse) > 0 {
		var cachedResp types.PoolsResponse
		if err := json.Unmarshal([]byte(cachedResponse), &cachedResp); err == nil {
			return &cachedResp, nil
		}
	}

	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	model := repository.NewEventSonexFactoryPoolCreatedModel(conn)
	data, meta, err := model.GetPoolList(l.ctx, req, l.svcCtx.Config.Sonexes[0].PositionManager)
	if err != nil {
		return &types.PoolsResponse{
			Code:    400,
			Data:    []*types.PoolsResponse_Pool{},
			Meta:    &types.PaginationMeta{},
			Message: err.Error(),
		}, nil
	}

	balanceFetcher := lp.NewTokenBalanceFetcher(l.svcCtx.Config.Chains[0])
	var wg sync.WaitGroup
	mu := &sync.Mutex{}

	for i := range data {
		wg.Add(1)
		go func(pool *types.PoolsResponse_Pool) {
			defer wg.Done()

			balance0, _ := balanceFetcher.GetTokenBalance(pool.Address, pool.Token0.Address, int64(pool.Token0.Decimals))
			balance1, _ := balanceFetcher.GetTokenBalance(pool.Address, pool.Token1.Address, int64(pool.Token1.Decimals))
			price0 := float64(0)
			price1 := float64(0)
			bestPrice0, _, err := global.GetPrices(l.ctx, nil, req.ChainId, pool.Token0.Address, true, 2, nil)
			if err == nil && bestPrice0 != nil {
				price0, _ = strconv.ParseFloat(bestPrice0.Value.String(), 64)
			}

			pool.Price0 = price0

			bestPrice1, _, err := global.GetPrices(l.ctx, nil, req.ChainId, pool.Token0.Address, true, 2, nil)
			if err == nil && bestPrice1 != nil {
				price1, _ = strconv.ParseFloat(bestPrice1.Value.String(), 64)
			}

			pool.Price1 = price1

			mu.Lock()
			pool.Token0Balance = balance0 * price0
			pool.Token1Balance = balance1 * price1
			pool.TvlUSD = pool.Token0Balance + pool.Token1Balance
			mu.Unlock()
		}(data[i])
	}

	wg.Wait()

	poolVolumes, err := model.GetPoolVolumes(l.ctx, l.svcCtx.Config.Sonexes[0].ChainId, true)
	if err == nil {
		for _, pool := range data {
			for _, volume := range poolVolumes {
				if pool.Address == volume.Address {
					factor0 := math.Pow(10, float64(pool.Token0.Decimals))
					factor1 := math.Pow(10, float64(pool.Token1.Decimals))
					pool.Volume = volume.Amounts0/factor0 + volume.Amounts1/factor1
					break
				}
			}
		}
	}

	response := &types.PoolsResponse{
		Code:    200,
		Data:    data,
		Meta:    meta,
		Message: "success",
	}

	responseBytes, err := json.Marshal(response)
	if err == nil {
		l.svcCtx.Cache.Set(cacheKey, responseBytes, 20*time.Second)
	}

	return response, nil
}
