package farm

import (
	"context"
	"strconv"
	"time"

	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"sonex_api_backend_go/internal/repository"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

type FarmPoolListLogic struct {
	logx.Logger
	ctx             context.Context
	svcCtx          *svc.ServiceContext
	chainId         uint64
	ethClient       *ethclient.Client
	factoryAddress  common.Address // UniswapV3Factory
	positionAddress common.Address // NonfungiblePositionManager
}

// Farm pool list
func NewFarmPoolListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FarmPoolListLogic {
	client, err := ethclient.DialContext(ctx, provider)
	if err != nil {
		logx.Errorf("[Xpool] failed to connect to the Ethereum client: %v", err)
		return nil
	}
	return &FarmPoolListLogic{
		Logger:          logx.WithContext(ctx),
		ctx:             ctx,
		svcCtx:          svcCtx,
		chainId:         uint64(svcCtx.Config.Xpools[0].ChainId),
		ethClient:       client,
		factoryAddress:  common.HexToAddress(svcCtx.Config.Sonexes[0].Factory),
		positionAddress: common.HexToAddress(svcCtx.Config.Sonexes[0].PositionManager),
	}
}

func (l *FarmPoolListLogic) FarmPoolList(req *types.FarmPoolListRequest) (resp *types.FarmPoolListResponse, err error) {
	// 获取用户参与的池子ID列表
	if req.StakedOnly == "true" && req.UserAddress != "" {
		userPools, err := l.getUserPools(req.UserAddress)
		if err != nil {
			return nil, err
		}
		req.UserPools = userPools
	}

	// 定期增量聚合数据到farm_pools表
	err = l.aggregateFarmPoolsData()
	if err != nil {
		return nil, err
	}

	// 获取池子列表
	model := repository.NewFarmPoolsModel(l.svcCtx.DB)
	pools, total, err := model.GetPoolList(l.ctx, req.ChainId, req.Keyword, req.UserPools, req.OrderBy, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	resp = &types.FarmPoolListResponse{
		List:   make([]types.FarmPoolData, len(pools)),
		Total:  total,
		Offset: int(req.Offset),
		Limit:  int(req.Limit),
	}

	resp.List = pools

	return resp, nil
}

func (l *FarmPoolListLogic) aggregateFarmPoolsData() error {
	const cacheKeyPrefix = "sonex_dex:farm:"
	const updateInterval = 5 * time.Minute

	// 检查是否需要更新
	lastUpdateKey := cacheKeyPrefix + "last_update_time"
	lastUpdateStr, err := l.svcCtx.Cache.Get(lastUpdateKey)
	if err == nil {
		lastUpdate, err := time.Parse(time.RFC3339, lastUpdateStr)
		if err == nil && time.Since(lastUpdate) < updateInterval {
			return nil
		}
	}

	// 获取event_sonex_factory_pool_created表的数据
	eventModel := repository.NewEventSonexFactoryPoolCreatedModel(l.svcCtx.DB)
	poolList, err := eventModel.FindAll(l.ctx)
	if err != nil {
		return err
	}

	// 获取token表的数据
	tokenModel := repository.NewSonexTokenModel(l.svcCtx.DB)
	tokens, err := tokenModel.FindAllToMap(l.ctx)
	if err != nil {
		return err
	}

	// 获取缓存中的pool_address集合
	poolAddressKey := cacheKeyPrefix + "pool_addresses"
	cachedPoolAddresses, err := l.svcCtx.Cache.SMembers(l.ctx, poolAddressKey).Result()
	if err != nil && err != redis.Nil {
		return err
	}
	cachedPoolAddressSet := make(map[string]struct{})
	for _, addr := range cachedPoolAddresses {
		cachedPoolAddressSet[addr] = struct{}{}
	}

	// 聚合数据
	farmPoolsModel := repository.NewFarmPoolsModel(l.svcCtx.DB)
	for _, pool := range poolList {
		if _, exists := cachedPoolAddressSet[pool.EventPool]; exists {
			continue
		}

		token0 := tokens[pool.EventToken0]
		token1 := tokens[pool.EventToken1]

		// 获取实时数据
		liquidityUSD, err := getLiquidityUSD(l.ctx, l.ethClient, l.chainId, pool.EventPool)
		if err != nil {
			return err
		}
		apr, err := getAPR(l.ctx, l.ethClient, pool.EventPool)
		if err != nil {
			return err
		}
		joinNum, err := getJoinNum(l.ctx, l.ethClient, pool.EventPool)
		if err != nil {
			return err
		}

		// 更新或插入数据到farm_pools表
		token0Symbol := ""
		if token0.Symbol.Valid {
			token0Symbol = token0.Symbol.String
		}
		token1Symbol := ""
		if token1.Symbol.Valid {
			token1Symbol = token1.Symbol.String
		}
		err = farmPoolsModel.Upsert(l.ctx, &repository.FarmPools{
			PoolAddress:    pool.EventPool,
			Token0Address:  pool.EventToken0,
			Token1Address:  pool.EventToken1,
			FeeRatio:       strconv.Itoa(int(pool.EventFee)),
			TickSpacing:    pool.EventTickSpacing,
			Token0Symbol:   token0Symbol,
			Token1Symbol:   token1Symbol,
			Token0Decimals: token0.Decimals,
			Token1Decimals: token1.Decimals,
			LiquidityUsd:   liquidityUSD,
			Apr:            apr,
			JoinNum:        int64(joinNum),
		})
		if err != nil {
			return err
		}

		// 将新的pool_address添加到缓存中
		l.svcCtx.Cache.SAdd(l.ctx, poolAddressKey, pool.EventPool)
	}

	// 更新缓存时间
	l.svcCtx.Cache.Set(l.ctx, lastUpdateKey, time.Now().Format(time.RFC3339), updateInterval)

	return nil
}

func getJoinNum(ctx context.Context, ethClient *ethclient.Client, poolAddress string) (int, error) {
	// ...实现获取参与人数的逻辑...
	return 0, nil
}

// 获取用户参与的池子ID列表
func (l *FarmPoolListLogic) getUserPools(userAddress string) ([]string, error) {
	// 简略逻辑，假设我们从数据库或缓存中获取用户参与的池子ID列表
	// 这里返回一个空的slice作为示例
	return []string{}, nil
}
