package farm

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/lib/cache"
	lib_common "sonex_api_backend_go/internal/lib/common"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/module/sonex"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"
	"sync"
	"time"

	"sonex_api_backend_go/internal/repository"

	"github.com/ahmetb/go-linq/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type FarmPoolListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

type CacheWarpper struct {
	List       []types.FarmPoolData
	TokenIds   map[string][]int64
	LastUpdate time.Time
}

// Farm pool list
func NewFarmPoolListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FarmPoolListLogic {
	logger := logx.WithContext(ctx).WithFields(logx.Field("module", "Farm"))
	return &FarmPoolListLogic{
		Logger: logger,
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FarmPoolListLogic) FarmPoolList(req *types.FarmPoolListRequest) (resp *types.FarmPoolListResponse, err error) {
	if req.ChainId <= 0 {
		return nil, fmt.Errorf("invalid parameter")
	}
	if req.PageSize < 1 {
		req.PageSize = 20
	}
	if req.PageIndex < 1 {
		req.PageIndex = 1
	}

	cacheKey := fmt.Sprintf("farm:pools:%d", req.ChainId)
	redisCache := cache.NewRedisCache(fmt.Sprintf("%s:%d", l.svcCtx.Config.Redis.Host, l.svcCtx.Config.Redis.Port), l.svcCtx.Config.Redis.Password, l.svcCtx.Config.Redis.DB)

	tokenIds := map[string][]int64{}
	// Check cache
	lastUpdate, err := redisCache.Get(cacheKey)
	if err == nil && lastUpdate != "" {
		var cacheWarpper CacheWarpper
		err := json.Unmarshal([]byte(lastUpdate), &cacheWarpper)
		if err == nil {
			resp = &types.FarmPoolListResponse{
				List: cacheWarpper.List,
			}
			tokenIds = cacheWarpper.TokenIds
			var respLock sync.Mutex
			conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
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
	if resp == nil {
		resp, tokenIds, err = l.getResponse(req)
		if err != nil {
			return nil, err
		}

		// Set cache
		respBytes, err := json.Marshal(CacheWarpper{List: resp.List, TokenIds: tokenIds, LastUpdate: time.Now().UTC()})
		if err != nil {
			logx.Errorf("Marshal response error: %v", err)
		}
		err = redisCache.Set(cacheKey, respBytes, time.Hour*24)
		if err != nil {
			logx.Errorf("Set cache error: %v", err)
		}
	}

	// Update status
	linq.From(resp.List).SelectT(func(item types.FarmPoolData) types.FarmPoolData {
		endTime := time.Unix(item.EndAt, 0)
		if endTime.Before(time.Now()) {
			item.Status = "closed"
		} else {
			item.Status = "active"
		}
		return item
	}).ToSlice(&resp.List)

	// Filter
	if req.UserAddress != "" && common.IsHexAddress(req.UserAddress) {
		conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
		eventFarmTokenStakedModel := repository.NewEventFarmTokenStakedModel(conn)
		validStakeRecords, err := eventFarmTokenStakedModel.QueryValidByChainIdUserAddress(l.ctx, req.ChainId, common.HexToAddress(req.UserAddress).Hex())
		if err != nil {
			return nil, err
		}
		userTokenIds := []int64{}
		linq.From(validStakeRecords).ForEachT(func(item repository.EventFarmTokenStaked) {
			userTokenIds = append(userTokenIds, item.EventTokenId)
		})
		linq.From(resp.List).WhereT(func(item types.FarmPoolData) bool {
			poolTokenIds, ok := tokenIds[item.PoolAddress]
			if !ok {
				return false
			}
			return linq.From(userTokenIds).AnyWithT(func(tokenId int64) bool {
				return linq.From(poolTokenIds).Contains(tokenId)
			})
		}).ToSlice(&resp.List)
	}
	if req.Status != "" && (req.Status == "active" || req.Status == "closed") {
		linq.From(resp.List).WhereT(func(item types.FarmPoolData) bool {
			return item.Status == req.Status
		}).ToSlice(&resp.List)
	}
	if req.OrderBy != "" {
		switch req.Status {
		case "hot":
			linq.From(resp.List).OrderByDescendingT(func(item types.FarmPoolData) int64 {
				return item.JoinNum
			}).ToSlice(&resp.List)
		case "apr":
			linq.From(resp.List).OrderByDescendingT(func(item types.FarmPoolData) lib_common.CompareableBigInt {
				apr, _ := new(big.Int).SetString(item.APR, 10)
				return lib_common.CompareableBigInt{Value: apr}
			}).ToSlice(&resp.List)
		}
	}

	// Page
	linq.From(resp.List).Skip(int((req.PageIndex - 1) * req.PageSize)).Take(int(req.PageSize)).ToSlice(&resp.List)

	return resp, nil
}

func (l *FarmPoolListLogic) getResponse(req *types.FarmPoolListRequest) (resp *types.FarmPoolListResponse, tokenIds map[string][]int64, err error) {
	sonexConf := l.svcCtx.Config.GetSonexConfig(req.ChainId)
	framConf := l.svcCtx.Config.GetFarmConfig(req.ChainId)
	stablecoin := common.HexToAddress(l.svcCtx.Config.GetChainConfig(req.ChainId).USD)
	stablecoinToken, err := global.GetToken(req.ChainId, stablecoin.Hex())
	if err != nil {
		return nil, nil, err
	}

	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	eventFarmIncentiveCreatedModel := repository.NewEventFarmIncentiveCreatedModel(conn)
	eventFarmTokenStakedModel := repository.NewEventFarmTokenStakedModel(conn)
	tokenPriceModel := repository.NewSonexTokenPriceRecordModel(conn)

	incentiveCreatedRecords, err := eventFarmIncentiveCreatedModel.QueryByChainIdStakerWithPool(l.ctx, req.ChainId, common.HexToAddress(framConf.Staker).Hex())
	if err != nil {
		return nil, nil, err
	}
	validStakeRecords, err := eventFarmTokenStakedModel.QueryValidByChainId(l.ctx, req.ChainId)
	if err != nil {
		return nil, nil, err
	}
	validStakeRecordMap := map[string][]repository.EventFarmTokenStaked{}
	linq.From(validStakeRecords).
		GroupByT(
			func(item repository.EventFarmTokenStaked) string { return item.EventIncentiveId },
			func(item repository.EventFarmTokenStaked) repository.EventFarmTokenStaked { return item },
		).
		ToMapByT(
			&validStakeRecordMap,
			func(item linq.Group) string { return item.Key.(string) },
			func(item linq.Group) []repository.EventFarmTokenStaked {
				value := []repository.EventFarmTokenStaked{}
				for _, v := range item.Group {
					value = append(value, v.(repository.EventFarmTokenStaked))
				}
				return value
			},
		)

	prices, err := tokenPriceModel.GetAllLatestHistoryPrices()
	if err != nil {
		return nil, nil, err
	}

	var wg sync.WaitGroup
	results := make(chan struct {
		Data     types.FarmPoolData
		TokenIds []int64
	}, len(incentiveCreatedRecords))

	for _, incincentiveCreatedRecord := range incentiveCreatedRecords {
		wg.Add(1)
		go func(record *repository.EventFarmIncentiveCreatedWithPool, ch chan struct {
			Data     types.FarmPoolData
			TokenIds []int64
		}) {
			defer wg.Done()

			token0, err := global.GetToken(req.ChainId, record.Token0)
			if err != nil {
				logx.Errorf("GetToken error: %v", err)
				return
			}
			token1, err := global.GetToken(req.ChainId, record.Token1)
			if err != nil {
				logx.Errorf("GetToken error: %v", err)
				return
			}
			price0 := prices[token0.Address]
			price1 := prices[token1.Address]

			validStakeRecords, ok := validStakeRecordMap[record.EventIncentiveId]
			if !ok {
				logx.Errorf("ValidStakeRecord not found: %s", record.EventIncentiveId)
				return
			}
			tokenIds := []*big.Int{}
			userAddressesMap := map[string]struct{}{}
			linq.From(validStakeRecords).ForEachT(func(item repository.EventFarmTokenStaked) {
				tokenIds = append(tokenIds, big.NewInt(item.EventTokenId))
				userAddressesMap[item.From] = struct{}{}
			})

			staticCaller, err := sonex.NewStaticCaller(l.svcCtx.Config)
			if err != nil {
				logx.Errorf("NewStaticCaller error: %v", err)
				return
			}
			positions, err := staticCaller.BatchGetPositions(l.ctx, req.ChainId, common.HexToAddress(sonexConf.PositionManager), tokenIds)
			if err != nil {
				logx.Errorf("BatchGetPositions error: %v", err)
				return
			}
			totalLiquidity := big.NewInt(0)
			totalLiquidityUSD := big.NewInt(0)
			for _, position := range positions {
				totalLiquidity.Add(totalLiquidity, position.Liquidity)
				amount0, amount1 := position.Amounts()
				totalLiquidityUSD.Add(totalLiquidityUSD, new(big.Int).Div(new(big.Int).Mul(amount0, price0), new(big.Int).Exp(big.NewInt(10), big.NewInt(stablecoinToken.Decimals), nil)))
				totalLiquidityUSD.Add(totalLiquidityUSD, new(big.Int).Div(new(big.Int).Mul(amount1, price1), new(big.Int).Exp(big.NewInt(10), big.NewInt(stablecoinToken.Decimals), nil)))
			}
			reward, _ := new(big.Int).SetString(record.EventReward, 10)
			startTime := time.Unix(record.EventStartTime, 0)
			endTime := time.Unix(record.EventEndTime, 0)
			apr := new(big.Int).Div(new(big.Int).Mul(reward, big.NewInt(365*24*60*60)), big.NewInt(endTime.Unix()-startTime.Unix()))
			var tids []int64
			linq.From(tokenIds).SelectT(func(item *big.Int) int64 { return item.Int64() }).ToSlice(&tids)
			ch <- struct {
				Data     types.FarmPoolData
				TokenIds []int64
			}{
				Data: types.FarmPoolData{
					ChainId:        incincentiveCreatedRecord.ChainId,
					PoolAddress:    record.EventPool,
					FeeRatio:       record.Fee,
					Token0Address:  token0.Address,
					Token0Decimals: token0.Decimals,
					Token0Symbol:   token0.Symbol,
					Token0Price:    price0.String(),
					Token1Address:  token1.Address,
					Token1Decimals: token1.Decimals,
					Token1Symbol:   token1.Symbol,
					Token1Price:    price1.String(),
					Liquidity:      totalLiquidity.String(),
					LiquidityUSD:   totalLiquidityUSD.String(),
					APR:            apr.String(),
					JoinNum:        int64(len(userAddressesMap)),
					StartAt:        record.EventStartTime,
					EndAt:          record.EventEndTime,
				},
				TokenIds: tids,
			}
		}(&incincentiveCreatedRecord, results)
	}
	wg.Wait()
	close(results)

	list := []types.FarmPoolData{}
	tokenIds = map[string][]int64{}
	for result := range results {
		list = append(list, result.Data)
		tokenIds[result.Data.PoolAddress] = result.TokenIds
	}

	return &types.FarmPoolListResponse{
		List: list,
	}, tokenIds, nil
}
