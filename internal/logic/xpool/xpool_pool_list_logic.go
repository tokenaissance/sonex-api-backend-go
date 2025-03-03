package xpool

import (
	"context"
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/handler/response"
	lib_common "sonex_api_backend_go/internal/lib/common"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/module/sonex"
	"sonex_api_backend_go/internal/module/xpool"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"
	"strings"
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
	eventXpoolPoolDepositedModel      repository.EventXpoolPoolDepositedModel
	eventXpoolPoolWithdrawedModel     repository.EventXpoolPoolWithdrawedModel
	tokenPriceModel                   repository.SonexTokenPriceRecordModel
}

type CacheWarpper struct {
	List       []types.XpoolPoolListResponse_Pool
	LastUpdate time.Time
}

// Xpool list
func NewXpoolPoolListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *XpoolPoolListLogic {
	conn := repository.NewConn(svcCtx.Config.PostgreSQL.DataSource())
	return &XpoolPoolListLogic{
		Logger:                            logx.WithContext(ctx),
		ctx:                               ctx,
		svcCtx:                            svcCtx,
		eventXpoolFactoryPoolCreatedModel: repository.NewEventXpoolFactoryPoolCreatedModel(conn),
		eventXpoolPoolDepositedModel:      repository.NewEventXpoolPoolDepositedModel(conn),
		eventXpoolPoolWithdrawedModel:     repository.NewEventXpoolPoolWithdrawedModel(conn),
		tokenPriceModel:                   repository.NewSonexTokenPriceRecordModel(conn),
	}
}

// XpoolUserEventType 用户操作事件类型
type XpoolUserEventType int64

const (
	XpoolUserEventTypeDeposited XpoolUserEventType = iota
	XpoolUserEventTypeWithdrawed
)

type XpoolDepositedEvent struct {
	Id          int64     // 主键
	CreatedAt   time.Time // 创建时间
	UpdatedAt   time.Time // 更新时间
	ChainId     int64     // 链ID
	BlockNumber int64     // 区块高度
	TxIndex     int64     // 交易索引
	TxHash      string    // 交易哈希
	TxAt        time.Time // 交易时间
	From        string    // 交易发送者
	To          string    // 交易接收者
	Value       string    // 事件交易发生的 ETH amount，如果交易不包含 ETH，则该值为 0
	EventUser   string    // 事件发送者
	EventAmount string    // deposited amount token 数量
	EventHeight int64     // deposited event 发生的区块高度
	EventValue  float64   // deposited amount token 的美元价值
}

type XpoolWithdrawedEvent struct {
	Id             int64     // 主键
	CreatedAt      time.Time // 创建时间
	UpdatedAt      time.Time // 更新时间
	ChainId        int64     // 链ID
	BlockNumber    int64     // 区块高度
	TxIndex        int64     // 交易索引
	TxHash         string    // 交易哈希
	TxAt           time.Time // 交易时间
	From           string    // 交易发送者
	To             string    // 交易接收者
	Value          string    // 事件交易发生的 ETH amount，如果交易不包含 ETH，则该值为 0
	EventUser      string    // 事件发送者
	EventWithdrawn string    // 事件提款数量
	EventClaimed   string    // 事件已领取数量
	EventHeight    int64     // 事件发生时的区块高度
	EventValue     float64   // 事件提款金额的美元价值
}

// XpoolUserOperationEvent 用户操作事件，根据 type 类型，返回不同的 event 类型
type XpoolUserOperationEvent struct {
	XpoolUserEventType XpoolUserEventType
	DepositedEvent     *XpoolDepositedEvent
	WithdrawedEvent    *XpoolWithdrawedEvent
}

