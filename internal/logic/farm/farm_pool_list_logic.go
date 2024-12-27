package farm

import (
	"context"
	"fmt"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"
	"time"

	"sonex_api_backend_go/internal/repository"

	farmCbk "sonex_api_backend_go/internal/chain_handler/callback/farm"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
)

type FarmPoolListLogic struct {
	logx.Logger
	ctx             context.Context
	svcCtx          *svc.ServiceContext
	chainId         int64
	ethClient       *ethclient.Client
	farmAddress     common.Address // UniswapV3Staker
	positionAddress common.Address // NonfungiblePositionManager
}

// Farm pool list
func NewFarmPoolListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FarmPoolListLogic {
	client, err := ethclient.DialContext(ctx, provider)
	if err != nil {
		logx.Errorf("[Farm] failed to connect to the Ethereum client: %v", err)
		return nil
	}
	return &FarmPoolListLogic{
		Logger:          logx.WithContext(ctx),
		ctx:             ctx,
		svcCtx:          svcCtx,
		chainId:         svcCtx.Config.Farm[0].ChainId,
		ethClient:       client,
		farmAddress:     common.HexToAddress(svcCtx.Config.Farm[0].Staker),
		positionAddress: common.HexToAddress(svcCtx.Config.Sonexes[0].PositionManager),
	}
}

func (l *FarmPoolListLogic) FarmPoolList(req *types.FarmPoolListRequest) (resp *types.FarmPoolListResponse, err error) {

	if req.ChainId == 0 {
		req.ChainId = l.chainId
	}

	if req.UserAddress != "" && !common.IsHexAddress(req.UserAddress) {
		return nil, fmt.Errorf("invalid parameter: user_address")
	}

	if req.Limit < 0 {
		req.Limit = 0
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	model := repository.NewFarmPoolsModel(l.svcCtx.DB)
	pools, total, err := model.GetPoolList(l.ctx, req)
	if err != nil {
		return nil, err
	}

	resp = &types.FarmPoolListResponse{
		List:   pools,
		Total:  total,
		Offset: req.Offset,
		Limit:  req.Limit,
	}

	// l.checkAndRefreshPools(req)

	return resp, nil
}

// 1小时内未更新的池子数据，异步更新
func (l *FarmPoolListLogic) checkAndRefreshPools(req *types.FarmPoolListRequest) {
	cacheKey := fmt.Sprintf("farm:pools:last_update:%d", req.ChainId)
	lastUpdate, err := l.svcCtx.Cache.Get(cacheKey)
	if err == redis.Nil || lastUpdate == "" || time.Since(time.Unix(cast.ToInt64(lastUpdate), 0)) > time.Hour {
		// 更新池子数据
		go func() {
			defer func() {
				if r := recover(); r != nil {
					logx.Errorf("Panic in async update pool data: %v", r)
				}
			}()

			model := repository.NewFarmPoolsModel(l.svcCtx.DB)
			outdatedPools, err := model.FindNotFreshPools(context.Background(), req.ChainId, time.Hour)
			if err != nil {
				logx.Errorf("Failed to get outdated pools: %v", err)
				return
			}

			if len(outdatedPools) > 0 {
				farmCallback := farmCbk.NewFarmCallback(l.svcCtx)
				for _, pool := range outdatedPools {
					if err := farmCallback.UpdateFarmPoolData(l.ethClient, pool.IncentiveId); err != nil {
						logx.Errorf("Failed to update pool %s data: %v", pool.IncentiveId, err)
					}
				}
			}

			// 更新最后更新时间
			if err := l.svcCtx.Cache.Set(cacheKey, time.Now().Unix(), 2*time.Hour); err != nil {
				logx.Errorf("Failed to update last update time: %v", err)
			}
		}()
	}
}
