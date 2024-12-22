package sonex

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TokensKlineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
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

	resp = new(types.SonexTokenKlineResponse)

	gconv.Structs(result, &resp.Data)

	return
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
