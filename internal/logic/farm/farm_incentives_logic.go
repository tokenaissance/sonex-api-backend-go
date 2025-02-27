package farm

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/cache"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/module/sonex"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type FarmIncentivesLogic struct {
	logx.Logger
	ctx                         context.Context
	svcCtx                      *svc.ServiceContext
	eventFarmTokenStakedModel   repository.EventFarmTokenStakedModel
	eventFarmTokenUnstakedModel repository.EventFarmTokenUnstakedModel
	eventIncentiveCreatedModel  repository.EventFarmIncentiveCreatedModel
}

type CacheWarpper struct {
	List       []types.FarmIncentivesResponse_Data
	TokenIds   map[int64][]int64
	LastUpdate time.Time
}

// Farm incentive list
func NewFarmIncentivesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FarmIncentivesLogic {
	conn := repository.NewConn(svcCtx.Config.PostgreSQL.DataSource())
	return &FarmIncentivesLogic{
		Logger:                      logx.WithContext(ctx),
		ctx:                         ctx,
		svcCtx:                      svcCtx,
		eventFarmTokenStakedModel:   repository.NewEventFarmTokenStakedModel(conn),
		eventFarmTokenUnstakedModel: repository.NewEventFarmTokenUnstakedModel(conn),
		eventIncentiveCreatedModel:  repository.NewEventFarmIncentiveCreatedModel(conn),
	}
}

