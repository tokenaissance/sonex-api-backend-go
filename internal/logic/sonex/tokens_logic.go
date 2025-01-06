package sonex

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"sync"
	"time"

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
	Data       []types.TokenWithExtra
	LastUpdate time.Time
}

func NewTokensLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TokensLogic {
	return &TokensLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TokensLogic) Tokens(req *types.TokensRequest) (resp *types.TokensResponse, err error) {
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

	// Check cache
	lastUpdate, err := redisCache.Get(cacheKey)
	if err == nil && lastUpdate != "" {
		var cacheWarpper CacheWarpper
		err := json.Unmarshal([]byte(lastUpdate), &cacheWarpper)
		if err == nil {
			resp = &types.TokensResponse{Tokens: cacheWarpper.Data}

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
	if resp == nil {
		data, err := l.getData(req)
		if err != nil {
			return nil, err
		}
		resp = &types.TokensResponse{Tokens: data}

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
		Data types.TokenWithExtra
		Err  error
	}, len(resp.Tokens))

	for _, item := range resp.Tokens {
		wg.Add(1)
		go func(item types.TokenWithExtra, ch chan struct {
			Data types.TokenWithExtra
			Err  error
		}) {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					ch <- struct {
						Data types.TokenWithExtra
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
					Data types.TokenWithExtra
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
				Data types.TokenWithExtra
				Err  error
			}{
				Data: item,
			}
		}(item, results)
	}
	wg.Wait()
	close(results)

	// Handle results
	resp.Tokens = []types.TokenWithExtra{}
	for result := range results {
		if result.Err != nil {
			return nil, result.Err
		}
		resp.Tokens = append(resp.Tokens, result.Data)
	}

	// Sorts
	if req.OrderBy != "" {
		switch req.OrderBy {
		case "1H_ASC":
			linq.From(resp.Tokens).OrderByT(func(item types.TokenWithExtra) float64 {
				return item.RiseFall1H.Percentage
			}).ToSlice(&resp.Tokens)
		case "1H_DESC":
			linq.From(resp.Tokens).OrderByDescendingT(func(item types.TokenWithExtra) float64 {
				return item.RiseFall1H.Percentage
			}).ToSlice(&resp.Tokens)
		case "1D_ASC":
			linq.From(resp.Tokens).OrderByT(func(item types.TokenWithExtra) float64 {
				return item.RiseFall1D.Percentage
			}).ToSlice(&resp.Tokens)
		case "1D_DESC":
			linq.From(resp.Tokens).OrderByDescendingT(func(item types.TokenWithExtra) float64 {
				return item.RiseFall1D.Percentage
			}).ToSlice(&resp.Tokens)
		case "Volume_ASC":
			linq.From(resp.Tokens).OrderByT(func(item types.TokenWithExtra) float64 {
				return item.VolumeUSD
			}).ToSlice(&resp.Tokens)
		case "Volume_DESC":
			linq.From(resp.Tokens).OrderByDescendingT(func(item types.TokenWithExtra) float64 {
				return item.VolumeUSD
			}).ToSlice(&resp.Tokens)
		}
	}

	// Page
	resp.Total = int64(len(resp.Tokens))
	if req.PageSize != nil && req.PageIndex != nil {
		if *req.PageSize < 1 {
			*req.PageSize = 20
		}
		if *req.PageIndex < 1 {
			*req.PageIndex = 1
		}

		linq.From(resp.Tokens).Skip(int((*req.PageIndex - 1) * *req.PageSize)).Take(int(*req.PageSize)).ToSlice(&resp.Tokens)
	}

	return resp, nil
}

func (l *TokensLogic) getData(req *types.TokensRequest) ([]types.TokenWithExtra, error) {
	chainConf := l.svcCtx.Config.GetChainConfig(req.ChainId)
	if chainConf == nil {
		return nil, fmt.Errorf("%s", "Miss chain configure")
	}

	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	eventSonexFactoryPoolCreatedModel := repository.NewEventSonexFactoryPoolCreatedModel(conn)
	sonexTokenModel := repository.NewSonexTokenModel(conn)

	logos, err := sonexTokenModel.QueryLogo(l.ctx)
	if err != nil {
		return nil, err
	}

	dateResult := []types.TokenWithExtra{}

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

	tokenMap := map[string]types.TokenWithExtra{}
	tokenPoolsMap := map[string][]types.Pool{}
	var dataLock sync.Mutex
	var wg sync.WaitGroup
	results := make(chan error, len(poolRecords))

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
				logo := logos[record.EventToken0]
				token0 = types.TokenWithExtra{
					Token: types.Token{
						Address:  token.Address,
						Name:     token.Name,
						Symbol:   token.Symbol,
						Decimals: token.Decimals,
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
				logo := logos[record.EventToken0]
				token1 = types.TokenWithExtra{
					Token: types.Token{
						Address:  token.Address,
						Name:     token.Name,
						Symbol:   token.Symbol,
						Decimals: token.Decimals,
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
		}(record, results)
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
