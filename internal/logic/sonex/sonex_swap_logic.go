package sonex

import (
	"context"
	"math"
	"math/big"

	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SonexSwapLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexSwapLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexSwapLogic {
	return &SonexSwapLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type SwapEventResponseWrapper struct {
	*types.SwapEventResponse
}

func (p *SwapEventResponseWrapper) GetData() interface{} {
	return p.Data
}

func (p *SwapEventResponseWrapper) GetMeta() interface{} {
	return p.Meta
}

func (p *SwapEventResponseWrapper) GetCode() interface{} {
	return p.Code
}

func (l *SonexSwapLogic) SonexSwap(req *types.SwapEventRequest) (resp *SwapEventResponseWrapper, err error) {
	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	model := repository.NewEventSonexPoolSwapModel(conn)
	data, meta, err := model.GetSwapEvents(l.ctx, req)
	if err != nil {
		return &SwapEventResponseWrapper{
			&types.SwapEventResponse{
				Code:    400,
				Data:    []*types.SwapEvent{},
				Meta:    &types.PaginationMeta{},
				Message: err.Error(),
			},
		}, nil
	}

	for _, item := range data {
		bestPrice, _, err := global.GetPrices(l.ctx, nil, req.ChainId, item.Token1.Address, true, 2, nil)
		if err == nil && bestPrice != nil {
			priceFloat := new(big.Float).SetInt(bestPrice.Value)
			amount1Positive := new(big.Float).SetFloat64(math.Abs(item.Amount1))
			decimalFactor := new(big.Float).SetFloat64(math.Pow10(int(item.Token1.Decimals)))
			usdValue := new(big.Float).Quo(new(big.Float).Mul(priceFloat, amount1Positive), decimalFactor)
			item.Usd, _ = usdValue.Float64()
			item.Amount1, _ = amount1Positive.Float64()
		}

	}

	return &SwapEventResponseWrapper{
		&types.SwapEventResponse{
			Code:    200,
			Data:    data,
			Meta:    meta,
			Message: "success",
		},
	}, nil
}