// GetUserOperationEvents 获取用户存款事件
func (l *XpoolPoolListLogic) GetUserOperationEvents(chainId int64, userAddress string) (m map[string][]*XpoolUserOperationEvent, err error) {
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

	// get user deposited events
	depositedEvents, err := l.eventXpoolPoolDepositedModel.GetDepositedEvents(l.ctx, chainId, userAddress)
	if err != nil {
		logx.Errorf("[XpoolPoolList]: Get user deposited events error: %s", err)
		return nil, err
	}

	withdrawedEvents, err := l.eventXpoolPoolWithdrawedModel.GetWithdrawedEvents(l.ctx, chainId, userAddress)
	if err != nil {
		logx.Errorf("[XpoolPoolList]: Get user withdrawed events error: %s", err)
		return nil, err
	}

	// Use a map to track unique addresses
	seenAddresses := make(map[string]bool)
	poolAddresses := make([]common.Address, 0)

	// Add addresses from deposited events
	for address := range depositedEvents {
		if !seenAddresses[address] {
			seenAddresses[address] = true
			poolAddresses = append(poolAddresses, common.HexToAddress(address))
		}
	}

	// Add addresses from withdrawn events
	for address := range withdrawedEvents {
		if !seenAddresses[address] {
			seenAddresses[address] = true
			poolAddresses = append(poolAddresses, common.HexToAddress(address))
		}
	}

	// get pool infos
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

	// get staking token addresses
	stakingTokenAddresses := make([]string, 0)
	linq.From(poolInfos).SelectT(func(kv linq.KeyValue) string {
		poolInfo, ok := kv.Value.(*xpool.PoolInfo)
		if !ok {
			logx.Errorf("[XpoolPoolList]: Invalid pool info: %s", kv.Value)
			return ""
		}
		return poolInfo.StakingToken.Hex()
	}).ToSlice(&stakingTokenAddresses)

	// get reward token addresses
	rewardTokenAddresses := make([]string, 0)
	linq.From(poolInfos).SelectT(func(kv linq.KeyValue) string {
		poolInfo, ok := kv.Value.(*xpool.PoolInfo)
		if !ok {
			logx.Errorf("[XpoolPoolList]: Invalid pool info: %s", kv.Value)
			return ""
		}
		return poolInfo.RewardToken.Hex()
	}).ToSlice(&rewardTokenAddresses)

	// get current prices
	prices, err := l.getCurrentPrices(l.ctx, chainId, append(stakingTokenAddresses, rewardTokenAddresses...))
	if err != nil {
		logx.Errorf("[XpoolPoolList]: Get current prices error: %s", err)
		return nil, err
	}

	m = make(map[string][]*XpoolUserOperationEvent, 0)
	// enumerate pool address and deposited events
	for address, events := range depositedEvents {
		info, ok := poolInfos[common.HexToAddress(address).Hex()]
		if !ok {
			logx.Errorf("[XpoolPoolList]: Pool info not found: %s, pools: %v", common.HexToAddress(address).Hex(), poolInfos)
			continue
		}

		// get staking token price
		stakingTokenPrice, exists := prices[info.StakingToken.Hex()]
		if !exists {
			logx.Errorf("[XpoolPoolList]: Staking token price not found: %s", info.StakingToken.Hex())
			continue
		}

		operationEvents, ok := m[address]
		if !ok {
			operationEvents = make([]*XpoolUserOperationEvent, 0)
		}
		// count deposited(stakeed) amount
		for _, event := range events {
			stakingTokenValue, err := sonex.NewTokenValue(event.EventAmount, int(info.StakingTokenDecimals), sonex.WithTokenPrice(stakingTokenPrice))
			if err != nil {
				logx.Errorf("[XpoolPoolList]: Invalid staking token value: %s", err)
				continue
			}

			eventValue, _ := stakingTokenValue.ValueFloat64()
			depositedEvent := &XpoolDepositedEvent{
				Id:          event.Id,
				CreatedAt:   event.CreatedAt,
				UpdatedAt:   event.UpdatedAt,
				ChainId:     event.ChainId,
				BlockNumber: event.BlockNumber,
				TxIndex:     event.TxIndex,
				TxHash:      event.TxHash,
				TxAt:        event.TxAt,
				From:        event.From,
				To:          event.To,
				Value:       event.Value,
				EventUser:   event.EventUser,
				EventAmount: event.EventAmount,
				EventHeight: event.EventHeight,
				EventValue:  eventValue,
			}
			operationEvent := &XpoolUserOperationEvent{
				XpoolUserEventType: XpoolUserEventTypeDeposited,
				DepositedEvent:     depositedEvent,
			}
			operationEvents = append(operationEvents, operationEvent)
		}
		// deposited events map
		m[address] = operationEvents
	}

	// enumerate pool address and deposited events
	for address, events := range withdrawedEvents {
		info, ok := poolInfos[common.HexToAddress(address).Hex()]
		if !ok {
			logx.Errorf("[XpoolPoolList]: Pool info not found: %s, pools: %v", common.HexToAddress(address).Hex(), poolInfos)
			continue
		}

		// get staking token price
		stakingTokenPrice, exists := prices[info.StakingToken.Hex()]
		if !exists {
			logx.Errorf("[XpoolPoolList]: Staking token price not found: %s", info.StakingToken.Hex())
			continue
		}

		operationEvents, ok := m[address]
		if !ok {
			operationEvents = make([]*XpoolUserOperationEvent, 0)
		}
		// count deposited(stakeed) amount
		for _, event := range events {
			stakingTokenValue, err := sonex.NewTokenValue(event.EventWithdrawn, int(info.StakingTokenDecimals), sonex.WithTokenPrice(stakingTokenPrice))
			if err != nil {
				logx.Errorf("[XpoolPoolList]: Invalid staking token value: %s", err)
				continue
			}

			eventValue, _ := stakingTokenValue.ValueFloat64()
			withdrawedEvent := &XpoolWithdrawedEvent{
				Id:             event.Id,
				CreatedAt:      event.CreatedAt,
				UpdatedAt:      event.UpdatedAt,
				ChainId:        event.ChainId,
				BlockNumber:    event.BlockNumber,
				TxIndex:        event.TxIndex,
				TxHash:         event.TxHash,
				TxAt:           event.TxAt,
				From:           event.From,
				To:             event.To,
				Value:          event.Value,
				EventUser:      event.EventUser,
				EventWithdrawn: event.EventWithdrawn,
				EventClaimed:   event.EventClaimed,
				EventHeight:    event.EventHeight,
				EventValue:     eventValue,
			}
			operationEvent := &XpoolUserOperationEvent{
				XpoolUserEventType: XpoolUserEventTypeWithdrawed,
				WithdrawedEvent:    withdrawedEvent,
			}
			operationEvents = append(operationEvents, operationEvent)
		}
		// deposited events map
		m[address] = operationEvents
	}

	return m, nil
}

