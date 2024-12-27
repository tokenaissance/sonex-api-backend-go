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
	lp "sonex_api_backend_go/internal/lib/pool"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type SonexPoolListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

type PoolsResponseWrapper struct {
	*types.PoolsResponse
}

func (p *PoolsResponseWrapper) GetData() interface{} {
	return p.Data
}

func (p *PoolsResponseWrapper) GetMeta() interface{} {
	return p.Meta
}

func (p *PoolsResponseWrapper) GetCode() interface{} {
	return p.Code
}

func NewSonexPoolListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexPoolListLogic {
	return &SonexPoolListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SonexPoolListLogic) SonexPoolList(req *types.PoolsRequest) (resp *PoolsResponseWrapper, err error) {
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
	cachedResponse, err := l.svcCtx.Cache.Get(cacheKey)
	if err == nil && len(cachedResponse) > 0 {
		var cachedResp PoolsResponseWrapper
		if err := json.Unmarshal([]byte(cachedResponse), &cachedResp); err == nil {
			// Launch a goroutine to refresh the cache in the background
			go l.refreshCache(req, cacheKey, &cachedResp)
			// Return the cached response immediately
			return &cachedResp, nil
		}
	}

	// If no valid cache, proceed to refresh and return the new response
	return l.refreshCache(req, cacheKey, nil)
}

func (l *SonexPoolListLogic) refreshCache(req *types.PoolsRequest, cacheKey string, cachedResp *PoolsResponseWrapper) (*PoolsResponseWrapper, error) {
	lockKey := fmt.Sprintf("lock:%s", cacheKey)
	ctx := context.Background()

	success, err := l.svcCtx.Cache.SetNX(lockKey, "locked", 30*time.Second)
	if err != nil {
		if cachedResp == nil {
			return nil, fmt.Errorf("failed to acquire lock: %w", err)
		}

		return cachedResp, nil
	}

	if !success {
		if cachedResp == nil {
			return nil, fmt.Errorf("cache refresh already in progress for key: %s", cacheKey)
		}

		return cachedResp, nil
	}

	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	model := repository.NewEventSonexFactoryPoolCreatedModel(conn)
	data, meta, err := model.GetPoolList(ctx, req, l.svcCtx.Config.Sonexes[0].PositionManager)
	if err != nil {
		return &PoolsResponseWrapper{
			PoolsResponse: &types.PoolsResponse{
				Code:    400,
				Data:    []*types.PoolsResponse_Pool{},
				Meta:    &types.PaginationMeta{},
				Message: err.Error(),
			},
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

			price0 := new(big.Float).SetFloat64(0)
			price1 := new(big.Float).SetFloat64(0)

			bestPrice0, _, err := global.GetPrices(ctx, nil, req.ChainId, common.HexToAddress(pool.Token0.Address).Hex(), true, 2, nil)
			if err == nil && bestPrice0 != nil {
				price0.SetInt(bestPrice0.Value)
			}

			factor0 := new(big.Float).SetFloat64(math.Pow10(int(pool.Token0.Decimals)))
			factor1 := new(big.Float).SetFloat64(math.Pow10(int(pool.Token1.Decimals)))

			poolPrice0 := new(big.Float).Quo(price0, factor0)
			pool.Price0, _ = poolPrice0.Float64()

			bestPrice1, _, err := global.GetPrices(ctx, nil, req.ChainId, common.HexToAddress(pool.Token1.Address).Hex(), true, 2, nil)
			if err == nil && bestPrice1 != nil {
				price1.SetInt(bestPrice1.Value)
			}

			poolPrice1 := new(big.Float).Quo(price1, factor1)
			pool.Price1, _ = poolPrice1.Float64()

			mu.Lock()
			token0Balance := new(big.Float).Mul(new(big.Float).SetFloat64(balance0), price0)
			token1Balance := new(big.Float).Mul(new(big.Float).SetFloat64(balance1), price1)
			pool.Token0Balance, _ = token0Balance.Float64()
			pool.Token1Balance, _ = token1Balance.Float64()
			tvlUSD := new(big.Float).Add(token0Balance, token1Balance)
			pool.TvlUSD, _ = tvlUSD.Float64()
			mu.Unlock()
		}(data[i])
	}

	wg.Wait()

	poolVolumes, err := model.GetPoolVolumes(ctx, l.svcCtx.Config.Sonexes[0].ChainId, true)
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

	response := &PoolsResponseWrapper{
		PoolsResponse: &types.PoolsResponse{
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
