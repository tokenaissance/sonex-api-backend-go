package sonex

import (
	"context"
	"fmt"

	"sonex_api_backend_go/internal/lib/cache"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/lib/quoter"
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

	prices, err := global.GetPrices(l.ctx, nil, req.ChainId, token.Hex(), req.IsCache)
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
