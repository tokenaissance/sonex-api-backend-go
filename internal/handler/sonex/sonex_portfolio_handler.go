package sonex

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"net/http"
	lib_common "sonex_api_backend_go/internal/lib/common"
	"sonex_api_backend_go/internal/logic/farm"
	"sonex_api_backend_go/internal/logic/sonex"
	"sonex_api_backend_go/internal/logic/xpool"
	sonex_module "sonex_api_backend_go/internal/module/sonex"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

const (
	portfolioTypeLiquidity string = "liquidity"
	portfolioTypeFarm      string = "farm"
	portfolioTypeXpool     string = "xpool"
)

type LiquidityList struct {
	Positions       []*types.Position                          // 流动性持仓列表
	OperationEvents map[int64][]*sonex.LiquidityOperationEvent // 流动性持仓操作事件列表
	// FarmIncentivesPools []repository.EventFarmIncentiveCreatedWithExtra // 有效并且当前正属于开启状态的 farm incentive pool列表
}

// xpoolsList 获取用户有效 xpool 列表 & 质押事件列表
type xpoolList struct {
	Pools           []types.XpoolPoolListResponse_Pool          // 质押池列表
	OperationEvents map[string][]*xpool.XpoolUserOperationEvent // 质押 & 提取事件列表
}

// FarmIncentiveList 获取用户有效 farm 列表 & 质押事件列表
type FarmIncentiveList struct {
	Pools           []types.FarmIncentivesResponse_Data            // 质押池列表
	OperationEvents map[int64][]*farm.FarmIncentivesOperationEvent // 质押 & 提取事件列表
}

func SonexPortfolioHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PortfolioRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 分页参数
		pageIndex := *req.PageIndex
		pageSize := *req.PageSize

		// 用户选中的 portfolio 类型，支持单选/复选
		liquiditySelected := true
		farmSelected := req.FarmSelected
		xpoolSelected := req.XpoolSelected

		// 用户所在时区，用于计算 assets rate
		timeZone := "Asia/Shanghai"

		// 计算用户选中的 portfolio 数量
		selectedCnt := 0
		if liquiditySelected {
			selectedCnt++
		}
		if farmSelected {
			selectedCnt++
		}
		if xpoolSelected {
			selectedCnt++
		}

		response := &types.PortfolioResponse{}

		// 持仓列表
		portfolioPositions := make([]*types.PortfolioPosition, 0)
		// 流动性持仓列表
		var liquidityList *LiquidityList
		// 持仓统计
		portfolioStats := &types.PortfolioStats{}

		// 可提收益
		liquidityClaimableFees := 0.0
		farmEarns := 0.0
		xpoolEarns := 0.0

		// 用户所在时区的当日 00:00
		todayBegin := getUserTimeZoneDayStart(timeZone)
		var previousDayTotalAssets float64
		// 每个 LP position 在当地时间 00:00 时的净值
		lpValueAtDayStart := make(map[int64]float64)
		// 每个 LP position 的当前净值
		lpCurrentValue := make(map[int64]float64)

		// 获取用户有效 positions 列表
		{
			if liquiditySelected || farmSelected {
				list, err := getLiquidity(r.Context(), svcCtx, req)
				if err != nil {
					logx.Errorf("[SonexPortfolioHandler]: SonexPositions error: %s, chainId: %d, address: %s", err.Error(), req.ChainId, req.Address)
					httpx.ErrorCtx(r.Context(), w, err)
					return
				}
				liquidityList = list

				// 1. 首先计算每个 LP position 在当地时间 00:00 时的净值
				// 计算每个 LP position 的当前净值

				for tokenId, events := range liquidityList.OperationEvents {
					var valueBeforeDayStart float64
					var currentValue float64

					for _, event := range events {
						if event.LiquidityOperationEventType == sonex.LiquidityOperationEventTypeIncrease {
							value, _ := event.IncreaseEvent.EventValue.Float64()
							if event.IncreaseEvent.TxAt.Before(todayBegin) {
								// 累计今天之前的净值变化
								valueBeforeDayStart += value
							}
							// 累计所有净值变化得到当前值
							currentValue += value
						} else if event.LiquidityOperationEventType == sonex.LiquidityOperationEventTypeDecrease {
							value, _ := event.DecreaseEvent.EventValue.Float64()
							if event.DecreaseEvent.TxAt.Before(todayBegin) {
								// 累计今天之前的净值变化
								valueBeforeDayStart -= value
							}
							// 累计所有净值变化得到当前值
							currentValue -= value
						}
					}

					lpValueAtDayStart[tokenId] = valueBeforeDayStart
					lpCurrentValue[tokenId] = currentValue
				}
			}

			if liquiditySelected {
				liquidityFrames := make([]*types.PortfolioLiquidityFrame, 0)
				// count in total position & liquidity claimable fees
				for _, position := range liquidityList.Positions {

					// 判断 liquidity position 是否支持 farm stake
					liquidityStakeable := false
					// for _, farmIncentive := range liquidityList.FarmIncentivesPools {
					// 	if common.HexToAddress(farmIncentive.ExtraToken0) == common.HexToAddress(position.Pool.Token0.Address) && common.HexToAddress(farmIncentive.ExtraToken1) == common.HexToAddress(position.Pool.Token0.Address) {
					// 		liquidityStakeable = true
					// 		break
					// 	}
					// }
					portfolioPositions = append(portfolioPositions, &types.PortfolioPosition{
						Type:              portfolioTypeLiquidity,
						LiquidityPosition: position,
						LiquidityStakable: liquidityStakeable,
					})
					portfolioStats.LiquidityPosition += position.Position
					liquidityClaimableFees += position.ClaimableFees

					txAt, err := time.Parse(time.RFC3339, position.TxAt)
					if err != nil {
						logx.Errorf("[SonexPortfolioHandler]: time.Parse error: %s, chainId: %d, address: %s", err.Error(), req.ChainId, req.Address)
						continue
					}
					// count in total position assets before today begin
					if txAt.UTC().Before(todayBegin) {
						previousDayTotalAssets += lpValueAtDayStart[position.TokenID]
					}

					// 计算 liquidity position 的饼图 frame 数据
					pool := position.Pool
					liquidityFrames = append(liquidityFrames, &types.PortfolioLiquidityFrame{
						Position:      position.Position,
						Token0Name:    pool.Token0.Name,
						Token0Symbol:  pool.Token0.Symbol,
						Token0Address: pool.Token0.Address,
						Token1Name:    pool.Token1.Name,
						Token1Symbol:  pool.Token1.Symbol,
						Token1Address: pool.Token1.Address,
					})
				}
				// 单选时，才会返回 liquidityFrames
				if selectedCnt == 1 {
					portfolioStats.LiquidityFrames = liquidityFrames
				}
			}
		}

		// 获取用户有效 xpool 列表
		{
			if xpoolSelected {
				xpoolList, err := getXPools(r.Context(), svcCtx, req)
				if err != nil {
					logx.Errorf("[SonexPortfolioHandler]: getXPools error: %s, chainId: %d, address: %s", err.Error(), req.ChainId, req.Address)
					httpx.ErrorCtx(r.Context(), w, err)
					return
				}

				// count in total xpool staked before today begin
				for _, events := range xpoolList.OperationEvents {
					for _, event := range events {
						// 质押事件
						if event.XpoolUserEventType == xpool.XpoolUserEventTypeDeposited {
							// deposited event transaction time
							txAt := event.DepositedEvent.TxAt
							// count in total xpool staked before today begin
							if txAt.Before(todayBegin) {
								previousDayTotalAssets += event.DepositedEvent.EventValue
							}
						}
						// 提取事件
						if event.XpoolUserEventType == xpool.XpoolUserEventTypeWithdrawed {
							// withdrawn event transaction time
							txAt := event.WithdrawedEvent.TxAt
							// count in total xpool staked before today begin
							if txAt.Before(todayBegin) {
								previousDayTotalAssets -= event.WithdrawedEvent.EventValue
							}
						}
					}
				}

				xpoolFrames := make([]*types.PortfolioXpoolFrame, 0)
				for _, xpool := range xpoolList.Pools {
					stakingTokenPrice, err := sonex_module.NewTokenPrice(xpool.StakingTokenPriceWithDecimals, int(xpool.StakingTokenPriceDecimals))
					if err != nil {
						logx.Errorf("[SonexPortfolioHandler]: stakingTokenPrice NewTokenPrice error: %s, chainId: %d, address: %s", err.Error(), req.ChainId, req.Address)
						continue
					}
					stakingTokenValue, err := sonex_module.NewTokenValue(xpool.Staked, int(xpool.StakingTokenDecimals), sonex_module.WithTokenPrice(stakingTokenPrice))
					if err != nil {
						logx.Errorf("[SonexPortfolioHandler]: NewTokenValue error: %s, chainId: %d, address: %s", err.Error(), req.ChainId, req.Address)
						continue
					}
					rewardTokenPrice, err := sonex_module.NewTokenPrice(xpool.RewardTokenPriceWithDecimals, int(xpool.RewardTokenPriceDecimals))
					if err != nil {
						logx.Errorf("[SonexPortfolioHandler]: rewardTokenPrice NewTokenPrice error: %s, chainId: %d, address: %s", err.Error(), req.ChainId, req.Address)
						continue
					}
					rewardTokenValue, err := sonex_module.NewTokenValue(xpool.Earned, int(xpool.RewardTokenDecimals), sonex_module.WithTokenPrice(rewardTokenPrice))
					if err != nil {
						logx.Errorf("[SonexPortfolioHandler]: rewardTokenValue NewTokenValue error: %s, chainId: %d, address: %s", err.Error(), req.ChainId, req.Address)
						continue
					}

					claimedRewardTokenValue, err := sonex_module.NewTokenValue(xpool.Claimed, int(xpool.RewardTokenDecimals), sonex_module.WithTokenPrice(rewardTokenPrice))
					if err != nil {
						logx.Errorf("[SonexPortfolioHandler]: claimedRewardTokenValue NewTokenValue error: %s, chainId: %d, address: %s", err.Error(), req.ChainId, req.Address)
						continue
					}

					xpoolStaked, _ := stakingTokenValue.ValueFloat64()
					claimedRewardTokenValueFloat64, _ := claimedRewardTokenValue.ValueFloat64()
					earnedValueFloat64, _ := rewardTokenValue.ValueFloat64()
					feeRevenue := fmt.Sprintf("%.2f", earnedValueFloat64+claimedRewardTokenValueFloat64)

					portfolioPositions = append(portfolioPositions, &types.PortfolioPosition{
						Type: portfolioTypeXpool,
						Xpool: &types.PortfolioXpool{
							PoolAddress:         xpool.PoolAddress,               // 质押池合约地址
							StakingTokenName:    xpool.StakingTokenName,          // 质押的代币名称
							StakingTokenSymbol:  xpool.StakingTokenSymbol,        // 质押的代币符号
							StakingTokenAddress: xpool.StakingTokenAddress,       // 质押的代币合约地址
							APR:                 xpool.APR,                       // 年化收益率
							Staked:              stakingTokenValue.ValueString(), // 用户已经质押的美元价值
							StakedFloat64:       xpoolStaked,                     // 用户已经质押的美元价值
							FeeRevenue:          feeRevenue,                      // 用户已经领取的收益 + 用户可以领取的收益
							Status:              xpool.Status,
							Earned:              rewardTokenValue.ValueString(), // 用户可以领取的收益
						},
					})

					portfolioStats.TotalXpoolStaked += xpoolStaked

					// 计算用户可领取的奖励代币价值
					claimableRewardTokenValueFloat64, _ := rewardTokenValue.ValueFloat64()
					xpoolEarns += claimableRewardTokenValueFloat64

					xpoolFrames = append(xpoolFrames, &types.PortfolioXpoolFrame{
						Staked:              xpoolStaked,
						StakingTokenName:    xpool.StakingTokenName,
						StakingTokenSymbol:  xpool.StakingTokenSymbol,
						StakingTokenAddress: xpool.StakingTokenAddress,
					})
				}

				// 单选时，才会返回 xpoolFrames
				if selectedCnt == 1 {
					portfolioStats.XpoolFrames = xpoolFrames
				}
			}
		}

		// 获取用户有效 farm 列表
		{
			if farmSelected {
				farmPools, err := getFarmPools(r.Context(), svcCtx, req, liquidityList.Positions)
				if err != nil {
					logx.Errorf("[SonexPortfolioHandler]: getFarmPools error: %s, chainId: %d, address: %s", err.Error(), req.ChainId, req.Address)
					httpx.ErrorCtx(r.Context(), w, err)
					return
				}

				// 2. 然后计算每个 Farm Staking 在当地时间 00:00 时的状态
				// 对于今天 00:00 之前的质押事件，将对应 LP 在当时的净值加入 previousDayTotalAssets
				// 对于今天 00:00 之前的解质押事件，将对应 LP 在当时的净值从 previousDayTotalAssets 中减去
				for _, events := range farmPools.OperationEvents {
					for _, event := range events {
						positionId := int64(0)
						if event.FarmIncentivesOperationEventType == farm.FarmIncentivesOperationEventTypeStaked {
							positionId = event.StakedEvent.EventTokenId
							if event.StakedEvent.TxAt.Before(todayBegin) {
								// 如果是今天之前的质押，计入前一天总资产
								if lpValue, exists := lpValueAtDayStart[positionId]; exists {
									previousDayTotalAssets += lpValue
								}
							}
						} else if event.FarmIncentivesOperationEventType == farm.FarmIncentivesOperationEventTypeUnstaked {
							positionId = event.UnstakedEvent.EventTokenId
							if event.UnstakedEvent.TxAt.Before(todayBegin) {
								// 如果是今天之前的解质押，从前一天总资产中减去
								if lpValue, exists := lpValueAtDayStart[positionId]; exists {
									previousDayTotalAssets -= lpValue
								}
							}
						}
					}
				}

				farmFrames := make([]*types.PortfolioFarmFrame, 0)
				for _, farm := range farmPools.Pools {
					feeRevenue := farm.TotalClaimedFeesUSDFloat64 + farm.TotalClaimableFeesUSDFloat64

					portfolioPositions = append(portfolioPositions, &types.PortfolioPosition{
						Type: portfolioTypeFarm,
						FarmPool: &types.PortfolioFarm{
							PoolAddress:          farm.PoolAddress,
							FeeRatio:             farm.FeeRatio,
							Token0Symbol:         farm.Token0Symbol,
							Token0Address:        farm.Token0Address,
							Token1Symbol:         farm.Token1Symbol,
							Token1Address:        farm.Token1Address,
							PositionValue:        farm.LiquidityUSD,
							PositionValueFloat64: farm.LiquidityUSDFloat64,
							FeeRevenue:           fmt.Sprintf("%.2f", feeRevenue),
							Earned:               farm.TotalClaimableFeesUSD,
						},
					})

					portfolioStats.TotalFarmStaked += farm.LiquidityUSDFloat64

					// 计算用户可领取的奖励代币价值
					farmEarns += farm.TotalClaimableFeesUSDFloat64

					farmFrames = append(farmFrames, &types.PortfolioFarmFrame{
						Position:      farm.LiquidityUSDFloat64,
						Token0Symbol:  farm.Token0Symbol,
						Token0Address: farm.Token0Address,
						Token1Symbol:  farm.Token1Symbol,
						Token1Address: farm.Token1Address,
					})
				}

				// 单选时，才会返回 farmFrames
				if selectedCnt == 1 {
					portfolioStats.FarmFrames = farmFrames
				}
			}
		}

		// 计算总资产和可提资产
		portfolioStats.TotalAssets = portfolioStats.LiquidityPosition + portfolioStats.TotalXpoolStaked + portfolioStats.TotalFarmStaked
		portfolioStats.ClaimableFees = liquidityClaimableFees + farmEarns + xpoolEarns

		// Calculate total assets rate, handling division by zero
		if previousDayTotalAssets > 0 {
			portfolioStats.TotalAssetsRate = fmt.Sprintf("%.2f%%", (portfolioStats.TotalAssets/previousDayTotalAssets-1.0)*100)
		} else if (previousDayTotalAssets == 0) && (portfolioStats.TotalAssets > 0) {
			portfolioStats.TotalAssetsRate = "100.00%"
		} else {
			portfolioStats.TotalAssetsRate = "0.00%"
		}
		// Calculate total assets change
		portfolioStats.TotalAssetsChange = math.Abs(portfolioStats.TotalAssets - previousDayTotalAssets)
		response.Stats = portfolioStats

		// sort portfolio positions
		sort(portfolioPositions)

		total := len(portfolioPositions)
		// pagination
		list := getPagination(portfolioPositions, pageIndex, pageSize)
		response.Positions = list
		response.Meta = types.PageResponse{
			PageSize:  pageSize,
			PageIndex: pageIndex,
			Total:     int64(total),
		}

		// hotfix to compatible with online version
		response.TotalPosition = portfolioStats.TotalAssets
		response.ClaimableFees = portfolioStats.ClaimableFees

		httpx.OkJsonCtx(r.Context(), w, response)
	}
}