func (l *FarmIncentivesLogic) FarmIncentives(req *types.FarmIncentivesRequest, positions ...[]*types.Position) (resp *types.FarmIncentivesResponse, err error) {
	if req.ChainId <= 0 {
		return resp, fmt.Errorf("invalid parameter")
	}

	farmConf := l.svcCtx.Config.GetFarmConfig(req.ChainId)
	chain := l.svcCtx.Config.GetChainConfig(req.ChainId)
	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DataSource())
	sonexConf := l.svcCtx.Config.GetSonexConfig(req.ChainId)

	redisCache := cache.NewRedisCache(fmt.Sprintf("%s:%d", l.svcCtx.Config.Redis.Host, l.svcCtx.Config.Redis.Port), l.svcCtx.Config.Redis.Password, l.svcCtx.Config.Redis.DB)

	resp, tokenIds, err := l.GetResponse(req, chain, conn, farmConf, sonexConf, redisCache)

	if err != nil {
		logx.Errorf("ERR from GetResponse:%v", err)
		return resp, err
	}

	if len(tokenIds) == 0 {
		logx.Info("no tokenid founds")
		return resp, nil
	}

	// Filter
	if !strings.EqualFold(req.UserAddress, "") && common.IsHexAddress(req.UserAddress) {

		if len(resp.List) > 0 {
			for i, item := range resp.List {
				bal, _, _, _, _ := getBalance(req.ChainId, item.IncentiveId, req.UserAddress, farmConf.Staker, conn)
				resp.List[i].Balance = bal
			}
		}

		usdAccuracy, _ := new(big.Float).SetString(resp.USD.Price)

		_ = usdAccuracy
		eventFarmTokenStakedModel := repository.NewEventFarmTokenStakedModel(conn)

		validStakeRecords, err := eventFarmTokenStakedModel.QueryValidByChainIdUserWithExtra(l.ctx, req.ChainId, common.HexToAddress(req.UserAddress).Hex())
		if err != nil {
			logx.Errorf("ERR:%v", err)
			return resp, err
		}

		userTokenIds := map[int64][]int64{}
		linq.From(validStakeRecords).GroupByT(
			func(item repository.EventFarmTokenStakedWithExtra) int64 { return item.EventIncentiveId },
			func(item repository.EventFarmTokenStakedWithExtra) repository.EventFarmTokenStakedWithExtra {
				return item
			},
		).ForEachT(func(item linq.Group) {
			tids := []int64{}
			linq.From(item.Group).ForEachT(func(recordA interface{}) {
				record := recordA.(repository.EventFarmTokenStakedWithExtra)
				tids = append(tids, record.EventTokenId)
			})
			userTokenIds[item.Key.(int64)] = tids
		})
		staked := map[int64]int64{}
		for _, data := range resp.List {
			poolTokenIds, ok := tokenIds[data.IncentiveId]
			if !ok {
				continue
			}
			for _, idItem := range userTokenIds[data.IncentiveId] {
				if linq.From(poolTokenIds).Contains(idItem) {
					staked[data.IncentiveId]++
				}
			}
		}

		var incentiveIds []int64
		linq.From(resp.List).SelectT(func(item types.FarmIncentivesResponse_Data) int64 { return item.IncentiveId }).ToSlice(&incentiveIds)
		for i, item := range resp.List {
			if val, ok := staked[item.IncentiveId]; ok {
				resp.List[i].Staked = val
			} else {
				resp.List[i].Staked = val
			}
			totalClaimableFeesString, totalClaimableFeesUSD, _, totalClaimableFeesUSDFloat64, err := getRewardsDetailInfo(farmConf.Staker, req.UserAddress, int64(item.IncentiveId), req.ChainId, chain, conn)

			if err != nil {
				return resp, err
			}
			resp.List[i].TotalClaimableFeesUSD = totalClaimableFeesUSD.String()
			resp.List[i].TotalClaimableFeesUSDFloat64 = totalClaimableFeesUSDFloat64

			logx.Infof("UserAddress:%v incentiveid:%v totalClaimableFeesString:%v ", req.UserAddress, item.IncentiveId, totalClaimableFeesString)

			totalClaimedAmount, err := getClainedAmount(req.UserAddress, req.ChainId, item.IncentiveId, conn)
			if err != nil {
				return resp, err
			}

			totalClaimedFeesUSD_bigFloat, _, err := getValueOfAmountToken(totalClaimedAmount, item.IncentiveId, req.ChainId, chain, conn)

			if err != nil {
				return resp, err
			}
			resp.List[i].TotalClaimedFeesUSD = totalClaimedFeesUSD_bigFloat.String()
			totalClaimedFeesUSDFloat64, _ := totalClaimedFeesUSD_bigFloat.Float64()
			resp.List[i].TotalClaimedFeesUSDFloat64 = totalClaimedFeesUSDFloat64
		}

	}
	if req.StakedOnly {
		linq.From(resp.List).WhereT(func(item types.FarmIncentivesResponse_Data) bool {
			return item.Staked > 0
		}).ToSlice(&resp.List)
	}
	if !strings.EqualFold(req.KeyWord, "") {
		linq.From(resp.List).WhereT(func(item types.FarmIncentivesResponse_Data) bool {
			return strings.Contains(item.Token0Symbol, req.KeyWord) || strings.Contains(item.Token1Symbol, req.KeyWord)

		}).ToSlice(&resp.List)
	}
	if !strings.EqualFold(req.Status, "") && (strings.EqualFold(req.Status, "active") || strings.EqualFold(req.Status, "closed")) {
		linq.From(resp.List).WhereT(func(item types.FarmIncentivesResponse_Data) bool {
			return item.Status == req.Status
		}).ToSlice(&resp.List)
	}

	// Sorts
	if !strings.EqualFold(req.OrderBy, "") {
		switch req.OrderBy {
		case "hot":
			sort.SliceStable(resp.List, func(i, j int) bool {
				if resp.List[i].JoinNum == resp.List[j].JoinNum {
					return resp.List[i].IncentiveId > resp.List[j].IncentiveId
				} else {
					return resp.List[i].JoinNum > resp.List[j].JoinNum
				}
			})

		case "apr":

			sort.SliceStable(resp.List, func(i, j int) bool {
				apr_i, _ := big.NewFloat(0).SetString(resp.List[i].APR)
				apr_j, _ := big.NewFloat(0).SetString(resp.List[j].APR)
				if apr_i.Cmp(apr_j) == 0 {
					return resp.List[i].IncentiveId > resp.List[j].IncentiveId
				} else {
					return apr_i.Cmp(apr_j) > 0
				}
			})

		default:
			sort.SliceStable(resp.List, func(i, j int) bool {
				if resp.List[i].JoinNum == resp.List[j].JoinNum {
					return resp.List[i].IncentiveId > resp.List[j].IncentiveId
				} else {
					return resp.List[i].JoinNum > resp.List[j].JoinNum
				}
			})
		}
	}

	// Total
	resp.Total = int64(len(resp.List))
	return resp, nil
}

