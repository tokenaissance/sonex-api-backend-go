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
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"
	"strconv"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

const (
	PortfolioTypeLiquidity string = "liquidity"
	PortfolioTypeFarm      string = "farm"
	PortfolioTypeXpool     string = "xpool"
)

// balance|revenue|stakedlp|apr|staked
const (
	OrderByBalance  string = "balance"
	OrderByRevenue  string = "revenue"
	OrderByStakedLP string = "stakedlp"
	OrderByAPR      string = "apr"
	OrderByStaked   string = "staked"
)

// order direction
const (
	OrderDirectionDESC string = "desc"
	OrderDirectionASC  string = "asc"
)

type LiquidityList struct {
	Positions           []*types.Position                               // 流动性持仓列表
	OperationEvents     map[int64][]*sonex.LiquidityOperationEvent      // 流动性持仓操作事件列表
	FarmIncentivesPools []repository.EventFarmIncentiveCreatedWithExtra // 有效并且当前正属于开启状态的 farm incentive pool列表
}

// xpoolsList 获取用户有效 xpool 列表 & 质押事件列表
type xpoolList struct {
	Pools           []types.XpoolPoolListResponse_Pool          // 质押池列表
	OperationEvents map[string][]*xpool.XpoolUserOperationEvent // 质押 & 提取事件列表
}

// FarmIncentiveList 获取用户有效 farm 列表 & 质押事件列表
type FarmIncentiveList struct {
	Incentives      *types.DashboardIncentiveWithStakedUsd         // 质押池列表
	OperationEvents map[int64][]*farm.FarmIncentivesOperationEvent // 质押 & 提取事件列表
}