// sort descending 排序 portfolio positions
func sort(portfolioPositions []*types.PortfolioPosition) {
	// sort
	linq.From(portfolioPositions).OrderByDescendingT(func(item *types.PortfolioPosition) lib_common.CompareableFloat {

		positionValue := float64(0)
		switch item.Type {
		case portfolioTypeLiquidity:
			positionValue = item.LiquidityPosition.Position
		case portfolioTypeFarm:
			positionValue = item.FarmPool.PositionValueFloat64
		case portfolioTypeXpool:
			positionValue = item.Xpool.StakedFloat64
		}

		return lib_common.CompareableFloat{Value: new(big.Float).SetFloat64(positionValue)}
	}).ToSlice(&portfolioPositions)
}

// getUserTimeZoneDayStart 获取用户所在时区的当日开始时间
func getUserTimeZoneDayStart(timeZone string) time.Time {
	// Parse timeZone
	loc, err := time.LoadLocation(timeZone)
	// if timeZone is invalid, fallback to UTC + 0
	if err != nil {
		logx.Errorf("[SonexPortfolioHandler]:getUserTimeZoneDayStart time.LoadLocation error: %s, timeZone: %s, fallback to UTC + 0", err.Error(), timeZone)
		loc = time.UTC
	}

	// get current time in user's time zone
	locNow := time.Now().In(loc)
	// get the start of the day in user's time zone
	localDayStart := time.Date(locNow.Year(), locNow.Month(), locNow.Day(), 0, 0, 0, 0, loc)

	return localDayStart
}

