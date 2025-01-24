package sonex

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"sync"
	"time"

	"sonex_api_backend_go/internal/handler/response"
	"sonex_api_backend_go/internal/lib/cache"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/ahmetb/go-linq/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type TokensLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

type CacheWarpper struct {
	Data       []types.TokensResponse_Token
	LastUpdate time.Time
}

func NewTokensLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TokensLogic {
	return &TokensLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TokensLogic) Tokens(req *types.TokensRequest) (resp *response.ApiResponse, err error) {
	if req.ChainId == 0 {
		return nil, fmt.Errorf("%s", "Invalid parameter")
	}

	chainConf := l.svcCtx.Config.GetChainConfig(req.ChainId)
	if chainConf == nil {
		return nil, fmt.Errorf("%s", "Miss chain configure")
	}
	usdToken, err := global.GetToken(req.ChainId, common.HexToAddress(chainConf.USD).Hex())
	if err != nil {
		return nil, err
	}
	usdAccuracy := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(usdToken.Decimals), nil))

	cacheKey := fmt.Sprintf("sonex:tokens:%d", req.ChainId)
	redisCache := cache.NewRedisCache(fmt.Sprintf("%s:%d", l.svcCtx.Config.Redis.Host, l.svcCtx.Config.Redis.Port), l.svcCtx.Config.Redis.Password, l.svcCtx.Config.Redis.DB)

	var dataResult []types.TokensResponse_Token
	// Check cache
	if !req.NoCache {
		lastUpdate, err := redisCache.Get(cacheKey)
		if err == nil && lastUpdate != "" {
			var cacheWarpper CacheWarpper
			err := json.Unmarshal([]byte(lastUpdate), &cacheWarpper)
			if err == nil {
				dataResult = cacheWarpper.Data
				conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
				eventSonexFactoryPoolCreatedModel := repository.NewEventSonexFactoryPoolCreatedModel(conn)
				record, err := eventSonexFactoryPoolCreatedModel.FindLatestByChainId(l.ctx, req.ChainId)
				switch err {
				case repository.ErrNotFound:
				case nil:
					if record.TxAt.UTC().After(cacheWarpper.LastUpdate) {
						resp = nil
					}
				default:
					return nil, err
				}
			}
		}
	}

	// No cache
	if dataResult == nil {
		data, err := l.getData(req)
		if err != nil {
			return nil, err
		}
		dataResult = data

		// Set cache
		if len(data) > 0 {
			respBytes, err := json.Marshal(CacheWarpper{Data: data, LastUpdate: time.Now().UTC()})
			if err != nil {
				logx.Errorf("Marshal response error: %v", err)
			}
			err = redisCache.Set(cacheKey, respBytes, time.Hour*24)
			if err != nil {
				logx.Errorf("Set cache error: %v", err)
			}
		}
	}

	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	tokenPriceModel := repository.NewSonexTokenPriceRecordModel(conn)
	eventSonexPoolSwapModel := repository.NewEventSonexPoolSwapModel(conn)
	tokenModel := repository.NewTokenModel(conn)

	whitelistTokens, err := tokenModel.QueryByChainId(l.ctx, req.ChainId)
	if err != nil {
		return nil, err
	}
	whitelistTokenMap := map[string]repository.Token{}
	linq.From(whitelistTokens).ToMapByT(
		&whitelistTokenMap,
		func(item repository.Token) string { return item.Address },
		func(item repository.Token) repository.Token { return item },
	)

	// Filter
	linq.From(dataResult).
		WhereT(func(item types.TokensResponse_Token) bool {
			_, ok := whitelistTokenMap[item.Address]
			return ok
		}).
		SelectT(func(item types.TokensResponse_Token) types.TokensResponse_Token {
			tokenInfo := whitelistTokenMap[item.Address]
			item.Logo = tokenInfo.Logo.String
			item.Sort = tokenInfo.Sort
			item.Hidden = tokenInfo.Hidden
			return item
		}).
		WhereT(func(item types.TokensResponse_Token) bool { return !item.Hidden }).ToSlice(&dataResult)

	// Update price
	volumes, err := eventSonexPoolSwapModel.SummarizeVolumesByChainId(l.ctx, req.ChainId)
	if err != nil {
		return nil, err
	}

	var tokenAddressesWithOutETH []common.Address
	linq.From(dataResult).
		WhereT(func(item types.TokensResponse_Token) bool { return item.Address != (common.Address{}).Hex() }).
		SelectT(func(item types.TokensResponse_Token) common.Address { return common.HexToAddress(item.Address) }).
		ToSlice(&tokenAddressesWithOutETH)
	var tokenAddressStrs []string
	linq.From(dataResult).SelectT(func(item types.TokensResponse_Token) string { return item.Address }).ToSlice(&tokenAddressStrs)
	totalSupplies, err := global.BatchGetTotalSupplies(l.ctx, req.ChainId, tokenAddressesWithOutETH)
	if err != nil {
		return nil, err
	}
	currentPriceRecords, err := tokenPriceModel.QueryLatestHistoryPricesByAddresses(l.ctx, tokenAddressStrs)
	if err != nil {
		return nil, err
	}
	last1DPriceRecords, err := tokenPriceModel.QueryLatestHistoryPricesByAddressesInterval(l.ctx, tokenAddressStrs, "1d")
	if err != nil {
		return nil, err
	}
	last1HPriceRecords, err := tokenPriceModel.QueryLatestHistoryPricesByAddressesInterval(l.ctx, tokenAddressStrs, "1h")
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	var results = make(chan struct {
		Data types.TokensResponse_Token
		Err  error
	}, len(dataResult))

	for _, item := range dataResult {
		wg.Add(1)
		go func(item types.TokensResponse_Token, ch chan struct {
			Data types.TokensResponse_Token
			Err  error
		}) {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					ch <- struct {
						Data types.TokensResponse_Token
						Err  error
					}{
						Err: fmt.Errorf("%s %s(panic)", item.Address, err),
					}
				}
			}()

			priceI, ok := currentPriceRecords[item.Address]
			if !ok {
				priceI = big.NewInt(0)
			}
			price, _ := global.PriceToUSD(req.ChainId, priceI)
			item.PriceUSD = fmt.Sprintf("%f", price)

			totalSupply, ok := totalSupplies[common.HexToAddress(item.Address)]
			if !ok {
				totalSupply = big.NewInt(0)
			}
			accuracy := new(big.Int).Exp(big.NewInt(10), big.NewInt(item.Decimals), nil)
			fdv, _ := global.PriceToUSD(req.ChainId, new(big.Int).Div(new(big.Int).Mul(totalSupply, priceI), accuracy))
			item.TotalSupply = totalSupply.String()
			item.FdvUSD = fmt.Sprintf("%f", fdv)

			for i, pool := range item.Pools {
				price0I, ok := currentPriceRecords[pool.Token0.Address]
				if !ok {
					price0I = big.NewInt(0)
				}
				price0, _ := global.PriceToUSD(req.ChainId, price0I)
				item.Pools[i].Token0.PriceUSD = fmt.Sprintf("%f", price0)
				price1I, ok := currentPriceRecords[pool.Token1.Address]
				if !ok {
					price1I = big.NewInt(0)
				}
				price1, _ := global.PriceToUSD(req.ChainId, price1I)
				item.Pools[i].Token1.PriceUSD = fmt.Sprintf("%f", price1)

				totalSupply0, ok := totalSupplies[common.HexToAddress(pool.Token0.Address)]
				if !ok {
					totalSupply0 = big.NewInt(0)
				}
				accuracy0 := new(big.Int).Exp(big.NewInt(10), big.NewInt(pool.Token0.Decimals), nil)
				fdv0, _ := global.PriceToUSD(req.ChainId, new(big.Int).Div(new(big.Int).Mul(totalSupply0, price0I), accuracy0))
				item.Pools[i].Token0.TotalSupply = totalSupply0.String()
				item.Pools[i].Token0.FdvUSD = fmt.Sprintf("%f", fdv0)
				totalSupply1, ok := totalSupplies[common.HexToAddress(pool.Token1.Address)]
				if !ok {
					totalSupply1 = big.NewInt(0)
				}
				accuracy1 := new(big.Int).Exp(big.NewInt(10), big.NewInt(pool.Token1.Decimals), nil)
				fdv1, _ := global.PriceToUSD(req.ChainId, new(big.Int).Div(new(big.Int).Mul(totalSupply1, price1I), accuracy1))
				item.Pools[i].Token1.TotalSupply = totalSupply1.String()
				item.Pools[i].Token1.FdvUSD = fmt.Sprintf("%f", fdv1)
			}

			volumeI, ok := volumes[item.Address]
			if !ok {
				volumeI = big.NewInt(0)
			}
			volumeI = new(big.Int).Div(new(big.Int).Mul(volumeI, priceI), accuracy)
			volume, _ := global.PriceToUSD(req.ChainId, volumeI)
			item.VolumeUSD = fmt.Sprintf("%f", volume)

			kline1D, err := tokenPriceModel.Kline(l.ctx, item.Address, "1d")
			if err != nil {
				ch <- struct {
					Data types.TokensResponse_Token
					Err  error
				}{
					Err: err,
				}
				return
			}
			item.Kline = []types.Kline{}
			currentPrice := big.NewInt(0)
			if currentPriceRecord, ok := currentPriceRecords[item.Address]; ok {
				currentPrice = currentPriceRecord
			}
			last1DPrice, _ := usdAccuracy.Int(nil)
			if last1DPriceRecord, ok := last1DPriceRecords[item.Address]; ok {
				last1DPrice = last1DPriceRecord
			}
			last1HPrice, _ := usdAccuracy.Int(nil)
			if last1HPriceRecord, ok := last1HPriceRecords[item.Address]; ok {
				last1HPrice = last1HPriceRecord
			}

			for _, tick := range kline1D {
				priceI, _ := new(big.Int).SetString(tick.Price, 10)
				price, _ := new(big.Float).Quo(new(big.Float).SetInt(priceI), usdAccuracy).Float64()
				t, _ := time.Parse("2006-01-02 15:04:05", tick.CreatedAt)
				item.Kline = append(item.Kline, types.Kline{
					Price:     fmt.Sprintf("%f", price),
					TimeStamp: t.Unix(),
				})
			}
			comparison1H, percentage1HF := l.comparePrices(currentPrice, last1HPrice)
			percentage1H, _ := percentage1HF.Float64()
			comparison1D, percentage1DF := l.comparePrices(currentPrice, last1DPrice)
			percentage1D, _ := percentage1DF.Float64()
			item.RiseFall1H = types.Comparison{Comparison: comparison1H, Percentage: fmt.Sprintf("%f", percentage1H)}
			item.RiseFall1D = types.Comparison{Comparison: comparison1D, Percentage: fmt.Sprintf("%f", percentage1D)}

			ch <- struct {
				Data types.TokensResponse_Token
				Err  error
			}{
				Data: item,
			}
		}(item, results)
	}
	wg.Wait()
	close(results)

	// Handle results
	dataResult = []types.TokensResponse_Token{}
	var eth *types.TokensResponse_Token
	var weth *types.TokensResponse_Token
	for result := range results {
		if result.Err != nil {
			return nil, result.Err
		}
		switch result.Data.Address {
		case (common.Address{}).Hex():
			eth = &result.Data
			continue
		case common.HexToAddress(chainConf.WETH).Hex():
			weth = &result.Data
		}
		dataResult = append(dataResult, result.Data)
	}

	// ETH
	if eth != nil && weth != nil {
		eth.TotalSupply = weth.TotalSupply
		eth.FdvUSD = weth.FdvUSD
		eth.PriceUSD = weth.PriceUSD
		eth.Pools = weth.Pools
		eth.VolumeUSD = weth.VolumeUSD
		eth.Kline = weth.Kline
		eth.RiseFall1H = weth.RiseFall1H
		eth.RiseFall1D = weth.RiseFall1D
		dataResult = append(dataResult, *eth)
	}

	// Sorts
	var query linq.OrderedQuery
	if req.OrderBy != "" {
		switch req.OrderBy {
		case "1H_ASC":
			query = linq.From(dataResult).OrderByT(func(item types.TokensResponse_Token) float64 {
				value, _ := strconv.ParseFloat(item.RiseFall1H.Percentage, 64)
				return value
			}).ThenByDescendingT(func(item types.TokensResponse_Token) int64 {
				return item.Sort
			})
		case "1H_DESC":
			query = linq.From(dataResult).OrderByDescendingT(func(item types.TokensResponse_Token) float64 {
				value, _ := strconv.ParseFloat(item.RiseFall1H.Percentage, 64)
				return value
			}).ThenByDescendingT(func(item types.TokensResponse_Token) int64 {
				return item.Sort
			})
		case "1D_ASC":
			query = linq.From(dataResult).OrderByT(func(item types.TokensResponse_Token) float64 {
				value, _ := strconv.ParseFloat(item.RiseFall1D.Percentage, 64)
				return value
			}).ThenByDescendingT(func(item types.TokensResponse_Token) int64 {
				return item.Sort
			})
		case "1D_DESC":
			query = linq.From(dataResult).OrderByDescendingT(func(item types.TokensResponse_Token) float64 {
				value, _ := strconv.ParseFloat(item.RiseFall1D.Percentage, 64)
				return value
			}).ThenByDescendingT(func(item types.TokensResponse_Token) int64 {
				return item.Sort
			})
		case "Volume_ASC":
			query = linq.From(dataResult).OrderByT(func(item types.TokensResponse_Token) float64 {
				value, _ := strconv.ParseFloat(item.VolumeUSD, 64)
				return value
			}).ThenByDescendingT(func(item types.TokensResponse_Token) int64 {
				return item.Sort
			})
		case "Volume_DESC":
			query = linq.From(dataResult).OrderByDescendingT(func(item types.TokensResponse_Token) float64 {
				value, _ := strconv.ParseFloat(item.VolumeUSD, 64)
				return value
			}).ThenByDescendingT(func(item types.TokensResponse_Token) int64 {
				return item.Sort
			})
		default:
			query = linq.From(dataResult).OrderByDescendingT(func(item types.TokensResponse_Token) int64 {
				return item.Sort
			})
		}
	} else {
		query = linq.From(dataResult).OrderByDescendingT(func(item types.TokensResponse_Token) int64 {
			return item.Sort
		})
	}
	query.ThenByT(func(item types.TokensResponse_Token) string { return item.Symbol }).ToSlice(&dataResult)

	// Page
	totalResult := int64(len(dataResult))
	var meta types.PageResponse
	if req.PageSize != nil && req.PageIndex != nil {
		if *req.PageSize < 1 {
			*req.PageSize = 20
		}
		if *req.PageIndex < 1 {
			*req.PageIndex = 1
		}

		linq.From(dataResult).Skip(int((*req.PageIndex - 1) * *req.PageSize)).Take(int(*req.PageSize)).ToSlice(&dataResult)
		meta = types.PageResponse{
			PageSize:  *req.PageSize,
			PageIndex: *req.PageIndex,
		}
	} else {
		meta = types.PageResponse{
			PageSize:  0,
			PageIndex: 0,
		}
	}
	meta.Total = totalResult

	return &response.ApiResponse{
		Data: dataResult,
		Meta: meta,
	}, nil
}