func (l *XpoolPoolListLogic) XpoolPoolList(req *types.XpoolPoolListRequest) (resp *response.ApiResponse, err error) {
	chainId := req.ChainId
	userAddress := req.UserAddress
	orderBy := req.OrderBy
	searchKeyword := req.Keyword
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
	prices, err := l.getCurrentPrices(l.ctx, chainId, append(stakingTokenAddresses, rewardTokenAddresses...))
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
		var stakingTokenPrice *sonex.TokenPrice
		var stakingTokenPriceString string
		var rewardTokenPrice *sonex.TokenPrice
		var rewardTokenPriceString string
		staked := new(big.Int)
		earned := new(big.Int)
		claimed := new(big.Int)

		// 获取用户在池子中的质押状态信息
		{
			addressPoolInfo, ok := addressPoolInfos[common.HexToAddress(pool.EventPool).Hex()]
			if ok {
				staked = addressPoolInfo.Staked
				earned = addressPoolInfo.EarnedReward
				claimed = addressPoolInfo.Claimed
			}
		}

		// staking token
		{
			price, exists := prices[info.StakingToken.Hex()]
			if !exists {
				logx.Errorf("[XpoolPoolList]: Staking token price not found: %s", info.StakingToken.Hex())
				continue
			}
			stakingTokenPrice = price
			stakingTokenPriceString = stakingTokenPrice.PriceWithoutDecimals()
			stakingTokenValue, err := sonex.NewTokenValue(info.TotalStaked.String(), int(info.StakingTokenDecimals), sonex.WithTokenPrice(stakingTokenPrice))
			if err != nil {
				logx.Errorf("[XpoolPoolList]: Invalid staking token value: %s", err)
				continue
			}
			tvlValue = stakingTokenValue.Value()
		}

		// reward token
		{
			price, exists := prices[info.RewardToken.Hex()]
			if !exists {
				logx.Errorf("[XpoolPoolList]: Reward token price not found: %s", info.RewardToken.Hex())
				continue
			}
			rewardTokenPrice = price
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
		}

		// APR
		if tvlValue.Cmp(big.NewFloat(1)) <= 0 {
			tvlValue = big.NewFloat(1)
		}
		apr := new(big.Float).Quo(rewardPerYearValue, tvlValue)

		aprFloat, _ := apr.Float64()
		aprPercentage, _ := new(big.Float).Mul(apr, big.NewFloat(100)).Float64()
		aprPercentageString := fmt.Sprintf("%.2f%%", aprPercentage)

		// status
		var status string
		{
			if info.IsDestroied {
				status = "closed"
			} else {
				status = "active"
			}
		}

		// filter staked only, if user address is not empty, only return staked xpool.
		{
			if (len(userAddress) > 0) && (staked.Cmp(big.NewInt(0)) <= 0) {
				continue
			}
		}

		// filter keyword，如果有搜索关键字，进行不区分大小写的模糊匹配过滤
		{
			if len(searchKeyword) > 0 && !strings.Contains(
				strings.ToLower(fmt.Sprintf("%s %s", info.StakingTokenName, info.StakingTokenSymbol)),
				strings.ToLower(searchKeyword)) {
				continue
			}
		}

		list = append(list, types.XpoolPoolListResponse_Pool{
			StakingTokenName:     info.StakingTokenName,                     // 质押的代币名称
			StakingTokenSymbol:   info.StakingTokenSymbol,                   // 质押的代币符号
			StakingTokenAddress:  info.StakingToken.Hex(),                   // 质押的代币合约地址
			StakingTokenDecimals: info.StakingTokenDecimals,                 // 质押的代币精度
			StakingTokenPrice:    stakingTokenPriceString,                   // 质押的代币价格
			RewardTokenName:      info.RewardTokenName,                      // 奖励的代币名称
			RewardTokenSymbol:    info.RewardTokenSymbol,                    // 奖励的代币符号
			RewardTokenAddress:   info.RewardToken.Hex(),                    // 奖励的代币合约地址
			RewardTokenDecimals:  info.RewardTokenDecimals,                  // 奖励的代币精度
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

			StakingTokenPriceWithDecimals: stakingTokenPrice.Price(),           // 质押的代币价格(包含精度)
			StakingTokenPriceDecimals:     int64(stakingTokenPrice.Decimals()), // 质押的代币价格精度
			RewardTokenPriceWithDecimals:  rewardTokenPrice.Price(),            // 奖励的代币价格(包含精度)
			RewardTokenPriceDecimals:      int64(rewardTokenPrice.Decimals()),  // 奖励的代币价格精度
			Claimed:                       claimed.String(),                    // 用户已领取的奖励代币数量
			CreatedAt:                     pool.CreatedAt.Format(time.RFC3339), // 创建时间
			UpdatedAt:                     pool.UpdatedAt.Format(time.RFC3339), // 更新时间
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

	// Pagination
	total := int64(len(list))
	list = l.getPagination(list, pageIndex, pageSize)

	return &response.ApiResponse{
		Data: list,
		Meta: types.PageResponse{
			PageSize:  *pageSize,
			PageIndex: *pageIndex,
			Total:     total, // 使用未分页前的总数
		},
	}, nil
}

func (l *XpoolPoolListLogic) getPagination(list []types.XpoolPoolListResponse_Pool, pageIndex *int64, pageSize *int64) []types.XpoolPoolListResponse_Pool {
	start := (*pageIndex - 1) * (*pageSize)
	if start < 0 {
		start = 0
	}
	end := start + *pageSize
	if start >= int64(len(list)) {
		list = []types.XpoolPoolListResponse_Pool{}
	} else {
		if end > int64(len(list)) {
			end = int64(len(list))
		}
		list = list[start:end]
	}
	return list
}

// getCurrentPrices 获取当前价格
func (l *XpoolPoolListLogic) getCurrentPrices(ctx context.Context, chainId int64, addressList []string) (prices map[string]*sonex.TokenPrice, err error) {
	sonexConf := l.svcCtx.Config.GetSonexConfig(chainId)
	if sonexConf == nil {
		return nil, fmt.Errorf("%s", "Miss sonex configure")
	}
	xTokenAddr := common.HexToAddress(sonexConf.XToken)

	chainConf := l.svcCtx.Config.GetChainConfig(chainId)
	if chainConf == nil {
		return nil, fmt.Errorf("Miss chain configure")
	}
	usdToken, err := global.GetToken(chainId, common.HexToAddress(chainConf.USD).Hex())
	if err != nil {
		return nil, err
	}

	// unique address list
	filterM := make(map[string]bool)
	uniqueAddr := make([]string, 0)
	for _, addr := range addressList {
		if _, ok := filterM[addr]; ok {
			continue
		}
		filterM[addr] = true
		uniqueAddr = append(uniqueAddr, addr)
	}

	addressList = uniqueAddr

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

	// set xToken price
	for _, addr := range addressList {
		// 1. db 价格存在，pass
		if _, ok := prices[addr]; ok {
			continue
		}

		// 设置默认价格为 1U
		tp, err := sonex.NewTokenPrice("1", 0)
		if err != nil {
			logx.Errorf("[XpoolPoolList]: Invalid token price: %s", err)
			return nil, err
		}

		// 2. db 价格不存在，当 addr 为 xToken 时，设置价格为默认 1U
		if xTokenAddr.Cmp(common.HexToAddress(addr)) == 0 {
			prices[addr] = tp
			continue
		}

		// 3. db 价格不存在，且不是 xToken，则从链上查询数据，如果获取不到，则将价格设置为 1U
		bestPrice0, _, _, err := global.GetPrices(ctx, nil, chainId, common.HexToAddress(addr).Hex(), true, 2)
		if err != nil {
			logx.Errorf("[XpoolPoolList]: getCurrentPrices GetPrices error: %s", err.Error())
		}
		if err == nil && bestPrice0 != nil {
			tp.SetPrice(bestPrice0.Value.String())
			tp.SetDecimals(int(usdToken.Decimals))
		}
		prices[addr] = tp
	}

	return prices, nil
}
