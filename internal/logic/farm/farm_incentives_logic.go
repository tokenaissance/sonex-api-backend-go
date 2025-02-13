package farm

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"sort"
	"strings"
	"sync"
	"time"

	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/cache"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/module/farm"
	"sonex_api_backend_go/internal/module/sonex"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/ahmetb/go-linq/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type FarmIncentivesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

type CacheWarpper struct {
	List       []types.FarmIncentivesResponse_Data
	TokenIds   map[int64][]int64
	LastUpdate time.Time
}

// Farm incentive list
func NewFarmIncentivesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FarmIncentivesLogic {
	return &FarmIncentivesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FarmIncentivesLogic) FarmIncentives(req *types.FarmIncentivesRequest, positions []*types.Position) (resp *types.FarmIncentivesResponse, err error) {
	if req.ChainId <= 0 {
		return nil, fmt.Errorf("invalid parameter")
	}

	farmConf := l.svcCtx.Config.GetFarmConfig(req.ChainId)
	chain := l.svcCtx.Config.GetChainConfig(req.ChainId)
	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	sonexConf := l.svcCtx.Config.GetSonexConfig(req.ChainId)

	cacheKey := fmt.Sprintf("farm:pools:%d", req.ChainId)
	redisCache := cache.NewRedisCache(fmt.Sprintf("%s:%d", l.svcCtx.Config.Redis.Host, l.svcCtx.Config.Redis.Port), l.svcCtx.Config.Redis.Password, l.svcCtx.Config.Redis.DB)

	tokenIds := map[int64][]int64{}
	// Check cache
	lastUpdate, err := redisCache.Get(cacheKey)
	if err == nil && lastUpdate != "" {
		var cacheWarpper CacheWarpper
		err := json.Unmarshal([]byte(lastUpdate), &cacheWarpper)
		if err == nil {
			resp = &types.FarmIncentivesResponse{
				List: cacheWarpper.List,
			}
			tokenIds = cacheWarpper.TokenIds
			var respLock sync.Mutex
			eventFarmIncentiveCreatedModel := repository.NewEventFarmIncentiveCreatedModel(conn)
			eventFarmIncentiveEndedModel := repository.NewEventFarmIncentiveEndedModel(conn)
			eventFarmTokenStakedModel := repository.NewEventFarmTokenStakedModel(conn)
			eventFarmTokenUnstakedModel := repository.NewEventFarmTokenUnstakedModel(conn)
			var wg sync.WaitGroup

			wg.Add(1)
			go func() {
				defer wg.Done()

				record, err := eventFarmIncentiveCreatedModel.FindLatestByChainId(l.ctx, req.ChainId)
				switch err {
				case repository.ErrNotFound:
					return
				case nil:
					if record.TxAt.UTC().After(cacheWarpper.LastUpdate) {
						respLock.Lock()
						defer respLock.Unlock()
						resp = nil
					}
				default:
					logx.Errorf("FindLatestByChainId error: %v", err)
					return
				}
			}()
			wg.Add(1)
			go func() {
				defer wg.Done()

				record, err := eventFarmIncentiveEndedModel.FindLatestByChainId(l.ctx, req.ChainId)
				switch err {
				case repository.ErrNotFound:
					return
				case nil:
					if record.TxAt.UTC().After(cacheWarpper.LastUpdate) {
						respLock.Lock()
						defer respLock.Unlock()
						resp = nil
					}
				default:
					logx.Errorf("FindLatestByChainId error: %v", err)
					return
				}
			}()
			wg.Add(1)
			go func() {
				defer wg.Done()

				record, err := eventFarmTokenStakedModel.FindLatestByChainId(l.ctx, req.ChainId)
				switch err {
				case repository.ErrNotFound:
					return
				case nil:
					if record.TxAt.UTC().After(cacheWarpper.LastUpdate) {
						respLock.Lock()
						defer respLock.Unlock()
						resp = nil
					}
				default:
					logx.Errorf("FindLatestByChainId error: %v", err)
					return
				}
			}()
			wg.Add(1)
			go func() {
				defer wg.Done()

				record, err := eventFarmTokenUnstakedModel.FindLatestByChainId(l.ctx, req.ChainId)
				switch err {
				case repository.ErrNotFound:
					return
				case nil:
					if record.TxAt.UTC().After(cacheWarpper.LastUpdate) {
						respLock.Lock()
						defer respLock.Unlock()
						resp = nil
					}
				default:
					logx.Errorf("FindLatestByChainId error: %v", err)
					return
				}
			}()
			wg.Wait()
		}
	}

	// No cache
	// 查看是否过期
	if resp == nil {
		resp, tokenIds, err = l.GetResponse(req, chain, conn, farmConf, sonexConf, redisCache)
		if len(tokenIds) == 0 {
			logx.Info("no tokenid founds")
			return nil, nil
		}
		if err != nil {
			logx.Errorf("ERR from GetResponse:%v", err)
			return nil, err
		}
		// Set cache
		// TODO  not in cache
		// respBytes, err := json.Marshal(CacheWarpper{List: resp.List, TokenIds: tokenIds, LastUpdate: time.Now().UTC()})
		// if err != nil {
		// 	logx.Errorf("Marshal response error: %v", err)
		// }
		// err = redisCache.Set(cacheKey, respBytes, time.Hour*24)
		// if err != nil {
		// 	logx.Errorf("Set cache error: %v", err)
		// }
	}

	// Filter
	if !strings.EqualFold(req.UserAddress, "") && common.IsHexAddress(req.UserAddress) {

		// loop assgin
		// TODO 优化
		if len(positions) > 0 {
			if len(resp.List) > 0 {
				for i, item := range resp.List {
					for _, vv := range positions {
						if strings.EqualFold(vv.Pool.Address, item.PoolAddress) {
							resp.List[i].Positions = append(resp.List[i].Positions, vv)
						}
					}
					// 可用余额
					resp.List[i].Balance = int64(len(resp.List[i].Positions))
					logx.Infof("RRRRRRRRRRRRRRRR:%v IncentiveId:%v ", len(resp.List[i].Positions), resp.List[i].IncentiveId)
				}
			}
		}

		usdAccuracy, _ := new(big.Float).SetString(resp.USD.Price)

		_ = usdAccuracy
		eventFarmTokenStakedModel := repository.NewEventFarmTokenStakedModel(conn)

		validStakeRecords, err := eventFarmTokenStakedModel.QueryValidByChainIdUserWithExtra(l.ctx, req.ChainId, common.HexToAddress(req.UserAddress).Hex())
		if err != nil {
			logx.Errorf("ERR:%v", err)
			return nil, err
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
		statcCaller, err := farm.NewStaticCaller(l.svcCtx.Config)
		if err != nil {
			logx.Errorf("ERR:%v", err)
			return nil, err
		}
		unclaimeds, err := statcCaller.BatchGetUnclaimeds(l.ctx, req.ChainId, common.HexToAddress(req.UserAddress), incentiveIds)
		if err != nil {
			logx.Errorf("ERR:%v", err)
			return nil, err
		}
		for i, item := range resp.List {
			if val, ok := staked[item.IncentiveId]; ok {
				resp.List[i].Staked = val

				// resp.List[i].Balance = int64(len(userTokenIds[item.IncentiveId])) - val
				// resp.List[i].Balance = int64(len(resp.List[i].Positions)) - val
				rewards, err := getRewardsAmount(farmConf.Staker, req.UserAddress, item.IncentiveId, *chain)
				if err != nil {
					return resp, err
				}

				logx.Infof("rewards :%v ", rewards.String())
				totalClaimableFeesUSD, _, err := getValueOfAmountToken(rewards, item.IncentiveId, req.ChainId, chain, conn)
				if err != nil {
					return resp, err
				}
				resp.List[i].TotalClaimableFeesUSD = totalClaimableFeesUSD.String()

				totalClaimableFeesUSDFloat64, _ := totalClaimableFeesUSD.Float64()
				resp.List[i].TotalClaimableFeesUSDFloat64 = totalClaimableFeesUSDFloat64
			} else {
				// unstaked
				// resp.List[i].Balance = int64(len(resp.List[i].Positions))
			}

			logx.Infof("SSSSA:%v B:%v c:%v", farmConf.Staker, req.UserAddress, item.IncentiveId)

			totalClaimedAmount, err := getClainedAmount(req.UserAddress, req.ChainId, item.IncentiveId, conn)
			if err != nil {
				return resp, err
			}

			totalClaimedFeesUSD_bigFloat, _, err := getValueOfAmountToken(totalClaimedAmount, item.IncentiveId, req.ChainId, chain, conn)

			if err != nil {
				return resp, err
			}
			resp.List[i].TotalClaimedFeesUSD = totalClaimedFeesUSD_bigFloat.String()
			_ = unclaimeds
			// if val, ok := unclaimeds[item.IncentiveId]; ok {
			// rewardTokenPrice, _ := new(big.Float).SetString(item.RewardToken.Price)
			// rewardTokenAccuracy := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(item.RewardToken.Decimals), nil))
			// unclaimedUSD := new(big.Float).Quo(new(big.Float).Quo(new(big.Float).Mul(new(big.Float).SetInt(val), rewardTokenPrice), rewardTokenAccuracy), usdAccuracy)
			// resp.List[i].EarnedUSD, _ = unclaimedUSD.Float64()
			// }
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
			linq.From(resp.List).OrderByDescendingT(func(item types.FarmIncentivesResponse_Data) int64 {
				return item.JoinNum
			}).ToSlice(&resp.List)
		case "apr":
			// linq.From(resp.List).OrderByDescendingT(func(item types.FarmIncentivesResponse_Data) lib_common.CompareableBigInt {
			// 	apr, _ := new(big.Int).SetString(item.APR, 10)
			// 	return lib_common.CompareableBigInt{Value: apr}
			// }).ToSlice(&resp.List)
			sort.SliceStable(resp.List, func(i, j int) bool {
				apr_i, _ := big.NewFloat(0).SetString(resp.List[i].APR)
				apr_j, _ := big.NewFloat(0).SetString(resp.List[j].APR)
				return apr_i.Cmp(apr_j) > 0
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
		return nil, nil, err
	}
	usdAccuracy := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(usdToken.Decimals), nil))

	eventFarmIncentiveCreatedModel := repository.NewEventFarmIncentiveCreatedModel(conn)
	eventFarmTokenStakedModel := repository.NewEventFarmTokenStakedModel(conn)
	tokenPriceModel := repository.NewSonexTokenPriceRecordModel(conn)

	incentiveCreatedRecords, err := eventFarmIncentiveCreatedModel.QueryByChainIdStakerWithExtra(l.ctx, req.ChainId, common.HexToAddress(farmConf.Staker).Hex())
	if err != nil {
		logx.Errorf("ERR1:%v ", err)
		return nil, nil, err
	}
	logx.Infof("incentiveCreatedRecords:%v len :%v  ", incentiveCreatedRecords, len(incentiveCreatedRecords))
	validStakeRecords, err := eventFarmTokenStakedModel.QueryValidByChainId(l.ctx, req.ChainId)
	if err != nil {
		logx.Errorf("ERR2:%v ", err)
		return nil, nil, err
	}
	validStakeRecordMap := map[int64][]repository.EventFarmTokenStaked{}
	logx.Infof("validStakeRecords:%v ", validStakeRecords)
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

	logx.Infof("validStakeRecordMap:%v ", validStakeRecordMap)
	prices, err := tokenPriceModel.QueryLatestHistoryPrices(l.ctx)
	if err != nil {
		logx.Errorf("ERR3:%v ", err)
		return nil, nil, err
	}

	logx.Infof("prices:%#v ", prices)
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
			price0 := prices[token0.Address]
			price1 := prices[token1.Address]

			_price0divUsdDecimals := new(big.Float).Quo(new(big.Float).SetInt(price0), new(big.Float).SetFloat64(math.Pow10(int(6))))
			_price1divUsdDecimals := new(big.Float).Quo(new(big.Float).SetInt(price1), new(big.Float).SetFloat64(math.Pow10(int(6))))

			// price0BigFloat :=  big.NewFloat()

			rewardToken, err := global.GetToken(req.ChainId, record.EventRewardToken)
			if err != nil {
				logx.Errorf("%s", fmt.Sprintf("Get extra token error: %s %s", record.EventRewardToken, err))
				return
			}
			logx.Infof("rewardToken:%#v", rewardToken)
			ratepersecond, _ := new(big.Float).SetString(record.ExtraRate)
			logx.Info("333333:", rewardToken.Address, "prices[rewardToken.Address]:", prices[rewardToken.Address])
			rewardPrice := new(big.Float).SetInt(prices[rewardToken.Address])
			logx.Info("4444444:")
			// rewardTokenAccuracy := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(rewardToken.Decimals), nil))
			rewardTokenAccuracy := new(big.Float).SetFloat64(math.Pow10(int(rewardToken.Decimals)))
			logx.Info("555555:")

			rewardUSD := new(big.Float).Quo(new(big.Float).Quo(new(big.Float).Mul(ratepersecond, rewardPrice), rewardTokenAccuracy), usdAccuracy)

			totalRewardUSD := new(big.Float).Mul(rewardUSD, big.NewFloat(365*24*60*60))
			logx.Info("77777777:")
			logx.Infof("666666:rewardToken:%#v ratepersecond :%v  rewardPrice:%v totalRewardUSD365: %v  ", rewardToken, ratepersecond, rewardPrice, totalRewardUSD)

			err = global.SetTotalReward365ForIncentiveId(cacher, global.Reward365ForIncentiveIdKey, incincentiveCreatedRecordItem.EventIncentiveId, totalRewardUSD.String())
			if err != nil {
				logx.Errorf("ERR:%v ", err)
			}
			var dataResult types.FarmIncentivesResponse_Data
			tokenIdsResult := []int64{}
			validStakeRecords, ok := validStakeRecordMap[record.EventIncentiveId]
			logx.Info("1111111111111111111111:", ok)
			var status string
			if record.ExtraEnded {
				status = "closed"
			} else {
				status = "active"
			}

			rewards, err := getRewardsAmount(farmConf.Staker, req.UserAddress, record.EventIncentiveId, *chainConf)
			if err != nil {
				logx.Errorf("getRewardsAmount:%v ", err)
			}

			logx.Infof("rewards :%v ", rewards.String())
			totalClaimableFeesUSD, _, err := getValueOfAmountToken(rewards, record.EventIncentiveId, req.ChainId, chainConf, conn)

			totalClaimableFeesUSDFloat64, _ := totalClaimableFeesUSD.Float64()

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
				logx.Infof("tokenIdsTotal:%v tokenIds:%v ", tokenIdsTotal, tokenIds)
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

				apr := getApr(record.EventPool, record.EventIncentiveId, totalRewardUSD, totalLiquidityUSD, record.EventWeight, farmConf.Staker, record.ExtraRewardToken, *chainConf, cacher)

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
					APR:                          apr,
					JoinNum:                      int64(len(userAddressesMap)),
					Status:                       status,
					StartAt:                      record.EventStartTime,
					Earned:                       big.NewInt(0).String(),
				}
				tokenIdsResult = tids
			} else {
				logx.Info("notoknotok")
				// apr := new(big.Float).Quo(totalRewardUSD, big.NewFloat(1))
				// 0 的时候  就应该按照1来算
				// 如果 totalliquidity 为0,那么按照 totalliquidity = 1 U 计算
				apr := getApr(record.EventPool, record.EventIncentiveId, totalRewardUSD, big.NewFloat(1), record.EventWeight, farmConf.Staker, record.ExtraRewardToken, *chainConf, cacher)

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
					APR:                          apr,
					JoinNum:                      0,
					Status:                       status,
					StartAt:                      record.EventStartTime,
					Earned:                       big.NewInt(0).String(),
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
			logx.Errorf("result error:%v", err)
			return nil, nil, result.Err
		}
		list = append(list, result.Data)
		tokenIds[result.Data.IncentiveId] = result.TokenIds
	}

	logx.Info(" tokenIds: ", tokenIds)

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
