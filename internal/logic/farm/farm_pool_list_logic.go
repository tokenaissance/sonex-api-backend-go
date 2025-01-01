package farm

import (
	"context"
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/lib/cache"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"
	"time"

	"sonex_api_backend_go/internal/repository"

	farmCbk "sonex_api_backend_go/internal/chain_handler/farm"

	"github.com/ahmetb/go-linq/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type FarmPoolListLogic struct {
	logx.Logger
	ctx                    context.Context
	svcCtx                 *svc.ServiceContext
	chainId                int64
	ethClient              *ethclient.Client
	cacher                 cache.ICache
	conn                   sqlx.SqlConn
	farmAddress            common.Address // UniswapV3Staker
	positionManagerAddress common.Address // NonfungiblePositionManager
}

// Farm pool list
func NewFarmPoolListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FarmPoolListLogic {

	logger := logx.WithContext(ctx).WithFields(logx.Field("module", "Farm"))
	return &FarmPoolListLogic{
		Logger:                 logger,
		ctx:                    ctx,
		svcCtx:                 svcCtx,
		chainId:                svcCtx.Config.Farm[0].ChainId,
		farmAddress:            common.HexToAddress(svcCtx.Config.Farm[0].Staker),
		positionManagerAddress: common.HexToAddress(svcCtx.Config.Sonexes[0].PositionManager),
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

	resp = &types.FarmPoolListResponse{
		List:   make([]types.FarmPoolData, 0),
		Total:  0,
		Offset: req.Offset,
		Limit:  req.Limit,
	}

	client, err := ethclient.DialContext(l.ctx, l.svcCtx.Config.GetChainConfig(req.ChainId).GetNode())
	if err != nil {
		return nil, err
	}
	l.ethClient = client
	redisCache := cache.NewRedisCache(fmt.Sprintf("%s:%d", l.svcCtx.Config.Redis.Host, l.svcCtx.Config.Redis.Port), l.svcCtx.Config.Redis.Password, l.svcCtx.Config.Redis.DB)
	l.cacher = redisCache
	l.conn = repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())

	model := repository.NewFarmPoolsModel(l.conn)
	pools, total, err := model.GetPoolList(l.ctx, req)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return
		}
		return nil, err
	}

	tokenMap := map[string]struct{}{}
	for _, pool := range pools {
		tokenMap[pool.Token0Address] = struct{}{}
		tokenMap[pool.Token1Address] = struct{}{}
	}
	tokens := []string{}
	for key, _ := range tokenMap {
		tokens = append(tokens, key)
	}

	tokenPriceModel := repository.NewSonexTokenPriceRecordModel(l.conn)
	prices, err := tokenPriceModel.GetLatestHistoryPrices(tokens)
	if err != nil {
		return nil, err
	}

	var list []types.FarmPoolData
	linq.From(pools).SelectT(func(item *repository.FarmPools) types.FarmPoolData {
		price0, ok := prices[item.Token0Address]
		if !ok {
			price0 = big.NewInt(0)
		}
		price1, ok := prices[item.Token1Address]
		if !ok {
			price1 = big.NewInt(0)
		}
		return types.FarmPoolData{
			Id:             item.Id,
			ChainId:        item.ChainId,
			PoolAddress:    item.PoolAddress,
			FeeRatio:       item.FeeRatio,
			Token0Address:  item.Token0Address,
			Token0Symbol:   item.Token0Symbol,
			Token0Decimals: item.Token0Decimals,
			Token0Price:    price0.String(),
			Token1Address:  item.Token1Address,
			Token1Symbol:   item.Token1Symbol,
			Token1Decimals: item.Token1Decimals,
			Token1Price:    price1.String(),
			Liquidity:      item.StakedLiquidity,
			LiquidityUSD:   item.StakedLiquidityUsd,
			APR:            item.Apr,
			JoinNum:        item.JoinNum,
			Status:         "active",
			CreatedAt:      item.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:      item.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		}
	}).ToSlice(&list)

	resp = &types.FarmPoolListResponse{
		List:   list,
		Total:  total,
		Offset: req.Offset,
		Limit:  req.Limit,
	}

	l.checkAndRefreshPools(req)

	return resp, nil
}

// 1小时内未更新的池子数据，异步更新
func (l *FarmPoolListLogic) checkAndRefreshPools(req *types.FarmPoolListRequest) {
	cacheKey := fmt.Sprintf("farm:pools:last_update:%d", req.ChainId)
	lastUpdate, err := l.cacher.Get(cacheKey)
	if err == redis.Nil || lastUpdate == "" || time.Since(time.Unix(cast.ToInt64(lastUpdate), 0)) > time.Hour {
		// 更新池子数据
		go func() {
			defer func() {
				if r := recover(); r != nil {
					logx.Errorf("Panic in async update pool data: %v", r)
				}
			}()

			model := repository.NewFarmPoolsModel(l.conn)
			outdatedPools, err := model.FindNotFreshPools(context.Background(), req.ChainId, time.Hour)
			if err != nil {
				logx.Errorf("Failed to get outdated pools: %v", err)
				return
			}

			if len(outdatedPools) > 0 {
				for _, pool := range outdatedPools {
					if err := farmCbk.UpdateFarmPoolData(l.svcCtx.Config, l.conn, l.ethClient, l.chainId, pool.IncentiveId); err != nil {
						logx.Errorf("Failed to update pool %s data: %v", pool.IncentiveId, err)
					}
				}
			}

			// 更新最后更新时间
			if err := l.cacher.Set(cacheKey, time.Now().Unix(), 2*time.Hour); err != nil {
				logx.Errorf("Failed to update last update time: %v", err)
			}
		}()
	}
}