// getPagination 获取分页数据
func getPagination(list []*types.PortfolioPosition, pageIndex int64, pageSize int64) []*types.PortfolioPosition {
	start := (pageIndex - 1) * pageSize
	if start < 0 {
		start = 0
	}
	end := start + pageSize
	if start >= int64(len(list)) {
		list = []*types.PortfolioPosition{}
	} else {
		if end > int64(len(list)) {
			end = int64(len(list))
		}
		list = list[start:end]
	}
	return list
}

// getLiquidity 获取用户有效 positions 列表
func getLiquidity(ctx context.Context, svcCtx *svc.ServiceContext, req types.PortfolioRequest) (*LiquidityList, error) {
	l := sonex.NewSonexPositionsLogic(ctx, svcCtx)
	// f := farm.NewFarmIncentivesLogic(ctx, svcCtx)

	type positionsResult struct {
		positions []*types.Position
		err       error
	}
	type operationEventsResult struct {
		operationEvents map[int64][]*sonex.LiquidityOperationEvent
		err             error
	}
	// type farmIncentivesPoolsResult struct {
	// 	pools []repository.EventFarmIncentiveCreatedWithExtra
	// 	err   error
	// }

	positionsChan := make(chan positionsResult, 1)
	operationEventsChan := make(chan operationEventsResult, 1)
	// farmIncentivesPoolsChan := make(chan farmIncentivesPoolsResult, 1)

	go func() {
		positions, err := l.SonexPositions(ctx, req.ChainId, req.Address)
		positionsChan <- positionsResult{positions: positions, err: err}
	}()

	go func() {
		operationEvents, err := l.GetLiquidityOperationEventsWithTokenValue(ctx, req.ChainId, req.Address)
		operationEventsChan <- operationEventsResult{operationEvents: operationEvents, err: err}
	}()

	// go func() {
	// 	pools, err := f.GetFarmIncentives(req.ChainId)
	// 	farmIncentivesPoolsChan <- farmIncentivesPoolsResult{pools: pools, err: err}
	// }()

	positions := <-positionsChan
	operationEvents := <-operationEventsChan
	// farmIncentivesPools := <-farmIncentivesPoolsChan

	if positions.err != nil {
		logx.Errorf("[SonexPortfolioHandler]:getPositions SonexPositions error: %s, chainId: %d, address: %s", positions.err.Error(), req.ChainId, req.Address)
		return nil, positions.err
	}
	if operationEvents.err != nil {
		logx.Errorf("[SonexPortfolioHandler]:getPositions GetLiquidityOperationEvents error: %s, chainId: %d, address: %s", operationEvents.err.Error(), req.ChainId, req.Address)
		return nil, operationEvents.err
	}
	// if farmIncentivesPools.err != nil {
	// 	logx.Errorf("[SonexPortfolioHandler]:getFarmPools GetFarmIncentives error: %s, chainId: %d, address: %s", farmIncentivesPools.err.Error(), req.ChainId, req.Address)
	// 	return nil, farmIncentivesPools.err
	// }

	return &LiquidityList{
		Positions:       positions.positions,
		OperationEvents: operationEvents.operationEvents,
		// FarmIncentivesPools: farmIncentivesPools.pools,
	}, nil
}

