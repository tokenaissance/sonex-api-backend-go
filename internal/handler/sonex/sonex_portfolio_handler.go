package sonex

import (
	"context"
	"errors"
	"math"
	"net/http"
	"sonex_api_backend_go/internal/logic/farm"
	"sonex_api_backend_go/internal/logic/sonex"
	"sonex_api_backend_go/internal/logic/xpool"
	sonex_module "sonex_api_backend_go/internal/module/sonex"
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

		portfolioType := req.Type

		response := &types.PortfolioResponse{}
		positions := make([]*types.Position, 0)

		// 获取用户有效 positions 列表
		{
			if portfolioType == "liquidity" || portfolioType == "all" {
				res, err := getPositions(r.Context(), svcCtx, req)
				if err != nil {
					logx.Errorf("[SonexPortfolioHandler]: SonexPositions error: %s, chainId: %d, address: %s, type: %s", err.Error(), req.ChainId, req.Address, req.Type)
					httpx.ErrorCtx(r.Context(), w, err)
					return
				}
				positions = res
				response.Positions = res
				// count in total position & liquidity claimable fees
				for _, position := range positions {
					response.TotalPosition += position.Position
					response.LiquidityClaimableFees += position.ClaimableFees
				}
			}
		}

		// 获取用户有效 xpool 列表
		{
			if portfolioType == "xpool" || portfolioType == "all" {
				xpools, err := getXPools(r.Context(), svcCtx, req)
				if err != nil {
					logx.Errorf("[SonexPortfolioHandler]: getXPools error: %s, chainId: %d, address: %s, type: %s", err.Error(), req.ChainId, req.Address, req.Type)
					httpx.ErrorCtx(r.Context(), w, err)
					return
				}
				pools := make([]*types.PortfolioXpool, 0)
				for _, xpool := range xpools {
					stakingTokenPrice, err := sonex_module.NewTokenPrice(xpool.StakingTokenPriceWithDecimals, int(xpool.StakingTokenPriceDecimals))
					if err != nil {
						logx.Errorf("[SonexPortfolioHandler]: NewTokenPrice error: %s, chainId: %d, address: %s, type: %s", err.Error(), req.ChainId, req.Address, req.Type)
						continue
					}
					stakingTokenValue, err := sonex_module.NewTokenValue(xpool.Staked, int(xpool.StakingTokenDecimals), sonex_module.WithTokenPrice(stakingTokenPrice))
					if err != nil {
						logx.Errorf("[SonexPortfolioHandler]: NewTokenValue error: %s, chainId: %d, address: %s, type: %s", err.Error(), req.ChainId, req.Address, req.Type)
						continue
					}
					pools = append(pools, &types.PortfolioXpool{
						PoolAddress:         xpool.PoolAddress,               // 质押池合约地址
						StakingTokenName:    xpool.StakingTokenName,          // 质押的代币名称
						StakingTokenSymbol:  xpool.StakingTokenSymbol,        // 质押的代币符号
						StakingTokenAddress: xpool.StakingTokenAddress,       // 质押的代币合约地址
						APR:                 xpool.APR,                       // 年化收益率
						Staked:              stakingTokenValue.ValueString(), // 用户已经质押的美元价值
						// FeeRevenue:          xpool.Earned,
					})
					xpoolStaked, _ := stakingTokenValue.ValueFloat64()
					response.TotalXpoolStaked += xpoolStaked
					// response.TotalXpoolEarned += xpool.Earned // todo: xpool earned 需要从 xpool 的 event 事件中获取
				}
				response.Xpools = pools
			}
		}

		// 获取用户有效 farm 列表
		{
			if portfolioType == "farm" || portfolioType == "all" {
				farmPools, err := getFarmPools(r.Context(), svcCtx, req, positions)
				if err != nil {
					logx.Errorf("[SonexPortfolioHandler]: getFarmPools error: %s, chainId: %d, address: %s, type: %s", err.Error(), req.ChainId, req.Address, req.Type)
					httpx.ErrorCtx(r.Context(), w, err)
					return
				}
				pools := make([]*types.PortfolioFarm, 0)
				for _, farm := range farmPools {
					pools = append(pools, &types.PortfolioFarm{
						PoolAddress:   farm.PoolAddress,
						FeeRatio:      farm.FeeRatio,
						Token0Symbol:  farm.Token0Symbol,
						Token0Address: farm.Token0Address,
						Token1Symbol:  farm.Token1Symbol,
						Token1Address: farm.Token1Address,
						PositionValue: farm.LiquidityUSD,
						FeeRevenue:    farm.TotalClaimedFeesUSD,
						Earned:        farm.TotalClaimableFeesUSD,
					})
					response.TotalFarmStaked += farm.LiquidityUSDFloat64
					response.FarmEarned += farm.TotalClaimableFeesUSDFloat64
				}
				response.FarmPools = pools
			}
		}

		for _, position := range positions {
			response.TotalPosition += position.Position
			response.ClaimableFees += position.ClaimableFees
		}

		// 计算总资产和可提资产
		response.TotalAssets = response.TotalPosition + response.TotalXpoolStaked + response.TotalFarmStaked
		response.ClaimableFees = response.LiquidityClaimableFees + response.FarmEarned + response.XpoolEarned
		response.TotalAssetsRate = "0.00%"

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
