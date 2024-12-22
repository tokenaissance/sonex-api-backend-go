package sonex

import (
	"context"
	"encoding/json"
	"math"
	"time"

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
