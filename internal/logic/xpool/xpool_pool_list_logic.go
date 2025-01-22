package xpool

import (
	"context"
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/handler/response"
	"sonex_api_backend_go/internal/module/xpool"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type XpoolPoolListLogic struct {
	logx.Logger
	ctx                               context.Context
	svcCtx                            *svc.ServiceContext
	eventXpoolFactoryPoolCreatedModel repository.EventXpoolFactoryPoolCreatedModel
	tokenPriceModel                   repository.SonexTokenPriceRecordModel
}

type CacheWarpper struct {
	List       []types.XpoolPoolListResponse_Pool
	LastUpdate time.Time
}

type GetPoolInfoResult struct {
	StakingToken         common.Address `json:"stakingToken"`
	StakingTokenName     common.Hash    `json:"stakingTokenName"`
	StakingTokenSymbol   common.Hash    `json:"stakingTokenSymbol"`
	StakingTokenDecimals uint8          `json:"stakingTokenDecimals"`
	RewardToken          common.Address `json:"rewardToken"`
	RewardTokenName      common.Hash    `json:"rewardTokenName"`
	RewardTokenSymbol    common.Hash    `json:"rewardTokenSymbol"`
	RewardTokenDecimals  uint8          `json:"rewardTokenDecimals"`
	RewardPerBlock       *big.Int       `json:"rewardPerBlock"`
}

// Xpool list
func NewXpoolPoolListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *XpoolPoolListLogic {
	conn := repository.NewConn(svcCtx.Config.PostgreSQL.DateSource())
	return &XpoolPoolListLogic{
		Logger:                            logx.WithContext(ctx),
		ctx:                               ctx,
		svcCtx:                            svcCtx,
		eventXpoolFactoryPoolCreatedModel: repository.NewEventXpoolFactoryPoolCreatedModel(conn),
		tokenPriceModel:                   repository.NewSonexTokenPriceRecordModel(conn),
	}
}

func (l *XpoolPoolListLogic) XpoolPoolListBack(req *types.XpoolPoolListRequest) (list []types.XpoolPoolListResponse_Pool, err error) {
	if req.ChainId == 0 {
		return nil, fmt.Errorf("%s", "Invalid parameter")
	}
	// if req.PageSize < 1 {
	// 	req.PageSize = 20
	// }
	// if req.PageIndex < 1 {
	// 	req.PageIndex = 1
	// }

	// cacheKey := fmt.Sprintf("xpool:pools:%d", req.ChainId)
	// redisCache := cache.NewRedisCache(fmt.Sprintf("%s:%d", l.svcCtx.Config.Redis.Host, l.svcCtx.Config.Redis.Port), l.svcCtx.Config.Redis.Password, l.svcCtx.Config.Redis.DB)

	// Check cache
	// lastUpdate, err := redisCache.Get(cacheKey)
	// if err == nil && lastUpdate != "" {
	// 	var cacheWarpper CacheWarpper
	// 	err := json.Unmarshal([]byte(lastUpdate), &cacheWarpper)
	// 	if err == nil {
	// 		resp = &types.XpoolPoolListResponse{
	// 			List: cacheWarpper.List,
	// 		}
	// 	}
	// }

	// No cache
	// if resp == nil {
	// 	resp, err = l.getResponse(req)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	// Set cache
	// 	// respBytes, err := json.Marshal(CacheWarpper{List: resp.List, TokenIds: tokenIds, LastUpdate: time.Now().UTC()})
	// 	// if err != nil {
	// 	// 	logx.Errorf("Marshal response error: %v", err)
	// 	// }
	// 	// err = redisCache.Set(cacheKey, respBytes, time.Hour*24)
	// 	// if err != nil {
	// 	// 	logx.Errorf("Set cache error: %v", err)
	// 	// }
	// }

	// resp.Total = int64(len(resp.List))

	// Filter
	if req.UserAddress != "" && common.IsHexAddress(req.UserAddress) {
		conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
		eventFarmTokenStakedModel := repository.NewEventFarmTokenStakedModel(conn)
		validStakeRecords, err := eventFarmTokenStakedModel.QueryValidByChainIdUser(l.ctx, req.ChainId, common.HexToAddress(req.UserAddress).Hex())
		if err != nil {
			return nil, err
		}
		_ = validStakeRecords
		// userTokenIds := []int64{}
		// linq.From(validStakeRecords).ForEachT(func(item repository.EventFarmTokenStaked) {
		// 	userTokenIds = append(userTokenIds, item.EventTokenId)
		// })
		// linq.From(resp.List).WhereT(func(item types.FarmPoolData) bool {
		// 	poolTokenIds, ok := tokenIds[item.PoolAddress]
		// 	if !ok {
		// 		return false
		// 	}
		// 	return linq.From(userTokenIds).AnyWithT(func(tokenId int64) bool {
		// 		return linq.From(poolTokenIds).Contains(tokenId)
		// 	})
		// }).ToSlice(&resp.List)
	}
	if req.OrderBy != "" {
		switch req.OrderBy {
		case "hot":
		case "apr":
		}
	}

	// Page
	// linq.From(resp.List).Skip(int((req.PageIndex - 1) * req.PageSize)).Take(int(req.PageSize)).ToSlice(&resp.List)

	return list, nil
}

func (l *XpoolPoolListLogic) XpoolPoolList(req *types.XpoolPoolListRequest) (resp *response.ApiResponse, err error) {
	chainId := req.ChainId
	// userAddress := req.UserAddress
	// orderBy := req.OrderBy
	// pageIndex := req.PageIndex
	// pageSize := req.PageSize

	chainConf := l.svcCtx.Config.GetChainConfig(chainId)
	if chainConf == nil {
		return nil, fmt.Errorf("%s", "Miss chain configure")
	}
	sonexConf := l.svcCtx.Config.GetSonexConfig(chainId)
	if sonexConf == nil {
		return nil, fmt.Errorf("%s", "Miss sonex configure")
	}
	xpoolConf := l.svcCtx.Config.GetXpoolConfig(chainId)
	if xpoolConf == nil {
		return nil, fmt.Errorf("%s", "Miss xpool configure")
	}

	availablePools, err := l.eventXpoolFactoryPoolCreatedModel.QueryAvailablePoolsByChainIdFactory(l.ctx, chainId, common.HexToAddress(xpoolConf.Factory).Hex())
	if err != nil {
		return nil, err
	}
	var poolAddresses []common.Address
	linq.From(availablePools).SelectT(func(item repository.EventXpoolFactoryPoolCreated) common.Address {
		return common.HexToAddress(item.EventPool)
	}).ToSlice(&poolAddresses)

	staticCaller, err := xpool.NewStaticCaller(l.svcCtx.Config)
	if err != nil {
		return nil, err
	}
	poolInfos, err := staticCaller.BatchGetPoolInfos(l.ctx, req.ChainId, poolAddresses)
	if err != nil {
		return nil, err
	}

	prices, err := l.tokenPriceModel.QueryLatestHistoryPrices(l.ctx)
	if err != nil {
		return nil, err
	}
	list := make([]types.XpoolPoolListResponse_Pool, 0)
	for _, pool := range availablePools {
		info, ok := poolInfos[common.HexToAddress(pool.EventPool).Hex()]
		if !ok {
			continue
		}

		stakingTokenPrice := prices[info.StakingToken.Hex()]
		stakingTokenAccuracy := new(big.Int).Exp(big.NewInt(10), big.NewInt(info.StakingTokenDecimals), nil)
		rewardTokenPrice := prices[info.RewardToken.Hex()]
		rewardTokenAccuracy := new(big.Int).Exp(big.NewInt(10), big.NewInt(info.RewardTokenDecimals), nil)
		tvl := new(big.Int).Div(new(big.Int).Mul(info.TotalStaked, stakingTokenPrice), stakingTokenAccuracy)
		if xpoolConf.BlockTime <= 0 {
			return nil, fmt.Errorf("%s", "Invalid block time configure")
		}
		blockNumberPerYear := new(big.Int).Div(new(big.Int).Mul(big.NewInt(365*24*60*60), big.NewInt(1000)), big.NewInt(xpoolConf.BlockTime))
		totalReward := new(big.Int).Div(new(big.Int).Mul(new(big.Int).Mul(info.RewardPerBlock, blockNumberPerYear), rewardTokenPrice), rewardTokenAccuracy)
		apr := new(big.Int).Div(totalReward, tvl)

		var status string
		if info.IsDestroied {
			status = "closed"
		} else {
			status = "active"
		}

		list = append(list, types.XpoolPoolListResponse_Pool{
			StakingToken:         info.StakingToken.Hex(),
			StakingTokenName:     info.StakingTokenName,
			StakingTokenSymbol:   info.StakingTokenSymbol,
			StakingTokenDecimals: info.StakingTokenDecimals,
			RewardToken:          info.RewardToken.Hex(),
			RewardTokenName:      info.RewardTokenName,
			RewardTokenSymbol:    info.RewardTokenSymbol,
			RewardTokenDecimals:  info.RewardTokenDecimals,
			RewardPerBlock:       info.RewardPerBlock.String(),
			TotalStaked:          info.TotalStaked.String(),
			BlockTime:            xpoolConf.BlockTime,
			TVL:                  tvl.String(),
			APR:                  apr.String(),
			Status:               status,
		})
	}

	return &response.ApiResponse{
		Data: list,
		Meta: types.PageResponse{
			PageSize:  *req.PageSize,
			PageIndex: *req.PageIndex,
			Total:     int64(len(list)),
		},
	}, nil
}
