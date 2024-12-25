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

	bestPrice, prices, err := global.GetPrices(l.ctx, nil, req.ChainId, token.Hex(), req.IsCache, 3, req.WantStablecoin)
	if err != nil {
		return nil, err
	}

	var bestPriceResult *types.TokenPriceResponse_Price
	if bestPrice != nil {
		var router []string
		linq.From(bestPrice.Pools).SelectT(func(pool quoter.Pool) string { return pool.Address.Hex() }).ToSlice(&router)
		bestPriceResult = &types.TokenPriceResponse_Price{
			Stablecoin: bestPrice.To,
			Router:     router,
			Value:      bestPrice.Value.String(),
		}
	}

	pricesResult := []types.TokenPriceResponse_Price{}
	linq.From(prices).
		SelectT(func(item cache.Price) types.TokenPriceResponse_Price {
			var router []string
			linq.From(item.Pools).SelectT(func(pool quoter.Pool) string { return pool.Address.Hex() }).ToSlice(&router)
			return types.TokenPriceResponse_Price{
				Stablecoin: item.To,
				Router:     router,
				Value:      item.Value.String(),
			}
		}).
		ToSlice(&pricesResult)

	return &types.TokenPriceResponse{
		BasePrice: bestPriceResult,
		Prices:    pricesResult,
	}, nil
}
