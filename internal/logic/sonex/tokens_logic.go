package sonex

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
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

	// No cache
	if dataResult == nil {
		data, err := l.getData(req)
		if err != nil {
			return nil, err
		}
		dataResult = data

		// Set cache
		respBytes, err := json.Marshal(CacheWarpper{Data: data, LastUpdate: time.Now().UTC()})
		if err != nil {
			logx.Errorf("Marshal response error: %v", err)
		}
		err = redisCache.Set(cacheKey, respBytes, time.Hour*24)
		if err != nil {
			logx.Errorf("Set cache error: %v", err)
		}
	}
	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	tokenPriceModel := repository.NewSonexTokenPriceRecordModel(conn)
	eventSonexPoolSwapModel := repository.NewEventSonexPoolSwapModel(conn)

	// Update price
	prices, err := tokenPriceModel.QueryLatestHistoryPrices()
	if err != nil {
		return nil, err
	}
	volumes, err := eventSonexPoolSwapModel.SummarizeVolumesByChainId(l.ctx, req.ChainId)
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

			priceI, ok := prices[item.Address]
			if !ok {
				priceI = big.NewInt(0)
			}
			price, _ := new(big.Float).Quo(new(big.Float).SetInt(priceI), usdAccuracy).Float64()
			item.PriceUSD = price
			for i, pool := range item.Pools {
				price0I, ok := prices[pool.Token0.Address]
				if !ok {
					price0I = big.NewInt(0)
				}
				price0, _ := new(big.Float).Quo(new(big.Float).SetInt(price0I), usdAccuracy).Float64()
				item.Pools[i].Token0.PriceUSD = price0
				price1I, ok := prices[pool.Token1.Address]
				if !ok {
					price1I = big.NewInt(0)
				}
				price1, _ := new(big.Float).Quo(new(big.Float).SetInt(price1I), usdAccuracy).Float64()
				item.Pools[i].Token1.PriceUSD = price1
			}

			volumeI, ok := volumes[item.Address]
			if !ok {
				volumeI = big.NewInt(0)
			}
			volume, _ := global.PriceToUSD(req.ChainId, volumeI)
			item.VolumeUSD = volume

			kline1D, err := tokenPriceModel.Kline(item.Address, "1d")
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
			var currentTimestamp int64 = 0
			last1DPrice := big.NewInt(0)
			if len(kline1D) > 0 {
				currentPrice, _ = new(big.Int).SetString(kline1D[len(kline1D)-1].Price, 10)
				t, _ := time.Parse("2006-01-02 15:04:05", kline1D[len(kline1D)-1].CreatedAt)
				currentTimestamp = t.Unix()
				last1DPrice, _ = new(big.Int).SetString(kline1D[len(kline1D)-1].Price, 10)
			}
			last1HPrice := big.NewInt(0)
			for _, tick := range kline1D {
				priceI, _ := new(big.Int).SetString(tick.Price, 10)
				price, _ := new(big.Float).Quo(new(big.Float).SetInt(priceI), usdAccuracy).Float64()
				t, _ := time.Parse("2006-01-02 15:04:05", tick.CreatedAt)
				item.Kline = append(item.Kline, types.Kline{
					Price:     price,
					TimeStamp: t.Unix(),
				})

				if last1HPrice.Cmp(big.NewInt(0)) == 0 && (t.Unix()-currentTimestamp)/3600 > 1 {
					last1HPrice, _ = new(big.Int).SetString(tick.Price, 10)
				}
			}
			comparison1H, percentage1HF := l.comparePrices(currentPrice, last1HPrice)
			percentage1H, _ := new(big.Float).Quo(percentage1HF, usdAccuracy).Float64()
			comparison1D, percentage1DF := l.comparePrices(currentPrice, last1DPrice)
			percentage1D, _ := new(big.Float).Quo(percentage1DF, usdAccuracy).Float64()
			item.RiseFall1H = types.Comparison{Comparison: comparison1H, Percentage: percentage1H}
			item.RiseFall1D = types.Comparison{Comparison: comparison1D, Percentage: percentage1D}

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
	for result := range results {
		if result.Err != nil {
			return nil, result.Err
		}
		dataResult = append(dataResult, result.Data)
	}

	// Sorts
	if req.OrderBy != "" {
		switch req.OrderBy {
		case "1H_ASC":
			linq.From(dataResult).OrderByT(func(item types.TokensResponse_Token) float64 {
				return item.RiseFall1H.Percentage
			}).ToSlice(&dataResult)
		case "1H_DESC":
			linq.From(dataResult).OrderByDescendingT(func(item types.TokensResponse_Token) float64 {
				return item.RiseFall1H.Percentage
			}).ToSlice(&dataResult)
		case "1D_ASC":
			linq.From(dataResult).OrderByT(func(item types.TokensResponse_Token) float64 {
				return item.RiseFall1D.Percentage
			}).ToSlice(&dataResult)
		case "1D_DESC":
			linq.From(dataResult).OrderByDescendingT(func(item types.TokensResponse_Token) float64 {
				return item.RiseFall1D.Percentage
			}).ToSlice(&dataResult)
		case "Volume_ASC":
			linq.From(dataResult).OrderByT(func(item types.TokensResponse_Token) float64 {
				return item.VolumeUSD
			}).ToSlice(&dataResult)
		case "Volume_DESC":
			linq.From(dataResult).OrderByDescendingT(func(item types.TokensResponse_Token) float64 {
				return item.VolumeUSD
			}).ToSlice(&dataResult)
		}
	}

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

	tokenMap := map[string]types.TokensResponse_Token{}
	tokenPoolsMap := map[string][]types.Pool{}
	var dataLock sync.Mutex
	var wg sync.WaitGroup
	results := make(chan error)

	for _, record := range poolRecords {
		info0, ok0 := whitelistTokenMap[record.EventToken0]
		info1, ok1 := whitelistTokenMap[record.EventToken1]
		if !ok0 && !ok1 {
			continue
		}

		wg.Add(1)
		go func(record repository.EventSonexFactoryPoolCreated, info0 *repository.Token, info1 *repository.Token, ch chan error) {
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

			dataLock.Lock()
			token0, ok0 := tokenMap[record.EventToken0]
			_, pOk0 := tokenPoolsMap[record.EventToken0]
			if !pOk0 {
				tokenPoolsMap[record.EventToken0] = []types.Pool{}
			}
			token1, ok1 := tokenMap[record.EventToken1]
			_, pOk1 := tokenPoolsMap[record.EventToken1]
			if !pOk1 {
				tokenPoolsMap[record.EventToken1] = []types.Pool{}
			}
			dataLock.Unlock()

			if !ok0 {
				token, err := global.GetToken(req.ChainId, record.EventToken0)
				if err != nil {
					ch <- err
				}
				var logo string
				if info0 != nil {
					logo = info0.Logo
				}
				token0 = types.TokensResponse_Token{
					TokenWithExtra: types.TokenWithExtra{
						Token: types.Token{
							Address:  token.Address,
							Name:     token.Name,
							Symbol:   token.Symbol,
							Decimals: token.Decimals,
						},
					},
					Logo: logo,
				}
			}
			pool.Token0 = token0.Token

			if !ok1 {
				token, err := global.GetToken(req.ChainId, record.EventToken1)
				if err != nil {
					ch <- err
				}
				var logo string
				if info1 != nil {
					logo = info1.Logo
				}
				token1 = types.TokensResponse_Token{
					TokenWithExtra: types.TokenWithExtra{
						Token: types.Token{
							Address:  token.Address,
							Name:     token.Name,
							Symbol:   token.Symbol,
							Decimals: token.Decimals,
						},
					},
					Logo: logo,
				}
			}
			pool.Token1 = token1.Token

			dataLock.Lock()
			tokenMap[record.EventToken0] = token0
			tokenMap[record.EventToken1] = token1
			tokenPoolsMap[record.EventToken0] = append(tokenPoolsMap[record.EventToken0], pool)
			tokenPoolsMap[record.EventToken1] = append(tokenPoolsMap[record.EventToken1], pool)
			dataLock.Unlock()
		}(record, &info0, &info1, results)
	}
	wg.Wait()
	close(results)

	for result := range results {
		if result != nil {
			return nil, err
		}
	}

	for key, value := range tokenMap {
		if pools, ok := tokenPoolsMap[key]; ok {
			value.Pools = pools
		}
		dateResult = append(dateResult, value)
	}

	return dateResult, nil
}

func (l *TokensLogic) comparePrices(currentPrice, comparePrice *big.Int) (string, *big.Float) {
	currentPriceF := new(big.Float).SetInt(currentPrice)
	comparePriceF := new(big.Float).SetInt(comparePrice)
	var comparison string
	var diff *big.Float
	if currentPrice.Cmp(comparePrice) > 0 {
		comparison = "UP"
	} else if currentPrice.Cmp(comparePrice) < 0 {
		comparison = "DOWN"
	} else {
		comparison = "SAME"
	}
	diff = new(big.Float).Sub(currentPriceF, comparePriceF)
	if comparePriceF.Cmp(big.NewFloat(0)) == 0 {
		comparePriceF = big.NewFloat(1)
	}
	percentage := new(big.Float).Quo(diff, comparePriceF)
	return comparison, percentage
}
