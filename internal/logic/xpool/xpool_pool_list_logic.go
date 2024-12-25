package xpool

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"sonex_api_backend_go/internal/module/xpool/utils"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/ahmetb/go-linq/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/sync/errgroup"
)

type XpoolPoolListLogic struct {
	logx.Logger
	ctx            context.Context
	svcCtx         *svc.ServiceContext
	chainId        int64
	factoryAddress common.Address
	ethClient      *ethclient.Client
}

func NewXpoolPoolListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *XpoolPoolListLogic {
	client, err := ethclient.DialContext(ctx, provider)
	if err != nil {
		logx.Errorf("[Xpool] failed to connect to the Ethereum client: %v", err)
		return nil
	}
	return &XpoolPoolListLogic{
		Logger:         logx.WithContext(ctx),
		ctx:            ctx,
		svcCtx:         svcCtx,
		chainId:        svcCtx.Config.Xpools[0].ChainId,
		factoryAddress: common.HexToAddress(svcCtx.Config.Xpools[0].Factory),
		ethClient:      client,
	}
}

func (l *XpoolPoolListLogic) XpoolPoolList(req *types.XpoolPoolListRequest) (resp []*types.XpoolPoolListResponse, total int64, err error) {
	if req.ChainId == 0 {
		req.ChainId = l.chainId
	}
	if req.UserAddress != "" && !common.IsHexAddress(req.UserAddress) {
		return nil, 0, fmt.Errorf("invalid parameter: user_address")
	}
	if req.Limit < 1 {
		req.Limit = 10
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	// 检查缓存中的池子列表
	if err := l.checkAndUpdatePoolListCache(); err != nil {
		logx.Errorf("failed to check and update pool list cache: %v", err)
	}

	// 更新池子的实时信息
	if err := l.updatePoolRealTimeInfo(); err != nil {
		logx.Errorf("failed to update pool real-time info: %v", err)
	}

	// 如果需要,更新用户质押池数据和资产
	if req.UserAddress != "" {
		if err := l.updateUserPoolsData(req.UserAddress); err != nil {
			logx.Errorf("failed to update user data: %v", err)
		}
	}

	if req.StakedOnly == "true" && req.UserAddress != "" {
		// 从缓存中获取用户的质押池列表
		poolListKey := fmt.Sprintf("sonex_dex:xpool:user_pools:%s", req.UserAddress)
		cachedPoolAddresses, err := l.svcCtx.Cache.SMembers(l.ctx, poolListKey).Result()
		if err == nil && len(cachedPoolAddresses) > 0 {
			req.UserPools = cachedPoolAddresses
		} else {
			req.UserPools = []string{}
		}
	}

	// 获取池子列表
	model := repository.NewXpoolPoolsModel(l.svcCtx.DB)
	results, total, err := model.GetPoolList(l.ctx, req)
	if err != nil {
		return nil, 0, err
	}
	linq.From(results).SelectT(func(item repository.XpoolPools) types.XpoolPoolListResponse {
		updatedAt := ""
		if item.UpdatedAt.Valid {
			updatedAt = item.UpdatedAt.Time.Format("2006-01-02 15:04:05")
		}
		return types.XpoolPoolListResponse{
			Id:                   item.Id,
			ChainId:              item.ChainId,
			StakingTokenName:     item.StakingTokenName,
			StakingTokenSymbol:   item.StakingTokenSymbol,
			StakingTokenAddress:  item.StakingTokenAddress,
			StakingTokenDecimals: item.StakingTokenDecimals,
			StakingTokenPrice:    item.StakingTokenPrice,
			RewardTokenName:      item.RewardTokenName,
			RewardTokenSymbol:    item.RewardTokenSymbol,
			RewardTokenAddress:   item.RewardTokenAddress,
			RewardTokenDecimals:  item.RewardTokenDecimals,
			RewardTokenPrice:     item.RewardTokenPrice,
			PoolAddress:          item.PoolAddress,
			RewardPerBlock:       item.RewardPerBlock,
			TotalStaked:          item.TotalStaked,
			BlockTime:            item.BlockTime,
			TVL:                  item.Tvl,
			APR:                  item.Apr,
			JoinNum:              item.JoinNum,
			Status:               item.Status,
			CreatedAt:            item.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:            updatedAt,
		}
	}).ToSlice(&resp)

	return resp, total, nil
}

// checkAndUpdatePoolListCache 检查并更新池子列表缓存
func (l *XpoolPoolListLogic) checkAndUpdatePoolListCache() error {
	const cacheKeyPrefix = "sonex_dex:xpool:"
	poolAddressesKey := cacheKeyPrefix + "pool_addresses"
	cachedPoolAddresses, err := l.svcCtx.Cache.SMembers(l.ctx, poolAddressesKey).Result()
	if err != nil || len(cachedPoolAddresses) == 0 {
		if err := l.updatePoolBasicInfo(); err != nil {
			return fmt.Errorf("failed to update pool basic info: %v", err)
		}
	} else {
		// 检查是否需要增量更新
		lastUpdateKey := cacheKeyPrefix + "last_update_time_basic"
		lastUpdateStr, err := l.svcCtx.Cache.Get(lastUpdateKey)
		if err == nil {
			lastUpdate, err := time.Parse(time.RFC3339, lastUpdateStr)
			if err == nil && time.Since(lastUpdate) >= 30*time.Minute {
				if err := l.updatePoolBasicInfo(); err != nil {
					return fmt.Errorf("failed to update pool basic info: %v", err)
				}
			}
		}
	}
	return nil
}

// updatePoolBasicInfo 更新池子的基本信息
func (l *XpoolPoolListLogic) updatePoolBasicInfo() error {
	const updateInterval = 30 * time.Minute
	const cacheKeyPrefix = "sonex_dex:xpool:"

	// 检查是否需要更新
	lastUpdateKey := cacheKeyPrefix + "last_update_time_basic"
	lastUpdateStr, err := l.svcCtx.Cache.Get(lastUpdateKey)
	if err == nil {
		lastUpdate, err := time.Parse(time.RFC3339, lastUpdateStr)
		if err == nil && time.Since(lastUpdate) < updateInterval {
			return nil
		}
	}

	factory, err := getFactoryContract(l.ctx, l.ethClient, l.factoryAddress)
	if err != nil {
		return fmt.Errorf("failed to get staking factory: %v", err)
	}

	poolAddresses, err := factory.GetAllPools(nil)
	if err != nil {
		return fmt.Errorf("failed to get all pools: %v", err)
	}

	poolAddressesKey := cacheKeyPrefix + "pool_addresses"
	cachedPoolAddresses, err := l.svcCtx.Cache.SMembers(l.ctx, poolAddressesKey).Result()
	if err == nil && len(cachedPoolAddresses) > 0 {
		newPoolAddresses := findNewAddresses(cachedPoolAddresses, poolAddresses)
		if len(newPoolAddresses) == 0 {
			return nil
		}
		poolAddresses = newPoolAddresses
	}

	// 增量更新池子的基本信息
	var wg sync.WaitGroup
	pools := make([]*repository.XpoolPools, len(poolAddresses))
	for i, poolAddress := range poolAddresses {
		wg.Add(1)
		go func(i int, poolAddress common.Address) {
			defer wg.Done()
			pool, err := l.getPoolInfo(l.ethClient, poolAddress)
			if err != nil {
				logx.Errorf("[Xpool] failed to get pool full info for %s: %v", poolAddress.Hex(), err)
				return
			}
			pools[i] = pool
		}(i, poolAddress)
	}
	wg.Wait()

	if len(pools) > 0 {
		model := repository.NewXpoolPoolsModel(l.svcCtx.DB)
		if err := model.BatchUpdatePools(l.ctx, pools); err != nil {
			return fmt.Errorf("[Xpool] failed to update pools: %v", err)
		}
	}

	poolAddressesStr := make([]string, 0, len(poolAddresses))
	for _, addr := range poolAddresses {
		poolAddressesStr = append(poolAddressesStr, addr.Hex())
	}

	l.svcCtx.Cache.Set(lastUpdateKey, time.Now().Format(time.RFC3339), 0)
	l.svcCtx.Cache.SAdd(l.ctx, poolAddressesKey, poolAddressesStr)

	return nil
}

// updatePoolRealTimeInfo 更新池子的实时信息
func (l *XpoolPoolListLogic) updatePoolRealTimeInfo() error {
	const updateInterval = 2 * time.Minute
	const cacheKeyPrefix = "sonex_dex:xpool:"

	// 检查是否需要更新
	lastUpdateKey := cacheKeyPrefix + "last_update_time_realtime"
	lastUpdateStr, err := l.svcCtx.Cache.Get(lastUpdateKey)
	if err == nil {
		lastUpdate, err := time.Parse(time.RFC3339, lastUpdateStr)
		if err == nil && time.Since(lastUpdate) < updateInterval {
			return nil
		}
	}

	// 从缓存中获取池子列表
	poolAddressesKey := cacheKeyPrefix + "pool_addresses"
	cachedPoolAddresses, err := l.svcCtx.Cache.SMembers(l.ctx, poolAddressesKey).Result()
	if err != nil {
		return fmt.Errorf("[Xpool] failed to get cached pool addresses: %v", err)
	}

	// 更新池子的实时信息
	var g errgroup.Group
	pools := make([]*repository.XpoolPools, len(cachedPoolAddresses))
	for i, poolAddress := range cachedPoolAddresses {
		i, poolAddress := i, poolAddress
		g.Go(func() error {
			pool, err := l.getPoolRealTimeInfo(l.ethClient, common.HexToAddress(poolAddress))
			if err != nil {
				logx.Errorf("failed to get pool real-time info for %s: %v", poolAddress, err)
				return err
			}
			pools[i] = pool
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}

	// 批量更新数据库
	if len(pools) > 0 {
		model := repository.NewXpoolPoolsModel(l.svcCtx.DB)
		if err := model.BatchUpdatePools(l.ctx, pools); err != nil {
			return fmt.Errorf("failed to update pools: %v", err)
		}
	}

	// 更新缓存时间
	l.svcCtx.Cache.Set(lastUpdateKey, time.Now().Format(time.RFC3339), 0)

	return nil
}

// getPoolInfo 获取池子的信息
func (l *XpoolPoolListLogic) getPoolInfo(ethClient *ethclient.Client, poolAddress common.Address) (*repository.XpoolPools, error) {
	factory, err := getFactoryContract(l.ctx, ethClient, l.factoryAddress)
	if err != nil {
		return nil, err
	}

	info, err := factory.GetPoolInfo(nil, poolAddress)
	if err != nil {
		return nil, err
	}

	// 获取代币价格
	prices, err := getTokenPrices(l.ctx, l.chainId, info.StakingToken, info.RewardToken)
	if err != nil {
		return nil, err
	}

	// 计算TVL
	stakingTokenPrice := decimal.NewFromBigInt(prices.StakingTokenPrice, 0)
	tvl := decimal.NewFromBigInt(big.NewInt(0), 0).Mul(stakingTokenPrice)

	return &repository.XpoolPools{
		ChainId:              int64(l.chainId),
		StakingTokenName:     info.StakingTokenName,
		StakingTokenSymbol:   info.StakingTokenSymbol,
		StakingTokenAddress:  info.StakingToken.Hex(),
		StakingTokenDecimals: int64(info.StakingTokenDecimals),
		StakingTokenPrice:    stakingTokenPrice.String(),
		RewardTokenName:      info.RewardTokenName,
		RewardTokenSymbol:    info.RewardTokenSymbol,
		RewardTokenAddress:   info.RewardToken.Hex(),
		RewardTokenDecimals:  int64(info.RewardTokenDecimals),
		RewardTokenPrice:     decimal.NewFromBigInt(prices.RewardTokenPrice, 0).String(),
		PoolAddress:          poolAddress.Hex(),
		RewardPerBlock:       info.RewardPerBlock.String(),
		BlockTime:            blockTime,
		Tvl:                  tvl.String(),
		Status:               "active",
	}, nil
}

// getPoolRealTimeInfo 获取池子的实时信息
func (l *XpoolPoolListLogic) getPoolRealTimeInfo(ethClient *ethclient.Client, poolAddress common.Address) (*repository.XpoolPools, error) {
	pool, err := getPoolContract(l.ctx, ethClient, poolAddress)
	if err != nil {
		return nil, err
	}

	// 实时数据
	totalStaked, err := pool.TotalStaked(nil)
	if err != nil {
		return nil, err
	}

	claimed, err := pool.TotalClaimed(nil)
	if err != nil {
		return nil, err
	}

	joinNum, err := pool.StakedNumber(nil)
	if err != nil {
		return nil, err
	}

	factory, err := getFactoryContract(l.ctx, ethClient, l.factoryAddress)
	if err != nil {
		return nil, err
	}

	poolInfo, err := factory.GetPoolInfo(nil, poolAddress)
	if err != nil {
		return nil, err
	}

	// 获取代币价格
	prices, err := getTokenPrices(l.ctx, l.chainId, poolInfo.StakingToken, poolInfo.RewardToken)
	if err != nil {
		return nil, err
	}

	// 计算TVL和APR
	stakingTokenPrice := decimal.NewFromBigInt(prices.StakingTokenPrice, 0)
	tvl := decimal.NewFromBigInt(totalStaked, 0).Mul(stakingTokenPrice)

	apr, _ := utils.CalculateAPR(
		prices.RewardTokenPrice,
		big.NewInt(0),
		blockTime,
		prices.StakingTokenPrice,
		totalStaked,
	)

	return &repository.XpoolPools{
		TotalStaked:  totalStaked.String(),
		TotalClaimed: claimed.TotalClaimed.String(),
		Tvl:          tvl.String(),
		Apr:          apr.String(),
		JoinNum:      joinNum.Int64(),
	}, nil
}

func (l *XpoolPoolListLogic) updateUserPoolsData(userAddr string) error {
	const cacheKeyPrefix = "sonex_dex:xpool:"
	const cacheExpiration = 2 * time.Minute

	// 检查缓存中是否存在用户的质押池数据
	poolListKey := fmt.Sprintf("%suser_pools:%s", cacheKeyPrefix, userAddr)
	cachedPoolAddresses, err := l.svcCtx.Cache.SMembers(l.ctx, poolListKey).Result()
	if err == nil && len(cachedPoolAddresses) > 0 {
		// 检查更新时间
		lastUpdateKey := fmt.Sprintf("%suser_pools_last_update:%s", cacheKeyPrefix, userAddr)
		lastUpdateStr, err := l.svcCtx.Cache.Get(lastUpdateKey)
		if err == nil {
			lastUpdate, err := time.Parse(time.RFC3339, lastUpdateStr)
			if err == nil && time.Since(lastUpdate) < cacheExpiration {
				return nil
			}
		}
	}

	factory, err := getFactoryContract(l.ctx, l.ethClient, l.factoryAddress)
	if err != nil {
		return err
	}

	pools, err := factory.GetPools(nil, common.HexToAddress(userAddr))
	if err != nil {
		return err
	}

	poolAddrs := make([]string, len(pools))
	for i, addr := range pools {
		poolAddrs[i] = addr.Hex()
	}
	l.svcCtx.Cache.SAdd(l.ctx, poolListKey, poolAddrs)

	// 更新用户质押池数据的更新时间
	lastUpdateKey := fmt.Sprintf("%suser_pools_last_update:%s", cacheKeyPrefix, userAddr)
	l.svcCtx.Cache.Set(lastUpdateKey, time.Now().Format(time.RFC3339), cacheExpiration)

	// 计算用户总资产
	totalStaked := decimal.NewFromInt(0)
	totalEarned := decimal.NewFromInt(0)

	for _, pool := range pools {
		stakingPool, err := getPoolContract(l.ctx, l.ethClient, pool)
		if err != nil {
			continue
		}

		// 获取质押金额
		staked, err := stakingPool.Staked(nil, common.HexToAddress(userAddr))
		if err != nil {
			continue
		}
		totalStaked = totalStaked.Add(decimal.NewFromBigInt(staked, 0))

		// 获取未领取收益
		unclaimed, err := stakingPool.Unclaimed(nil)
		if err != nil {
			continue
		}
		totalEarned = totalEarned.Add(decimal.NewFromBigInt(unclaimed.Unclaimed, 0))
	}

	// 缓存用户资产数据
	assetsKey := fmt.Sprintf("%suser_assets:%s", cacheKeyPrefix, userAddr)
	l.svcCtx.Cache.HSet(assetsKey, "balance", totalStaked.Add(totalEarned).String())
	l.svcCtx.Cache.HSet(assetsKey, "staked", totalStaked.String())
	l.svcCtx.Cache.HSet(assetsKey, "earned", totalEarned.String())
	l.svcCtx.Cache.Expire(assetsKey, cacheExpiration)

	return nil
}

func findNewAddresses(cachedAddresses []string, newAddresses []common.Address) []common.Address {
	cachedMap := make(map[string]struct{})
	for _, addr := range cachedAddresses {
		cachedMap[addr] = struct{}{}
	}

	var result []common.Address
	for _, addr := range newAddresses {
		if _, exists := cachedMap[addr.Hex()]; !exists {
			result = append(result, addr)
		}
	}
	return result
}
