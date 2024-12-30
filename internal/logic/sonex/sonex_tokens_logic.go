package sonex

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"sync"
	"time"

	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/zeromicro/go-zero/core/logx"
)

type SonexTokensLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexTokensLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexTokensLogic {
	return &SonexTokensLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type TokensResponseWrapper struct {
	*types.TokensResponse
}

func (p *TokensResponseWrapper) GetData() interface{} {
	return p.Data
}

func (p *TokensResponseWrapper) GetMeta() interface{} {
	return p.Meta
}

func (p *TokensResponseWrapper) GetCode() interface{} {
	return p.Code
}

func (l *SonexTokensLogic) SonexTokens(req *types.TokenListRequest) (resp *TokensResponseWrapper, err error) {
	const cacheKey = "sonex_dex:uniswap:tokens"

	cachedResponse, err := l.svcCtx.Cache.Get(cacheKey)
	if err == nil && len(cachedResponse) > 0 {
		var cachedResp TokensResponseWrapper
		if err := json.Unmarshal([]byte(cachedResponse), &cachedResp); err == nil {
			// Launch a goroutine to refresh the cache in the background
			go l.refreshCache(req, cacheKey, &cachedResp)
			// Return the cached response immediately
			return &cachedResp, nil
		}
	}

	return l.refreshCache(req, cacheKey, nil)
}

func (l *SonexTokensLogic) refreshCache(req *types.TokenListRequest, cacheKey string, cachedResp *TokensResponseWrapper) (*TokensResponseWrapper, error) {
	lockKey := fmt.Sprintf("lock:%s", cacheKey)
	ctx := context.Background()

	var lockAcquired bool

	if cachedResp != nil {
		success, err := l.svcCtx.Cache.SetNX(lockKey, "locked", 30*time.Second)
		if err != nil {
			return cachedResp, nil
		}

		if !success {
			return cachedResp, nil
		}

		lockAcquired = true
	}

	defer func() {
		if lockAcquired {
			l.svcCtx.Cache.Del(lockKey)
		}
	}()

	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	model := repository.NewSonexTokenModel(conn)

	tokenPriceModel := repository.NewSonexTokenPriceRecordModel(conn)

	data, meta, err := model.GetTokens(ctx, req, l.svcCtx.Config.Sonexes[0].PositionManager)

	if err != nil {
		return &TokensResponseWrapper{
			&types.TokensResponse{
				Code:    400,
				Data:    []*types.Token_TokenDetail{},
				Meta:    &types.PaginationMeta{},
				Message: err.Error(),
			},
		}, nil
	}

	wethAddress := l.svcCtx.Config.Chains[0].WETH

	tokenVolumes, err := model.GetTokenVolumes(ctx)
	if err == nil {
		var wg sync.WaitGroup
		var mu sync.Mutex

		for _, token := range data {
			wg.Add(1)
			go func(token *types.Token_TokenDetail) { // Replace TokenType with the actual type of token
				defer wg.Done()

				address := token.Address
				if token.Symbol == "ETH" {
					address = wethAddress
				}

				currentPrice, _ := tokenPriceModel.GetCurrentPrice(address)
				lastHourPrice, _ := tokenPriceModel.GetLastHourPrice(address)
				lastDayPrice, _ := tokenPriceModel.GetYesterdayPrice(address)
				lastHourComparison := comparePrices(currentPrice, lastHourPrice)
				lastDayComparison := comparePrices(currentPrice, lastDayPrice)
				kline, err := tokenPriceModel.Kline(address, "1d")
				if err == nil && len(kline) > 0 {
					token.Kline = make([]*types.Token_Kline, 0)
					for _, record := range kline {
						t, _ := gtime.StrToTime(record.CreatedAt)
						token.Kline = append(token.Kline, &types.Token_Kline{
							Price:     record.Price,
							TimeStamp: t.Timestamp(),
						})
					}
				}

				token.HourPriceComparison.Comparison = lastHourComparison.Comparison
				token.HourPriceComparison.Percentage = lastHourComparison.Percentage
				token.DayPriceComparison.Comparison = lastDayComparison.Comparison
				token.DayPriceComparison.Percentage = lastDayComparison.Percentage

				factor := new(big.Float).SetFloat64(math.Pow10(int(token.Decimals)))
				price := new(big.Float).SetFloat64(0)
				factorDai := new(big.Float).SetFloat64(math.Pow10(18))

				bestPrice, _, err := global.GetPrices(ctx, nil, req.ChainId, address, true, 2, &l.svcCtx.Config.Chains[0].USD)

				if err == nil && bestPrice != nil {
					price.SetInt(bestPrice.Value)
				}

				tokenPrice := new(big.Float).Quo(price, factorDai)
				mu.Lock()
				token.Price = tokenPrice.Text('f', -1)
				mu.Unlock()

				for _, volume := range tokenVolumes {
					if address == volume.Address {
						volumeAmounts := new(big.Float).SetFloat64(volume.Amounts)
						tokenVolume := new(big.Float).Quo(volumeAmounts, factor)
						mu.Lock()
						token.Volume = tokenVolume.Text('f', -1)
						fdv := new(big.Float).Mul(tokenPrice, tokenVolume)
						token.Fdv = fdv.Text('f', -1)
						mu.Unlock()
						break
					}
				}
			}(token)
		}
		wg.Wait()
	}

	response := &TokensResponseWrapper{
		&types.TokensResponse{
			Code:    200,
			Data:    data,
			Meta:    meta,
			Message: "success",
		},
	}

	responseBytes, err := json.Marshal(response)
	if err == nil {
		l.svcCtx.Cache.Set(cacheKey, responseBytes, 0)
	}

	return response, nil
}