func (l *TokensLogic) getData(req *types.TokensRequest) ([]types.TokensResponse_Token, error) {
	chainConf := l.svcCtx.Config.GetChainConfig(req.ChainId)
	if chainConf == nil {
		return nil, fmt.Errorf("%s", "Miss chain configure")
	}

	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	eventSonexFactoryPoolCreatedModel := repository.NewEventSonexFactoryPoolCreatedModel(conn)
	tokenModel := repository.NewTokenModel(conn)

	whitelistTokens, err := tokenModel.QueryByChainId(l.ctx, req.ChainId)
	if err != nil {
		return nil, err
	}
	whitelistTokenMap := map[string]repository.Token{}
	linq.From(whitelistTokens).ToMapByT(
		&whitelistTokenMap,
		func(item repository.Token) string { return item.Address },
		func(item repository.Token) repository.Token { return item },
	)

	dateResult := []types.TokensResponse_Token{}

	poolRecords, err := eventSonexFactoryPoolCreatedModel.QueryByChainId(l.ctx, req.ChainId)
	if err != nil {
		return nil, err
	}

	linq.From(poolRecords).GroupByT(
		func(item repository.EventSonexFactoryPoolCreated) string { return item.EventPool },
		func(item repository.EventSonexFactoryPoolCreated) repository.EventSonexFactoryPoolCreated {
			return item
		},
	)

	tokenWithPoolsMap := map[string]types.TokensResponse_Token{}
	tokenPoolsMap := map[string][]types.Pool{}
	var dataLock sync.Mutex
	var wg sync.WaitGroup
	results := make(chan error)

	for _, record := range poolRecords {
		wg.Add(1)
		go func(record repository.EventSonexFactoryPoolCreated, ch chan error) {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					ch <- fmt.Errorf("%s %s(panic)", record.EventPool, err)
				}
			}()

			pool := types.Pool{
				Address: record.EventPool,
				Fee:     record.EventFee,
			}
			tokenInfo0, err := global.GetToken(req.ChainId, record.EventToken0)
			if err != nil {
				ch <- err
			}
			tokenInfo1, err := global.GetToken(req.ChainId, record.EventToken1)
			if err != nil {
				ch <- err
			}

			dataLock.Lock()
			token0, ok0 := tokenWithPoolsMap[record.EventToken0]
			if !ok0 {
				token0 = types.TokensResponse_Token{
					TokenWithExtra: types.TokenWithExtra{
						Token: types.Token{
							Address:  tokenInfo0.Address,
							Name:     tokenInfo0.Name,
							Symbol:   tokenInfo0.Symbol,
							Decimals: tokenInfo0.Decimals,
						},
					},
				}
			}
			_, pOk0 := tokenPoolsMap[record.EventToken0]
			if !pOk0 {
				tokenPoolsMap[record.EventToken0] = []types.Pool{}
			}
			pool.Token0 = token0.Token

			token1, ok1 := tokenWithPoolsMap[record.EventToken1]
			if !ok1 {
				token1 = types.TokensResponse_Token{
					TokenWithExtra: types.TokenWithExtra{
						Token: types.Token{
							Address:  tokenInfo1.Address,
							Name:     tokenInfo1.Name,
							Symbol:   tokenInfo1.Symbol,
							Decimals: tokenInfo1.Decimals,
						},
					},
				}
			}
			_, pOk1 := tokenPoolsMap[record.EventToken1]
			if !pOk1 {
				tokenPoolsMap[record.EventToken1] = []types.Pool{}
			}
			pool.Token1 = token1.Token
			tokenWithPoolsMap[record.EventToken0] = token0
			tokenWithPoolsMap[record.EventToken1] = token1
			tokenPoolsMap[record.EventToken0] = append(tokenPoolsMap[record.EventToken0], pool)
			tokenPoolsMap[record.EventToken1] = append(tokenPoolsMap[record.EventToken1], pool)
			dataLock.Unlock()
		}(record, results)
	}
	wg.Wait()
	close(results)

	for result := range results {
		if result != nil {
			return nil, err
		}
	}

	for _, wt := range whitelistTokens {
		data, tok := tokenWithPoolsMap[wt.Address]
		pools, pok := tokenPoolsMap[wt.Address]
		if tok && pok {
			data.Pools = pools
		} else {
			data = types.TokensResponse_Token{
				TokenWithExtra: types.TokenWithExtra{
					Token: types.Token{
						Address:  wt.Address,
						Name:     wt.Name,
						Symbol:   wt.Symbol,
						Decimals: wt.Decimals,
					},
					Pools: []types.Pool{},
				},
			}
		}
		dateResult = append(dateResult, data)
	}

	return dateResult, nil
}

func (l *TokensLogic) comparePrices(currentPrice, comparePrice *big.Int) (string, *big.Float) {
	currentPriceF := new(big.Float).SetInt(currentPrice)
	comparePriceF := new(big.Float).SetInt(comparePrice)
	var comparison string
	var percentage *big.Float
	diff := new(big.Float).Sub(currentPriceF, comparePriceF)
	if diff.Cmp(big.NewFloat(0)) > 0 {
		comparison = "UP"
		percentage = new(big.Float).Quo(diff, comparePriceF)
	} else if diff.Cmp(big.NewFloat(0)) < 0 {
		comparison = "DOWN"
		percentage = new(big.Float).Quo(diff, comparePriceF)
	} else {
		comparison = "SAME"
		percentage = big.NewFloat(0)
	}
	return comparison, percentage
}
