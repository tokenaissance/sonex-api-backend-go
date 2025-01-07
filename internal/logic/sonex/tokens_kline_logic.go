package sonex

import (
	"context"
	"math/big"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/zeromicro/go-zero/core/logx"
)

type TokensKlineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

type PriceComparisonResult struct {
	Comparison string  `json:"comparison"` // 比较结果：UP、DOWN、SAME
	Percentage float64 `json:"percentage"` // 涨跌百分比
}

func NewTokensKlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TokensKlineLogic {
	return &TokensKlineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TokensKlineLogic) TokensKline(req *types.SonexTokenKlineRequest) (resp *types.SonexTokenKlineResponse, err error) {
	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	model := repository.NewSonexTokenPriceRecordModel(conn)
	result, err := model.Kline(req.Address, req.Interval)
	if err != nil {
		logx.Errorf("[TokensKline]: Kline err %s", err)
		return nil, err
	}

	resp = &types.SonexTokenKlineResponse{}

	gconv.Structs(result, &resp.Records)

	for i, v := range resp.Records {
		priceI, _ := new(big.Int).SetString(resp.Records[i].Price, 10)
		priceUSD, _ := global.PriceToUSD(req.ChainId, priceI)
		resp.Records[i].PriceUSD = priceUSD
		resp.Records[i].TimeStamp = gtime.New(v.CreatedAt).Timestamp()
	}
	// todo: cache kline

	return resp, nil
}

// 根据 interval 获取查询条件的时间范围
func getIntervalDuration(interval string) string {
	switch interval {
	case "1d":
		return "1 day"
	case "1w":
		return "1 week"
	case "1m":
		return "1 month"
	case "1y":
		return "1 year"
	default:
		return "1 day" // 默认按1天查询
	}
}

func (l *TokensKlineLogic) ComparePriceHour(tokenAddress string) PriceComparisonResult {
	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	model := repository.NewSonexTokenPriceRecordModel(conn)
	currentPrice, _ := model.GetCurrentPrice(tokenAddress)
	lastHourPrice, _ := model.GetLastHourPrice(tokenAddress)
	lastHourComparison := comparePrices(currentPrice, lastHourPrice)

	return lastHourComparison
}
func (l *TokensKlineLogic) ComparePriceDay(tokenAddress string) PriceComparisonResult {
	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	model := repository.NewSonexTokenPriceRecordModel(conn)
	currentPrice, _ := model.GetCurrentPrice(tokenAddress)
	yesterdayPrice, _ := model.GetYesterdayPrice(tokenAddress)
	yesterdayComparison := comparePrices(currentPrice, yesterdayPrice)

	return yesterdayComparison
}

func comparePrices(currentPrice, comparePrice float64) PriceComparisonResult {
	if currentPrice > comparePrice {
		percentage := ((currentPrice - comparePrice) / comparePrice) * 100
		return PriceComparisonResult{
			Comparison: "UP",
			Percentage: percentage,
		}
	} else if currentPrice < comparePrice {
		percentage := ((comparePrice - currentPrice) / comparePrice) * 100
		return PriceComparisonResult{
			Comparison: "DOWN",
			Percentage: percentage,
		}
	} else {
		return PriceComparisonResult{
			Comparison: "SAME",
			Percentage: 0,
		}
	}
}
