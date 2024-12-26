package sonex

import (
	"context"
	"encoding/json"
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

			price0 := new(big.Float).SetFloat64(0)
			price1 := new(big.Float).SetFloat64(0)

			bestPrice0, _, err := global.GetPrices(l.ctx, nil, req.ChainId, common.HexToAddress(pool.Token0.Address).Hex(), true, 2, nil)
			//			fmt.Println("bestPrice0", pool.Token0.Address, bestPrice0.Value.String())
			if err == nil && bestPrice0 != nil {
				price0.SetInt(bestPrice0.Value)
			}

			factor0 := new(big.Float).SetFloat64(math.Pow10(int(pool.Token0.Decimals)))
			factor1 := new(big.Float).SetFloat64(math.Pow10(int(pool.Token1.Decimals)))

			poolPrice0 := new(big.Float).Quo(price0, factor0)
			pool.Price0, _ = poolPrice0.Float64() // Convert to float64 for assignment

			bestPrice1, _, err := global.GetPrices(l.ctx, nil, req.ChainId, common.HexToAddress(pool.Token1.Address).Hex(), true, 2, nil)
			//fmt.Println("bestPrice1", pool.Token1.Address, bestPrice1.Value.String())
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
