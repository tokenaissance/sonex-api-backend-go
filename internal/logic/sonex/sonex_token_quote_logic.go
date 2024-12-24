package sonex

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/big"

	"sonex_api_backend_go/internal/lib/cache"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/lib/quoter"
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
	amount, _ := new(big.Int).SetString(req.Amount, 10)
	var maxLength int64 = 3

	var prices []cache.Price
	switch req.Type {
	case "INPUT":
		prices, err = global.QuoteExactInput(l.ctx, nil, req.ChainId, fromToken.Hex(), toToken.Hex(), amount, maxLength)
		if err != nil {
			return nil, err
		}
	case "OUTPUT":
		prices, err = global.QuoteExactOutput(l.ctx, nil, req.ChainId, fromToken.Hex(), toToken.Hex(), amount, maxLength)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("%s", "Invalid type")
	}

	var bestRouterResult *types.TokenQuoteResponse_Router
	maxValue := common.Big0
	for _, price := range prices {
		if maxValue.Cmp(price.Value) < 0 {
			maxValue = price.Value
			var router []string
			linq.From(price.Pools).SelectT(func(pool quoter.Pool) string { return pool.Address.Hex() }).ToSlice(&router)
			var sqrtPriceX96AfterList []string
			linq.From(price.SqrtPriceX96AfterList).SelectT(func(sqrtPriceX96After *big.Int) string { return sqrtPriceX96After.String() }).ToSlice(&sqrtPriceX96AfterList)
			var initializedTicksCrossedList []int64
			linq.From(price.InitializedTicksCrossedList).SelectT(func(initializedTicksCrossed uint32) int64 { return int64(initializedTicksCrossed) }).ToSlice(&initializedTicksCrossedList)
			bestRouterResult = &types.TokenQuoteResponse_Router{
				Target:                      price.To,
				Router:                      router,
				Path:                        base64.StdEncoding.EncodeToString(price.Path),
				Price:                       price.Value.String(),
				SqrtPriceX96AfterList:       sqrtPriceX96AfterList,
				InitializedTicksCrossedList: initializedTicksCrossedList,
				GasEstimate:                 price.GasEstimate.String(),
			}
		}
	}

	routerResult := []types.TokenQuoteResponse_Router{}
	linq.From(prices).
		SelectT(func(item cache.Price) types.TokenQuoteResponse_Router {
			var router []string
			linq.From(item.Pools).SelectT(func(pool quoter.Pool) string { return pool.Address.Hex() }).ToSlice(&router)
			var sqrtPriceX96AfterList []string
			linq.From(item.SqrtPriceX96AfterList).SelectT(func(sqrtPriceX96After *big.Int) string { return sqrtPriceX96After.String() }).ToSlice(&sqrtPriceX96AfterList)
			var initializedTicksCrossedList []int64
			linq.From(item.InitializedTicksCrossedList).SelectT(func(initializedTicksCrossed uint32) int64 { return int64(initializedTicksCrossed) }).ToSlice(&initializedTicksCrossedList)
			return types.TokenQuoteResponse_Router{
				Target:                      item.To,
				Router:                      router,
				Path:                        base64.StdEncoding.EncodeToString(item.Path),
				Price:                       item.Value.String(),
				SqrtPriceX96AfterList:       sqrtPriceX96AfterList,
				InitializedTicksCrossedList: initializedTicksCrossedList,
				GasEstimate:                 item.GasEstimate.String(),
			}
		}).
		ToSlice(&routerResult)

	return &types.TokenQuoteResponse{
		BestRouter: bestRouterResult,
		Routers:    routerResult,
	}, nil
}
