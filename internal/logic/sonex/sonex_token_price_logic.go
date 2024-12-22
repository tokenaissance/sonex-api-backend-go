package sonex

import (
	"context"
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

type SonexTokenPriceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexTokenPriceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexTokenPriceLogic {
	return &SonexTokenPriceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SonexTokenPriceLogic) SonexTokenPrice(req *types.TokenPriceRequest) (resp *types.TokenPriceResponse, err error) {
	if req.ChainId == 0 || req.Token == "" || !common.IsHexAddress(req.Token) {
		return nil, fmt.Errorf("%s", "Invalid parameter")
	}
	token := common.HexToAddress(req.Token)

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

	prices, err := global.GetPrices(pools, req.ChainId, token.Hex(), req.IsCache)
	if err != nil {
		return nil, err
	}

	result := []types.TokenPriceResponse_Price{}
	linq.From(prices).
		SelectT(func(item cache.Price) types.TokenPriceResponse_Price {
			var routers []string
			linq.From(item.Pools).SelectT(func(pool quoter.Pool) string { return pool.Address.Hex() }).ToSlice(&routers)
			return types.TokenPriceResponse_Price{
				Stablecoin: item.To,
				Router:     routers,
				Value:      item.Value.String(),
			}
		}).
		ToSlice(&result)

	return &types.TokenPriceResponse{
		Prices: result,
	}, nil
}