// Create a struct to hold all results
type portfolioResults struct {
	liquidityList *LiquidityList
	xpoolList     *xpoolList
	incentives    *FarmIncentiveList
	err           error
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

		// 排序参数
		orderBy := req.OrderBy
		orderDirection := req.OrderDirection
		if len(orderBy) <= 0 {
			switch req.PortfolioType {
			case PortfolioTypeLiquidity:
				orderBy = OrderByBalance
			case PortfolioTypeFarm:
				orderBy = OrderByBalance
			case PortfolioTypeXpool:
				orderBy = OrderByStaked
			default:
				orderBy = OrderByBalance
			}
		}
		if len(orderDirection) <= 0 {
			// default desc
			orderDirection = OrderDirectionDESC
		}

		// 用户选中的 portfolio 类型，仅支持单选
		portfolioType := req.PortfolioType

		// 用户所在时区，用于计算 assets rate
		timeZone := req.TimeZone

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

		// 总的已经领取收益
		totalFeeRevenue := 0.0

		// 用户所在时区的当日 00:00
		todayBegin := getUserTimeZoneDayStart(timeZone)
		var previousDayTotalAssets float64
		// 每个 LP position 在当地时间 00:00 时的净值
		lpValueAtDayStart := make(map[int64]float64)
		// 每个 LP position 的当前净值
		lpCurrentValue := make(map[int64]float64)

		// Create channels for each result
		liquidityChan := make(chan portfolioResults, 1)
		xpoolChan := make(chan portfolioResults, 1)
		farmChan := make(chan portfolioResults, 1)

		// Launch goroutine for getLiquidity
		go func() {
			list, err := getLiquidity(r.Context(), svcCtx, req)
			liquidityChan <- portfolioResults{liquidityList: list, err: err}
		}()

		// Launch goroutine for getXPools
		go func() {
			list, err := getXPools(r.Context(), svcCtx, req)
			xpoolChan <- portfolioResults{xpoolList: list, err: err}
		}()

		// Launch goroutine for getFarmPools
		go func() {
			incentives, err := getFarmPools(r.Context(), svcCtx, req)
			farmChan <- portfolioResults{incentives: incentives, err: err}
		}()

		liquidityResult := <-liquidityChan // Get liquidity results
		xpoolResult := <-xpoolChan         // Get xpool results
		farmResult := <-farmChan           // Get farm results
		if liquidityResult.err != nil {
			logx.Errorf("[SonexPortfolioHandler]: SonexPositions error: %s, chainId: %d, address: %s",
				liquidityResult.err.Error(), req.ChainId, req.Address)
			httpx.ErrorCtx(r.Context(), w, liquidityResult.err)
			return
		}

		if xpoolResult.err != nil {
			logx.Errorf("[SonexPortfolioHandler]: getXPools error: %s, chainId: %d, address: %s",
				xpoolResult.err.Error(), req.ChainId, req.Address)
			httpx.ErrorCtx(r.Context(), w, xpoolResult.err)
			return
		}

		if farmResult.err != nil {
			logx.Errorf("[SonexPortfolioHandler]: getFarmPools error: %s, chainId: %d, address: %s",
				farmResult.err.Error(), req.ChainId, req.Address)
			httpx.ErrorCtx(r.Context(), w, farmResult.err)
			return
		}

		// 处理 liquidity 逻辑
		{
			liquidityList = liquidityResult.liquidityList
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

			// count in total position & liquidity claimable fees
			for _, position := range liquidityList.Positions {
				if portfolioType == PortfolioTypeLiquidity {
					// 判断 liquidity position 是否支持 farm stake
					liquidityStakeable := false
					for _, farmIncentive := range liquidityList.FarmIncentivesPools {
						if common.HexToAddress(farmIncentive.ExtraToken0).Cmp(common.HexToAddress(position.Pool.Token0.Address)) == 0 &&
							common.HexToAddress(farmIncentive.ExtraToken1).Cmp(common.HexToAddress(position.Pool.Token1.Address)) == 0 &&
							farmIncentive.ExtraFee == strconv.FormatInt(int64(position.Pool.FeeTier), 10) {
							liquidityStakeable = true
							break
						}
					}
					portfolioPositions = append(portfolioPositions, &types.PortfolioPosition{
						Type:              PortfolioTypeLiquidity,
						LiquidityPosition: position,
						LiquidityStakable: liquidityStakeable,
					})
				}
				portfolioStats.LiquidityPosition += position.Position
				liquidityClaimableFees += position.ClaimableFees
				totalFeeRevenue += position.PNL

				txAt, err := time.Parse(time.RFC3339, position.TxAt)
				if err != nil {
					logx.Errorf("[SonexPortfolioHandler]: time.Parse error: %s, chainId: %d, address: %s", err.Error(), req.ChainId, req.Address)
					continue
				}
				// count in total position assets before today begin
				if txAt.UTC().Before(todayBegin) {
					previousDayTotalAssets += lpValueAtDayStart[position.TokenID]
				}
			}
		}

		// 获取用户有效 xpool 列表
		{
			xpoolList := xpoolResult.xpoolList

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
				earnedRewardTokenValue, err := sonex_module.NewTokenValue(xpool.Earned, int(xpool.RewardTokenDecimals), sonex_module.WithTokenPrice(rewardTokenPrice))
				if err != nil {
					logx.Errorf("[SonexPortfolioHandler]: earnedRewardTokenValue NewTokenValue error: %s, chainId: %d, address: %s", err.Error(), req.ChainId, req.Address)
					continue
				}

				claimedRewardTokenValue, err := sonex_module.NewTokenValue(xpool.Claimed, int(xpool.RewardTokenDecimals), sonex_module.WithTokenPrice(rewardTokenPrice))
				if err != nil {
					logx.Errorf("[SonexPortfolioHandler]: claimedRewardTokenValue NewTokenValue error: %s, chainId: %d, address: %s", err.Error(), req.ChainId, req.Address)
					continue
				}

				xpoolStaked, _ := stakingTokenValue.ValueFloat64()
				claimedRewardTokenValueFloat64, _ := claimedRewardTokenValue.ValueFloat64()
				earnedValueFloat64, _ := earnedRewardTokenValue.ValueFloat64()
				feeRevenueValue := earnedRewardTokenValue.AmountAdd(claimedRewardTokenValue.Amount())
				feeRevenue := feeRevenueValue.ValueString()

				totalFeeRevenue += earnedValueFloat64 + claimedRewardTokenValueFloat64

				if portfolioType == PortfolioTypeXpool {
					portfolioPositions = append(portfolioPositions, &types.PortfolioPosition{
						Type: PortfolioTypeXpool,
						Xpool: &types.PortfolioXpool{
							PoolAddress:         xpool.PoolAddress,               // 质押池合约地址
							StakingTokenName:    xpool.StakingTokenName,          // 质押的代币名称
							StakingTokenSymbol:  xpool.StakingTokenSymbol,        // 质押的代币符号
							StakingTokenAddress: xpool.StakingTokenAddress,       // 质押的代币合约地址
							StakingTokenLogo:    xpool.StakingTokenLogo,          // 质押的代币logo
							RewardTokenLogo:     xpool.RewardTokenLogo,           // 奖励的代币logo
							APR:                 xpool.APR,                       // 年化收益率
							Staked:              stakingTokenValue.ValueString(), // 用户已经质押的美元价值
							FeeRevenue:          feeRevenue,                      // 用户已经领取的收益 + 用户可以领取的收益
							FeeRevenueAmount:    feeRevenueValue.AmountWithoutDecimals(),
							Status:              xpool.Status,
							Earned:              earnedRewardTokenValue.ValueString(), // 用户可以领取的收益
							EarnedAmount:        earnedRewardTokenValue.AmountWithoutDecimals(),
							StakedFloat64:       xpoolStaked,    // 用户已经质押的美元价值
							APRFlot64:           xpool.APRFloat, // 年化收益率 float64，辅助排序用
						},
					})
				}

				portfolioStats.TotalXpoolStaked += xpoolStaked

				// 计算用户可领取的奖励代币价值
				claimableRewardTokenValueFloat64, _ := earnedRewardTokenValue.ValueFloat64()
				xpoolEarns += claimableRewardTokenValueFloat64
			}
		}

		// 获取用户有效 farm 列表
		{
			incentives := farmResult.incentives

			// 2. 然后计算每个 Farm Staking 在当地时间 00:00 时的状态
			// 对于今天 00:00 之前的质押事件，将对应 LP 在当时的净值加入 previousDayTotalAssets
			// 对于今天 00:00 之前的解质押事件，将对应 LP 在当时的净值从 previousDayTotalAssets 中减去
			for _, events := range incentives.OperationEvents {
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

			for _, incentive := range incentives.Incentives.DashboardIncentives {
				feeRevenue := incentive.RevenueFloat64

				totalFeeRevenue += feeRevenue

				stakedUsdFloat64, err := strconv.ParseFloat(incentive.StakedUsd, 64)
				if err != nil {
					logx.Errorf("[SonexPortfolioHandler]: incentive.StakedUsd %s ParseFloat error: %s, chainId: %d, address: %s", incentive.StakedUsd, err.Error(), req.ChainId, req.Address)
					continue
				}

				if portfolioType == PortfolioTypeFarm {
					portfolioPositions = append(portfolioPositions, &types.PortfolioPosition{
						Type: PortfolioTypeFarm,
						FarmPool: &types.PortfolioFarm{
							PoolAddress:          incentive.PoolAddress,
							FeeRatio:             incentive.FeeRatio,
							Token0Symbol:         incentive.Token0Symbol,
							Token0Address:        incentive.Token0Address,
							Token0Logo:           incentive.Token0Logo,
							Token1Symbol:         incentive.Token1Symbol,
							Token1Address:        incentive.Token1Address,
							Token1Logo:           incentive.Token1Logo,
							PositionValue:        incentive.StakedUsd,
							PositionValueFloat64: stakedUsdFloat64,
							FeeRevenue:           fmt.Sprintf("%.2f", feeRevenue),
							FeeRevenueAmount:     incentive.RevenueAmount,
							Earned:               incentive.TotalClaimableFeesUSD,
							EarnedAmount:         incentive.TotalClaimableFeesAmount,
							StakedLPNum:          incentive.Staked, // 质押的LP个数
						},
					})
				}

				portfolioStats.TotalFarmStaked += stakedUsdFloat64

				// 计算用户可领取的奖励代币价值
				farmEarns += incentive.TotalClaimableFeesUSDFloat64
			}
		}

		// 计算总资产和可提资产
		portfolioStats.TotalAssets = portfolioStats.LiquidityPosition + portfolioStats.TotalXpoolStaked + portfolioStats.TotalFarmStaked
		// 尚未被领取的收益
		portfolioStats.ClaimableFees = liquidityClaimableFees + farmEarns + xpoolEarns
		// 已经被领取的收益 + 尚未被领取的收益
		portfolioStats.Earns = totalFeeRevenue
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
		sort(portfolioPositions, orderBy, orderDirection)

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

// sort descending 排序 portfolio positions, orderBy 排序字段, direction 排序方向
func sort(portfolioPositions []*types.PortfolioPosition, orderBy string, direction string) {
	// sort
	if direction == OrderDirectionDESC {
		linq.From(portfolioPositions).OrderByDescendingT(func(item *types.PortfolioPosition) lib_common.CompareableFloat {

			sortValue := float64(0)
			switch item.Type {
			case PortfolioTypeLiquidity:
				if orderBy == OrderByBalance {
					sortValue = item.LiquidityPosition.Position
				}
				if orderBy == OrderByRevenue {
					sortValue = item.LiquidityPosition.PNL
				}
			case PortfolioTypeFarm:
				if orderBy == OrderByBalance {
					sortValue = item.FarmPool.PositionValueFloat64
				}
				if orderBy == OrderByStakedLP {
					sortValue = float64(item.FarmPool.StakedLPNum)
				}
			case PortfolioTypeXpool:
				if orderBy == OrderByStaked {
					sortValue = item.Xpool.StakedFloat64
				}
				if orderBy == OrderByAPR {
					sortValue = item.Xpool.APRFlot64
				}
			}

			return lib_common.CompareableFloat{Value: new(big.Float).SetFloat64(sortValue)}
		}).ToSlice(&portfolioPositions)
	} else {
		linq.From(portfolioPositions).OrderByT(func(item *types.PortfolioPosition) lib_common.CompareableFloat {

			sortValue := float64(0)
			switch item.Type {
			case PortfolioTypeLiquidity:
				if orderBy == OrderByBalance {
					sortValue = item.LiquidityPosition.Position
				}
				if orderBy == OrderByRevenue {
					sortValue = item.LiquidityPosition.PNL
				}
			case PortfolioTypeFarm:
				if orderBy == OrderByBalance {
					sortValue = item.FarmPool.PositionValueFloat64
				}
				if orderBy == OrderByStakedLP {
					sortValue = float64(item.FarmPool.StakedLPNum)
				}
			case PortfolioTypeXpool:
				if orderBy == OrderByStaked {
					sortValue = item.Xpool.StakedFloat64
				}
				if orderBy == OrderByAPR {
					sortValue = item.Xpool.APRFlot64
				}
			}

			return lib_common.CompareableFloat{Value: new(big.Float).SetFloat64(sortValue)}
		}).ToSlice(&portfolioPositions)
	}
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
	address := common.HexToAddress(req.Address).Hex()
	l := sonex.NewSonexPositionsLogic(ctx, svcCtx)
	f := farm.NewFarmIncentivesLogic(ctx, svcCtx)

	type positionsResult struct {
		positions []*types.Position
		err       error
	}
	type operationEventsResult struct {
		operationEvents map[int64][]*sonex.LiquidityOperationEvent
		err             error
	}
	type farmIncentivesPoolsResult struct {
		pools []repository.EventFarmIncentiveCreatedWithExtra
		err   error
	}

	positionsChan := make(chan positionsResult, 1)
	operationEventsChan := make(chan operationEventsResult, 1)
	farmIncentivesPoolsChan := make(chan farmIncentivesPoolsResult, 1)

	go func() {
		positions, err := l.SonexPositions(ctx, req.ChainId, address)
		positionsChan <- positionsResult{positions: positions, err: err}
	}()

	go func() {
		operationEvents, err := l.GetLiquidityOperationEventsWithTokenValue(ctx, req.ChainId, address)
		operationEventsChan <- operationEventsResult{operationEvents: operationEvents, err: err}
	}()

	go func() {
		pools, err := f.GetFarmIncentives(req.ChainId)
		farmIncentivesPoolsChan <- farmIncentivesPoolsResult{pools: pools, err: err}
	}()

	positions := <-positionsChan
	operationEvents := <-operationEventsChan
	farmIncentivesPools := <-farmIncentivesPoolsChan

	if positions.err != nil {
		logx.Errorf("[SonexPortfolioHandler]:getPositions SonexPositions error: %s, chainId: %d, address: %s", positions.err.Error(), req.ChainId, req.Address)
		return nil, positions.err
	}
	if operationEvents.err != nil {
		logx.Errorf("[SonexPortfolioHandler]:getPositions GetLiquidityOperationEvents error: %s, chainId: %d, address: %s", operationEvents.err.Error(), req.ChainId, req.Address)
		return nil, operationEvents.err
	}
	if farmIncentivesPools.err != nil {
		logx.Errorf("[SonexPortfolioHandler]:getFarmPools GetFarmIncentives error: %s, chainId: %d, address: %s", farmIncentivesPools.err.Error(), req.ChainId, req.Address)
		return nil, farmIncentivesPools.err
	}

	return &LiquidityList{
		Positions:           positions.positions,
		OperationEvents:     operationEvents.operationEvents,
		FarmIncentivesPools: farmIncentivesPools.pools,
	}, nil
}

// getXPools 获取用户有效 xpool 列表
func getXPools(ctx context.Context, svcCtx *svc.ServiceContext, req types.PortfolioRequest) (*xpoolList, error) {
	address := common.HexToAddress(req.Address).Hex()
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
			UserAddress: address,
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
		operationEvents, err := l.GetUserOperationEvents(req.ChainId, address)
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
func getFarmPools(ctx context.Context, svcCtx *svc.ServiceContext, req types.PortfolioRequest) (*FarmIncentiveList, error) {
	address := common.HexToAddress(req.Address).Hex()

	l := farm.NewFarmIncentivesLogic(ctx, svcCtx)
	f := farm.NewFarmUnstakePositionLogic(ctx, svcCtx)

	type farmIncentivesResult struct {
		incentives *types.DashboardIncentiveWithStakedUsd
		err        error
	}

	type farmIncentivesOperationEventsResult struct {
		operationEvents map[int64][]*farm.FarmIncentivesOperationEvent
		err             error
	}

	incentivesChan := make(chan farmIncentivesResult, 1)
	operationEventsChan := make(chan farmIncentivesOperationEventsResult, 1)

	go func() {
		incentives, err := f.GetFarmIncentives(svcCtx.Config, req.ChainId, address)
		incentivesChan <- farmIncentivesResult{incentives: incentives, err: err}
	}()

	go func() {
		// 用户质押&解质押事件
		operationEvents, err := l.GetFarmIncentivesOperationEvents(req.ChainId, address)
		operationEventsChan <- farmIncentivesOperationEventsResult{operationEvents: operationEvents, err: err}
	}()

	incentives := <-incentivesChan
	operationEvents := <-operationEventsChan

	if incentives.err != nil {
		logx.Errorf("[SonexPortfolioHandler]:getFarmPools FarmIncentives error: %s, chainId: %d, address: %s", incentives.err.Error(), req.ChainId, req.Address)
		return nil, incentives.err
	}
	if operationEvents.err != nil {
		logx.Errorf("[SonexPortfolioHandler]:getFarmPools GetFarmIncentivesOperationEvents error: %s, chainId: %d, address: %s", operationEvents.err.Error(), req.ChainId, req.Address)
		return nil, operationEvents.err
	}

	return &FarmIncentiveList{
		Incentives:      incentives.incentives,
		OperationEvents: operationEvents.operationEvents,
	}, nil
}