// getXPools 获取用户有效 xpool 列表
func getXPools(ctx context.Context, svcCtx *svc.ServiceContext, req types.PortfolioRequest) (*xpoolList, error) {
	l := xpool.NewXpoolPoolListLogic(ctx, svcCtx)

	type poolsResult struct {
		pools []types.XpoolPoolListResponse_Pool
		err   error
	}
	type operationEventsResult struct {
		operationEvents map[string][]*xpool.XpoolUserOperationEvent
		err             error
	}

	poolsChan := make(chan poolsResult, 1)
	operationEventsChan := make(chan operationEventsResult, 1)

	go func() {
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
			poolsChan <- poolsResult{err: fmt.Errorf("[SonexPortfolioHandler]:getXPools XpoolPoolList error: %s, chainId: %d, address: %s", err.Error(), req.ChainId, req.Address)}
			return
		}

		list, ok := pools.Data.([]types.XpoolPoolListResponse_Pool)
		if !ok {
			poolsChan <- poolsResult{err: fmt.Errorf("[SonexPortfolioHandler]:getXPools XpoolPoolList data type error, chainId: %d, address: %s", req.ChainId, req.Address)}
			return
		}
		poolsChan <- poolsResult{pools: list}
	}()

	go func() {
		operationEvents, err := l.GetUserOperationEvents(req.ChainId, req.Address)
		if err != nil {
			operationEventsChan <- operationEventsResult{err: fmt.Errorf("[SonexPortfolioHandler]:getXPools GetUserOperationEvents error: %s, chainId: %d, address: %s", err.Error(), req.ChainId, req.Address)}
			return
		}
		operationEventsChan <- operationEventsResult{operationEvents: operationEvents}
	}()

	pools := <-poolsChan
	operationEvents := <-operationEventsChan

	if pools.err != nil {
		return nil, pools.err
	}
	if operationEvents.err != nil {
		return nil, operationEvents.err
	}

	return &xpoolList{
		Pools:           pools.pools,
		OperationEvents: operationEvents.operationEvents,
	}, nil
}