func (l *FarmIncentivesLogic) GetResponse(req *types.FarmIncentivesRequest, chainConf *config.ChainConfig, conn sqlx.SqlConn, farmConf *config.FarmConfig, sonexConf *config.SonexConfig, cacher cache.ICache) (resp *types.FarmIncentivesResponse, tokenIds map[int64][]int64, err error) {
	usdToken, err := global.GetToken(req.ChainId, chainConf.USD)
	if err != nil {
		logx.Errorf("ERR:%v ", err)
		return resp, nil, err
	}
	usdAccuracy := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(usdToken.Decimals), nil))

	eventFarmIncentiveCreatedModel := repository.NewEventFarmIncentiveCreatedModel(conn)
	eventFarmTokenStakedModel := repository.NewEventFarmTokenStakedModel(conn)
	tokenPriceModel := repository.NewSonexTokenPriceRecordModel(conn)

	incentiveCreatedRecords, err := eventFarmIncentiveCreatedModel.QueryByChainIdStakerWithExtra(l.ctx, req.ChainId, common.HexToAddress(farmConf.Staker).Hex())
	if err != nil {
		logx.Errorf("ERR1:%v ", err)
		return resp, nil, err
	}
	validStakeRecords, err := eventFarmTokenStakedModel.QueryValidByChainId(l.ctx, req.ChainId)
	if err != nil {
		logx.Errorf("ERR2:%v ", err)
		return resp, nil, err
	}
	validStakeRecordMap := map[int64][]repository.EventFarmTokenStaked{}
	linq.From(validStakeRecords).
		GroupByT(
			func(item repository.EventFarmTokenStaked) int64 { return item.EventIncentiveId },
			func(item repository.EventFarmTokenStaked) repository.EventFarmTokenStaked { return item },
		).
		ToMapByT(
			&validStakeRecordMap,
			func(item linq.Group) int64 { return item.Key.(int64) },
			func(item linq.Group) []repository.EventFarmTokenStaked {
				value := []repository.EventFarmTokenStaked{}
				for _, v := range item.Group {
					value = append(value, v.(repository.EventFarmTokenStaked))
				}
				return value
			},
		)

	prices, err := tokenPriceModel.QueryLatestHistoryPrices(l.ctx)
	if err != nil {
		logx.Errorf("ERR3:%v ", err)
		return resp, nil, err
	}

	var wg sync.WaitGroup
	results := make(chan struct {
		Data     types.FarmIncentivesResponse_Data
		TokenIds []int64
		Err      error
	}, len(incentiveCreatedRecords))

	for _, incincentiveCreatedRecordItem := range incentiveCreatedRecords {
		wg.Add(1)
		go func(record *repository.EventFarmIncentiveCreatedWithExtra, ch chan struct {
			Data     types.FarmIncentivesResponse_Data
			TokenIds []int64
			Err      error
		}) {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					ch <- struct {
						Data     types.FarmIncentivesResponse_Data
						TokenIds []int64
						Err      error
					}{
						Err: fmt.Errorf("%d %s(panic)", record.EventIncentiveId, err),
					}
				}
			}()

			token0, err := global.GetToken(req.ChainId, record.ExtraToken0)
			if err != nil {
				logx.Errorf("%s", fmt.Sprintf("Get token0 error: %s %s", record.ExtraToken0, err))
				return
			}
			token1, err := global.GetToken(req.ChainId, record.ExtraToken1)
			if err != nil {
				logx.Errorf("%s", fmt.Sprintf("Get token1 error: %s %s", record.ExtraToken1, err))
				return
			}
			price0, ok := prices[token0.Address]
			if !ok {
				logx.Errorf("Get price error:%v ", token0.Address)
				return
			}
			price1, ok := prices[token1.Address]
			if !ok {
				logx.Errorf("Get price error:%v ", token1.Address)
				return
			}

			_price0divUsdDecimals := new(big.Float).Quo(new(big.Float).SetInt(price0), new(big.Float).SetFloat64(math.Pow10(int(6))))
			_price1divUsdDecimals := new(big.Float).Quo(new(big.Float).SetInt(price1), new(big.Float).SetFloat64(math.Pow10(int(6))))

			// price0BigFloat :=  big.NewFloat()

			rewardToken, err := global.GetToken(req.ChainId, record.EventRewardToken)
			if err != nil {
				logx.Errorf("%s", fmt.Sprintf("Get extra token error: %s %s", record.EventRewardToken, err))
				return
			}
			ratepersecond, _ := new(big.Float).SetString(record.ExtraRate)
			val, ok := prices[rewardToken.Address]
			if !ok {
				val = big.NewInt(1000000) //  1U
			}
			rewardPrice := new(big.Float).SetInt(val)
			// rewardTokenAccuracy := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(rewardToken.Decimals), nil))
			rewardTokenAccuracy := new(big.Float).SetFloat64(math.Pow10(int(rewardToken.Decimals)))

			rewardUSD := new(big.Float).Quo(new(big.Float).Quo(new(big.Float).Mul(ratepersecond, rewardPrice), rewardTokenAccuracy), usdAccuracy)

			totalRewardUSD := new(big.Float).Mul(rewardUSD, big.NewFloat(365*24*60*60))

			err = global.SetTotalReward365ForIncentiveId(cacher, global.Reward365ForIncentiveIdKey, record.EventIncentiveId, totalRewardUSD.String())
			if err != nil {
				logx.Errorf("SetTotalReward365ForIncentiveId ERR:%v ", err)
				return
			}
			var dataResult types.FarmIncentivesResponse_Data
			tokenIdsResult := []int64{}
			var status string
			if record.ExtraEnded {
				status = "closed"
			} else {
				status = "active"
			}

			rewards, err := getRewardsAmount(farmConf.Staker, req.UserAddress, record.EventIncentiveId, *chainConf)
			if err != nil {
				logx.Errorf("getRewardsAmount ERR:%v ", err)
				return
			}

			totalClaimableFeesUSD, _, err := getValueOfAmountToken(rewards, record.EventIncentiveId, req.ChainId, chainConf, conn)

			if err != nil {
				logx.Errorf("getValueOfAmountToken:ERR %v ", err)
				return
			}

			totalClaimableFeesUSDFloat64, _ := totalClaimableFeesUSD.Float64()

			eventWeightString := big.NewInt(record.EventWeight).String()
			totalWeightString, err := SetRewardPoolInfo(farmConf.Staker, record.EventRewardToken, *chainConf, cacher)
			if err != nil {
				logx.Errorf("SetRewardPoolInfo Err:%v ", err)
				return
			}
			if strings.EqualFold(totalWeightString, "0") {
				totalWeightString = "1000"
			}
			rewardTokenLogo := ""
			if !strings.EqualFold(farmConf.XRewardToken, "") {
				rewardTokenLogo = farmConf.XRewardToken
			}
			validStakeRecords, ok := validStakeRecordMap[record.EventIncentiveId]
			if ok {
				tokenIds := []*big.Int{}
				userAddressesMap := map[string]struct{}{}
				linq.From(validStakeRecords).ForEachT(func(item repository.EventFarmTokenStaked) {
					tokenIds = append(tokenIds, big.NewInt(item.EventTokenId))
					userAddressesMap[item.From] = struct{}{}
				})

				staticCaller, err := sonex.NewStaticCaller(l.svcCtx.Config)
				if err != nil {
					logx.Error(fmt.Errorf("NewStaticCaller error: %v", err))
					return
				}

				tokenIdsTotal, err := getIncentiveInfoTokenIds(farmConf.Staker, record.EventIncentiveId, *chainConf)
				if err != nil {
					logx.Error(fmt.Errorf("getIncentiveInfoTokenIds error: %v", err))
					return
				}
				positions, err := staticCaller.BatchGetPositions(l.ctx, req.ChainId, common.HexToAddress(sonexConf.PositionManager), tokenIdsTotal)
				if err != nil {
					logx.Error(fmt.Errorf("BatchGetPositions error: %v", err))
					return
				}
				totalLiquidity := big.NewInt(0)
				totalLiquidityUSD := big.NewFloat(0)
				for _, position := range positions {
					totalLiquidity.Add(totalLiquidity, position.Liquidity)
					amount0, amount1 := position.Amounts()
					totalLiquidityUSD.Add(totalLiquidityUSD, new(big.Float).Quo(new(big.Float).Mul(new(big.Float).SetInt(amount0), _price0divUsdDecimals), new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(token0.Decimals), nil))))
					totalLiquidityUSD.Add(totalLiquidityUSD, new(big.Float).Quo(new(big.Float).Mul(new(big.Float).SetInt(amount1), _price1divUsdDecimals), new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(token1.Decimals), nil))))
				}

				apr := getApr("listok", record.EventPool, record.EventIncentiveId, totalRewardUSD, totalLiquidityUSD, record.EventWeight, farmConf.Staker, record.ExtraRewardToken, *chainConf, cacher)

				var tids []int64
				linq.From(tokenIds).SelectT(func(item *big.Int) int64 { return item.Int64() }).ToSlice(&tids)

				totalLiquidityUSDValue, _ := totalLiquidityUSD.Float64()

				dataResult = types.FarmIncentivesResponse_Data{
					IncentiveId:                  record.EventIncentiveId,
					PoolAddress:                  record.EventPool,
					FeeRatio:                     record.ExtraFee,
					RewardTokenAddress:           rewardToken.Address,
					RewardTokenDecimals:          rewardToken.Decimals,
					RewardTokenSymbol:            rewardToken.Symbol,
					RewardTokenPrice:             rewardPrice.String(),
					Token0Address:                token0.Address,
					Token0Decimals:               token0.Decimals,
					Token0Symbol:                 token0.Symbol,
					Token0Price:                  price0.String(),
					Token1Address:                token1.Address,
					Token1Decimals:               token1.Decimals,
					Token1Symbol:                 token1.Symbol,
					Token1Price:                  price1.String(),
					Liquidity:                    totalLiquidity.String(),
					LiquidityUSD:                 totalLiquidityUSD.String(),
					LiquidityUSDFloat64:          totalLiquidityUSDValue,
					TotalClaimableFeesUSDFloat64: totalClaimableFeesUSDFloat64,
					Weight:                       eventWeightString,
					TotalWeight:                  totalWeightString,
					APR:                          apr,
					JoinNum:                      int64(len(userAddressesMap)),
					Status:                       status,
					StartAt:                      record.EventStartTime,
					Earned:                       big.NewInt(0).String(),
					RewardTokenLogo:              rewardTokenLogo,
				}
				tokenIdsResult = tids
			} else {
				// apr := new(big.Float).Quo(totalRewardUSD, big.NewFloat(1))
				// 0 的时候  就应该按照1来算
				// 如果 totalliquidity 为0,那么按照 totalliquidity = 1 U 计算
				apr := getApr("listnotoknotok", record.EventPool, record.EventIncentiveId, totalRewardUSD, big.NewFloat(1), record.EventWeight, farmConf.Staker, record.ExtraRewardToken, *chainConf, cacher)

				dataResult = types.FarmIncentivesResponse_Data{
					IncentiveId:                  record.EventIncentiveId,
					PoolAddress:                  record.EventPool,
					FeeRatio:                     record.ExtraFee,
					RewardTokenAddress:           rewardToken.Address,
					RewardTokenDecimals:          rewardToken.Decimals,
					RewardTokenSymbol:            rewardToken.Symbol,
					RewardTokenPrice:             rewardPrice.String(),
					Token0Address:                token0.Address,
					Token0Decimals:               token0.Decimals,
					Token0Symbol:                 token0.Symbol,
					Token0Price:                  price0.String(),
					Token1Address:                token1.Address,
					Token1Decimals:               token1.Decimals,
					Token1Symbol:                 token1.Symbol,
					Token1Price:                  price1.String(),
					Liquidity:                    big.NewInt(0).String(),
					LiquidityUSD:                 big.NewInt(0).String(),
					TotalClaimableFeesUSDFloat64: totalClaimableFeesUSDFloat64,
					TotalClaimableFeesUSD:        totalClaimableFeesUSD.String(),
					Weight:                       eventWeightString,
					TotalWeight:                  totalWeightString,
					APR:                          apr,
					JoinNum:                      0,
					Status:                       status,
					StartAt:                      record.EventStartTime,
					Earned:                       big.NewInt(0).String(),
					RewardTokenLogo:              rewardTokenLogo,
				}
			}

			ch <- struct {
				Data     types.FarmIncentivesResponse_Data
				TokenIds []int64
				Err      error
			}{
				Data:     dataResult,
				TokenIds: tokenIdsResult,
				Err:      err,
			}
		}(&incincentiveCreatedRecordItem, results)
	}
	wg.Wait()
	close(results)

	list := []types.FarmIncentivesResponse_Data{}
	tokenIds = map[int64][]int64{}
	for result := range results {
		if result.Err != nil {
			logx.Errorf("result error:%v", result.Err)
			return resp, tokenIds, result.Err
		}
		list = append(list, result.Data)
		tokenIds[result.Data.IncentiveId] = result.TokenIds
	}

	return &types.FarmIncentivesResponse{
		USD: types.FarmIncentivesResponse_Token{
			Address:  usdToken.Address,
			Decimals: usdToken.Decimals,
			Symbol:   usdToken.Symbol,
			Price:    usdAccuracy.String(),
		},
		List: list,
	}, tokenIds, nil
}

