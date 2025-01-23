package xpool

import (
	"context"
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/handler/response"
	lib_common "sonex_api_backend_go/internal/lib/common"
	"sonex_api_backend_go/internal/module/sonex"
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

func (l *XpoolPoolListLogic) XpoolPoolList(req *types.XpoolPoolListRequest) (resp *response.ApiResponse, err error) {
	chainId := req.ChainId
	userAddress := req.UserAddress
	orderBy := req.OrderBy
	pageIndex := req.PageIndex
	pageSize := req.PageSize

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

	// get available xpools
	availablePools, err := l.eventXpoolFactoryPoolCreatedModel.QueryAvailablePoolsByChainIdFactory(l.ctx, chainId, common.HexToAddress(xpoolConf.Factory).Hex())
	if err != nil {
		logx.Errorf("[XpoolPoolList]: Get available pools error: %s", err)
		return nil, err
	}
	var poolAddresses []common.Address
	linq.From(availablePools).SelectT(func(item repository.EventXpoolFactoryPoolCreated) common.Address {
		return common.HexToAddress(item.EventPool)
	}).ToSlice(&poolAddresses)

	staticCaller, err := xpool.NewStaticCaller(l.svcCtx.Config)
	if err != nil {
		logx.Errorf("[XpoolPoolList]: New static caller error: %s", err)
		return nil, err
	}
	poolInfos, err := staticCaller.BatchGetPoolInfos(l.ctx, chainId, poolAddresses)
	if err != nil {
		logx.Errorf("[XpoolPoolList]:BatchGetPoolInfos Batch get pool infos error: %s", err)
		return nil, err
	}
	stakingTokenAddresses := make([]string, 0)
	linq.From(poolInfos).SelectT(func(kv linq.KeyValue) string {
		poolInfo, ok := kv.Value.(*xpool.PoolInfo)
		if !ok {
			logx.Errorf("[XpoolPoolList]: Invalid pool info: %s", kv.Value)
			return ""
		}
		return poolInfo.StakingToken.Hex()
	}).ToSlice(&stakingTokenAddresses)

	rewardTokenAddresses := make([]string, 0)
	linq.From(poolInfos).SelectT(func(kv linq.KeyValue) string {
		poolInfo, ok := kv.Value.(*xpool.PoolInfo)
		if !ok {
			logx.Errorf("[XpoolPoolList]: Invalid pool info: %s", kv.Value)
			return ""
		}
		return poolInfo.RewardToken.Hex()
	}).ToSlice(&rewardTokenAddresses)

	// user wallet address pool stats
	addressPoolInfos := make(map[string]*xpool.AddressPoolInfo, 0)
	if len(userAddress) > 0 {
		addressPoolInfos, err = staticCaller.BatchGetAddressPoolInfos(l.ctx, chainId, poolAddresses, common.HexToAddress(userAddress))
		if err != nil {
			logx.Errorf("[XpoolPoolList]: BatchGetAddressPoolInfos error: %s", err)
			return nil, err
		}
	}

	// 获取当前价格
	prices, err := l.getCurrentPrices(l.ctx, append(stakingTokenAddresses, rewardTokenAddresses...))
	if err != nil {
		logx.Errorf("[XpoolPoolList]: Get current prices error: %s", err)
		return nil, err
	}
	list := make([]types.XpoolPoolListResponse_Pool, 0)
	for _, pool := range availablePools {
		info, ok := poolInfos[common.HexToAddress(pool.EventPool).Hex()]
		if !ok {
			logx.Errorf("[XpoolPoolList]: Pool info not found: %s, pools: %v", common.HexToAddress(pool.EventPool).Hex(), poolInfos)
			continue
		}

		tvlValue := new(big.Float)
		rewardPerYearValue := new(big.Float)
		var stakingTokenPriceString string
		var rewardTokenPriceString string
		staked := new(big.Int)
		earned := new(big.Int)

		// 获取用户在池子中的质押状态信息
		{
			addressPoolInfo, ok := addressPoolInfos[common.HexToAddress(pool.EventPool).Hex()]
			if ok {
				staked = addressPoolInfo.Staked
				earned = addressPoolInfo.EarnedReward
			}
		}

		// staking token
		{
			stakingTokenPrice, exists := prices[info.StakingToken.Hex()]
			if !exists {
				logx.Errorf("[XpoolPoolList]: Staking token price not found: %s", info.StakingToken.Hex())
				continue
			}
			stakingTokenPriceString = stakingTokenPrice.PriceWithoutDecimals()
			stakingTokenValue, err := sonex.NewTokenValue(info.TotalStaked.String(), int(info.StakingTokenDecimals), sonex.WithTokenPrice(stakingTokenPrice))
			if err != nil {
				logx.Errorf("[XpoolPoolList]: Invalid staking token value: %s", err)
				continue
			}
			tvlValue = stakingTokenValue.Value()

			logx.Infof("[XpoolPoolList]: Staking token value: %f", tvlValue)
		}

		// reward token
		{
			rewardTokenPrice, exists := prices[info.RewardToken.Hex()]
			if !exists {
				logx.Errorf("[XpoolPoolList]: Reward token price not found: %s", info.RewardToken.Hex())
				continue
			}
			rewardTokenPriceString = rewardTokenPrice.PriceWithoutDecimals()
			if xpoolConf.BlockTime <= 0 {
				logx.Errorf("[XpoolPoolList]: Invalid block time configure")
				continue
			}
			blockNumberPerYear := new(big.Int).Div(new(big.Int).Mul(big.NewInt(365*24*60*60), big.NewInt(1000)), big.NewInt(xpoolConf.BlockTime))
			rewardAmountPerYear := new(big.Int).Mul(info.RewardPerBlock, blockNumberPerYear)
			rewardTokenValue, err := sonex.NewTokenValue(rewardAmountPerYear.String(), int(info.RewardTokenDecimals), sonex.WithTokenPrice(rewardTokenPrice))
			if err != nil {
				logx.Errorf("[XpoolPoolList]: Invalid reward token value: %s", err)
				continue
			}
			rewardPerYearValue = rewardTokenValue.Value()
			logx.Infof("[XpoolPoolList]: Reward token value: %f", rewardPerYearValue)
		}

		// APR
		apr := new(big.Float).Quo(rewardPerYearValue, tvlValue)
		if tvlValue.Cmp(big.NewFloat(0)) == 0 {
			apr = big.NewFloat(0)
		}
		aprFloat, _ := apr.Float64()
		aprPercentage, _ := new(big.Float).Mul(apr, big.NewFloat(100)).Float64()
		aprPercentageString := fmt.Sprintf("%.2f%%", aprPercentage)

		// status
		var status string
		if info.IsDestroied {
			status = "closed"
		} else {
			status = "active"
		}

		list = append(list, types.XpoolPoolListResponse_Pool{
			StakingTokenName:     info.StakingTokenName,                     // 质押的代币名称
			StakingTokenSymbol:   info.StakingTokenSymbol,                   // 质押的代币符号
			StakingTokenAddress:  info.StakingToken.Hex(),                   // 质押的代币合约地址
			StakingTokenDecimals: info.StakingTokenDecimals,                 // 质押的代币合约精度
			StakingTokenPrice:    stakingTokenPriceString,                   // 质押的代币价格
			RewardTokenName:      info.RewardTokenName,                      // 奖励的代币名称
			RewardTokenSymbol:    info.RewardTokenSymbol,                    // 奖励的代币符号
			RewardTokenAddress:   info.RewardToken.Hex(),                    // 奖励的代币合约地址
			RewardTokenDecimals:  info.RewardTokenDecimals,                  // 奖励的代币合约精度
			RewardTokenPrice:     rewardTokenPriceString,                    // 奖励的代币价格
			PoolAddress:          common.HexToAddress(pool.EventPool).Hex(), // 质押池合约地址
			RewardPerBlock:       info.RewardPerBlock.String(),              // 出块奖励
			TotalStaked:          info.TotalStaked.String(),                 // 质押池的质押币总数量
			BlockTime:            xpoolConf.BlockTime,                       // 出块时长 单位ms
			TVL:                  tvlValue.String(),                         // 质押池的质押总价值（USD）
			APR:                  aprPercentageString,                       // 年化收益率 取 12.22% 格式
			APRFloat:             aprFloat,                                  // 年化收益率，辅助排序用
			JoinNum:              info.StakedNumber.Int64(),                 // 参与人数
			Staked:               staked.String(),                           // 用户已质押代币数量
			Earned:               earned.String(),                           // 用户已获得奖励代币数量
			Status:               status,                                    // 运行状态
			CreatedAt:            pool.CreatedAt.Format(time.RFC3339),       // 创建时间
			UpdatedAt:            pool.UpdatedAt.Format(time.RFC3339),       // 更新时间
		})
	}

	// Sorts
	switch orderBy {
	case "hot":
		linq.From(list).OrderByDescendingT(func(item types.XpoolPoolListResponse_Pool) int64 {
			return item.JoinNum
		}).ToSlice(&list)
	case "apr":
		linq.From(list).OrderByDescendingT(func(item types.XpoolPoolListResponse_Pool) lib_common.CompareableFloat {
			return lib_common.CompareableFloat{Value: new(big.Float).SetFloat64(item.APRFloat)}
		}).ToSlice(&list)
	case "tvl":
		// do not sort
	}

	return &response.ApiResponse{
		Data: list,
		Meta: types.PageResponse{
			PageSize:  *pageSize,
			PageIndex: *pageIndex,
			Total:     int64(len(list)),
		},
	}, nil
}

// getCurrentPrices 获取当前价格
func (l *XpoolPoolListLogic) getCurrentPrices(ctx context.Context, addressList []string) (prices map[string]*sonex.TokenPrice, err error) {
	records, err := l.tokenPriceModel.GetCurrentPrices(ctx, addressList)
	if err != nil {
		logx.Errorf("[XpoolPoolList]: Get current prices error: %s", err)
		return nil, err
	}
	prices = make(map[string]*sonex.TokenPrice, 0)
	for _, record := range records {
		tp, err := sonex.NewTokenPrice(record.Price, record.UsdDecimals)
		if err != nil {
			logx.Errorf("[XpoolPoolList]: Invalid token price: %s", err)
			return nil, err
		}
		prices[record.TokenAddress] = tp
	}

	// todo db 获取不到的话，从链上获取价格
	return prices, nil
}
