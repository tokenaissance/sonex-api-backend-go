package sonex

import (
	"context"
	"errors"
	"math"
	"net/http"
	"sonex_api_backend_go/internal/logic/farm"
	"sonex_api_backend_go/internal/logic/sonex"
	"sonex_api_backend_go/internal/logic/xpool"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SonexPortfolioHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PortfolioRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := sonex.NewSonexPositionsLogic(r.Context(), svcCtx)
		positions, err := l.SonexPositions(r.Context(), req.ChainId, req.Address)
		if err != nil {
			logx.Errorf("[SonexPortfolioHandler]: SonexPositions error: %s, chainId: %d, address: %s", err.Error(), req.ChainId, req.Address)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		response := &types.PortfolioResponse{}
		for _, position := range positions {
			response.TotalPosition += position.Position
			response.ClaimableFees += position.ClaimableFees
		}

		httpx.OkJsonCtx(r.Context(), w, response)
	}
}

// getPositions 获取用户有效 positions 列表
func getPositions(ctx context.Context, svcCtx *svc.ServiceContext, req types.PortfolioRequest) ([]*types.Position, error) {
	l := sonex.NewSonexPositionsLogic(ctx, svcCtx)
	positions, err := l.SonexPositions(ctx, req.ChainId, req.Address)
	if err != nil {
		logx.Errorf("[getPositions]: SonexPositions error: %s, chainId: %d, address: %s, type: %s", err.Error(), req.ChainId, req.Address, req.Type)
		return nil, err
	}

	return positions, nil
}

// getXPools 获取用户有效 xpool 列表
func getXPools(ctx context.Context, svcCtx *svc.ServiceContext, req types.PortfolioRequest) ([]types.XpoolPoolListResponse_Pool, error) {
	l := xpool.NewXpoolPoolListLogic(ctx, svcCtx)
	pageIndex := int64(1)
	pageSize := int64(math.MaxInt64)
	pools, err := l.XpoolPoolList(&types.XpoolPoolListRequest{
		ChainRequest: types.ChainRequest{
			ChainId: req.ChainId,
		},
		UserAddress: req.Address,
		PageRequest: types.PageRequest{
			PageIndex: &pageIndex,
			PageSize:  &pageSize,
		},
	})
	if err != nil {
		logx.Errorf("[getXPools]: XpoolPoolList error: %s, chainId: %d, address: %s, type: %s", err.Error(), req.ChainId, req.Address, req.Type)
		return nil, err
	}

	list, ok := pools.Data.([]types.XpoolPoolListResponse_Pool)
	if !ok {
		logx.Errorf("[getXPools]: XpoolPoolList data type error, chainId: %d, address: %s, type: %s", req.ChainId, req.Address, req.Type)
		return nil, errors.New("XpoolPoolList data type error")
	}

	return list, nil
}

// getFarmPools 获取用户有效 farm 列表
func getFarmPools(ctx context.Context, svcCtx *svc.ServiceContext, req types.PortfolioRequest, positions []*types.Position) ([]types.FarmIncentivesResponse_Data, error) {
	l := farm.NewFarmIncentivesLogic(ctx, svcCtx)

	pools, err := l.FarmIncentives(&types.FarmIncentivesRequest{
		ChainRequest: types.ChainRequest{
			ChainId: req.ChainId,
		},
		Status:      "active",
		UserAddress: req.Address,
		StakedOnly:  true,
	}, positions)
	if err != nil {
		logx.Errorf("[getFarmPools]: FarmIncentives error: %s, chainId: %d, address: %s, type: %s", err.Error(), req.ChainId, req.Address, req.Type)
		return nil, err
	}

	return pools.List, nil
}