type FarmIncentivesOperationEventType int64

const (
	FarmIncentivesOperationEventTypeStaked FarmIncentivesOperationEventType = iota
	FarmIncentivesOperationEventTypeUnstaked
)

type FarmIncentivesOperationEvent struct {
	FarmIncentivesOperationEventType FarmIncentivesOperationEventType
	StakedEvent                      *repository.EventFarmTokenStaked
	UnstakedEvent                    *repository.EventFarmTokenUnstaked
}

// GetFarmIncentivesOperationEvents 获取 user address 质押/解质押事件, key 为 incentive id
func (l *FarmIncentivesLogic) GetFarmIncentivesOperationEvents(chainId int64, userAddress string) (m map[int64][]*FarmIncentivesOperationEvent, err error) {
	type stakedResult struct {
		events []*repository.EventFarmTokenStaked
		err    error
	}
	type unstakedResult struct {
		events []*repository.EventFarmTokenUnstaked
		err    error
	}

	stakedChan := make(chan stakedResult, 1)
	unstakedChan := make(chan unstakedResult, 1)

	// Concurrent fetching using channels
	go func() {
		events, err := l.eventFarmTokenStakedModel.FindByChainIdUser(l.ctx, chainId, userAddress)
		stakedChan <- stakedResult{events: events, err: err}
	}()

	go func() {
		events, err := l.eventFarmTokenUnstakedModel.FindByChainIdUser(l.ctx, chainId, userAddress)
		unstakedChan <- unstakedResult{events: events, err: err}
	}()

	// Receive results from channels
	stakedRes := <-stakedChan
	unstakedRes := <-unstakedChan

	if stakedRes.err != nil {
		logx.Errorf("GetFarmIncentivesOperationEvents: GetStakedEvents error: %v", stakedRes.err)
		return nil, stakedRes.err
	}

	if unstakedRes.err != nil {
		logx.Errorf("GetFarmIncentivesOperationEvents: GetUnstakedEvents error: %v", unstakedRes.err)
		return nil, unstakedRes.err
	}

	// Initialize the map
	m = make(map[int64][]*FarmIncentivesOperationEvent)

	// Process staked events
	for _, stakedEvent := range stakedRes.events {
		m[stakedEvent.EventIncentiveId] = append(m[stakedEvent.EventIncentiveId], &FarmIncentivesOperationEvent{
			FarmIncentivesOperationEventType: FarmIncentivesOperationEventTypeStaked,
			StakedEvent:                      stakedEvent,
		})
	}

	// Process unstaked events
	for _, unstakedEvent := range unstakedRes.events {
		m[unstakedEvent.EventIncentiveId] = append(m[unstakedEvent.EventIncentiveId], &FarmIncentivesOperationEvent{
			FarmIncentivesOperationEventType: FarmIncentivesOperationEventTypeUnstaked,
			UnstakedEvent:                    unstakedEvent,
		})
	}

	return m, nil
}

// GetFarmIncentives 获取有效并且当前正属于开启状态的 farm incentive pool列表
func (l *FarmIncentivesLogic) GetFarmIncentives(chainId int64) (m []repository.EventFarmIncentiveCreatedWithExtra, err error) {
	farmConf := l.svcCtx.Config.GetFarmConfig(chainId)
	if farmConf == nil {
		return nil, fmt.Errorf("%s", "Miss farm configure")
	}

	records, err := l.eventIncentiveCreatedModel.QueryByChainIdStakerWithExtra(l.ctx, chainId, common.HexToAddress(farmConf.Staker).Hex())
	if err != nil {
		logx.Errorf("GetFarmIncentivesPools: QueryByChainIdStakerWithExtra error: %v", err)
		return nil, err
	}
	return records, nil
}
