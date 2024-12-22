package sonex

import (
	"context"
	"encoding/base64"
	"fmt"

	"sonex_api_backend_go/internal/lib/cache"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/lib/quoter"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/ahmetb/go-linq/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type SonexTokenQuoteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexTokenQuoteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexTokenQuoteLogic {
	return &SonexTokenQuoteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SonexTokenQuoteLogic) SonexTokenQuote(req *types.TokenQuoteRequest) (resp *types.TokenQuoteResponse, err error) {
	if req.ChainId == 0 ||
		req.From == "" || !common.IsHexAddress(req.From) ||
		req.To == "" || !common.IsHexAddress(req.To) {
		return nil, fmt.Errorf("%s", "Invalid parameter")
	}
	fromToken := common.HexToAddress(req.From)
	toToken := common.HexToAddress(req.To)

	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	eventSonexFactoryPoolCreatedModel := repository.NewEventSonexFactoryPoolCreatedModel(conn)
	records, err := eventSonexFactoryPoolCreatedModel.QueryByChainId(l.ctx, req.ChainId)
	if err != nil {
		return nil, err
	}
	pools := []quoter.Pool{}
	linq.From(records).
		SelectT(func(item repository.EventSonexFactoryPoolCreated) quoter.Pool {
			return quoter.Pool{
				Address: common.HexToAddress(item.EventPool),
				Token0:  common.HexToAddress(item.EventToken0),
				Token1:  common.HexToAddress(item.EventToken1),
				Fee:     item.EventFee,
			}
		}).
		ToSlice(&pools)

	prices, err := global.Quote(pools, req.ChainId, fromToken.Hex(), toToken.Hex(), req.IsCache)
	if err != nil {
		return nil, err
	}

	result := []types.TokenQuoteResponse_Price{}
	linq.From(prices).
		SelectT(func(item cache.Price) types.TokenQuoteResponse_Price {
			var routers []string
			linq.From(item.Pools).SelectT(func(pool quoter.Pool) string { return pool.Address.Hex() }).ToSlice(&routers)
			return types.TokenQuoteResponse_Price{
				Target: item.To,
				Router: routers,
				Path:   base64.StdEncoding.EncodeToString(item.Path),
				Value:  item.Value.String(),
			}
		}).
		ToSlice(&result)

	return &types.TokenQuoteResponse{
		Prices: result,
	}, nil
}