// getFarmPools 获取用户有效 farm 列表
func getFarmPools(ctx context.Context, svcCtx *svc.ServiceContext, req types.PortfolioRequest, positions []*types.Position) (*FarmIncentiveList, error) {
	l := farm.NewFarmIncentivesLogic(ctx, svcCtx)

	type farmIncentivesResult struct {
		pools []types.FarmIncentivesResponse_Data
		err   error
	}

	type farmIncentivesOperationEventsResult struct {
		operationEvents map[int64][]*farm.FarmIncentivesOperationEvent
		err             error
	}

	poolsChan := make(chan farmIncentivesResult, 1)
	operationEventsChan := make(chan farmIncentivesOperationEventsResult, 1)

	go func() {
		pools, err := l.FarmIncentives(&types.FarmIncentivesRequest{
			ChainRequest: types.ChainRequest{
				ChainId: req.ChainId,
			},
			Status:      "active",
			UserAddress: req.Address,
			StakedOnly:  true,
		}, positions)
		poolsChan <- farmIncentivesResult{pools: pools.List, err: err}
	}()

	go func() {
		// 用户质押&解质押事件
		operationEvents, err := l.GetFarmIncentivesOperationEvents(req.ChainId, req.Address)
		operationEventsChan <- farmIncentivesOperationEventsResult{operationEvents: operationEvents, err: err}
	}()

	pools := <-poolsChan
	operationEvents := <-operationEventsChan

	if pools.err != nil {
		logx.Errorf("[SonexPortfolioHandler]:getFarmPools FarmIncentives error: %s, chainId: %d, address: %s", pools.err.Error(), req.ChainId, req.Address)
		return nil, pools.err
	}
	if operationEvents.err != nil {
		logx.Errorf("[SonexPortfolioHandler]:getFarmPools GetFarmIncentivesOperationEvents error: %s, chainId: %d, address: %s", operationEvents.err.Error(), req.ChainId, req.Address)
		return nil, operationEvents.err
	}

	return &FarmIncentiveList{
		Pools:           pools.pools,
		OperationEvents: operationEvents.operationEvents,
	}, nil
}
