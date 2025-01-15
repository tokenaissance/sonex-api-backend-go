package sonex

import (
	"context"
	"fmt"
	"math/big"
	"sync"

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
		req.To == "" || !common.IsHexAddress(req.To) ||
		req.Amount == "" {
		return nil, fmt.Errorf("%s", "Invalid parameter")
	}
	amount, ok := new(big.Int).SetString(req.Amount, 10)
	if !ok {
		return nil, fmt.Errorf("%s", "Invalid parameter")
	}

	chainConf := l.svcCtx.Config.GetChainConfig(req.ChainId)
	if chainConf == nil {
		return nil, fmt.Errorf("%s", "Miss chain configure")
	}

	fromAddress := common.HexToAddress(req.From)
	toAddress := common.HexToAddress(req.To)
	stablecoinAddress := common.HexToAddress(chainConf.USD)

	if amount.Cmp(big.NewInt(0)) <= 0 {
		return &types.TokenQuoteResponse{
			SwapRouter: l.svcCtx.Config.GetSonexConfig(req.ChainId).SwapRouter,
			BestRouter: &types.TokenQuoteResponse_Router{
				TokenIn:                     common.HexToAddress(req.From).Hex(),
				TokenOut:                    common.HexToAddress(req.To).Hex(),
				AmountIn:                    "0",
				AmountOut:                   "0",
				Stablecoin:                  stablecoinAddress.Hex(),
				PriceUSDIn:                  "0",
				PriceUSDOut:                 "0",
				Router:                      []types.TokenQuoteResponse_Router_Pool{},
				Path:                        "0x",
				SqrtPriceX96AfterList:       []string{},
				InitializedTicksCrossedList: []int64{},
				GasEstimate:                 "0",
				PriceImpact:                 "0%",
				FeeRate:                     "0%",
			},
			Routers: []types.TokenQuoteResponse_Router{},
		}, nil
	}

	var maxLength int64 = 2
	fromPriceUSD := big.NewInt(0)
	toPriceUSD := big.NewInt(0)
	var prices []cache.Price

	var wg sync.WaitGroup
	results := make(chan error, 2)

	wg.Add(1)
	go func(ch chan error) {
		defer wg.Done()
		defer func() {
			if err := recover(); err != nil {
				ch <- fmt.Errorf("%s(panic)", err)
			}
		}()

		conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
		tokenPriceModel := repository.NewSonexTokenPriceRecordModel(conn)
		prices, err := tokenPriceModel.QueryLatestHistoryPricesByAddresses([]string{fromAddress.Hex(), toAddress.Hex()})
		if err != nil {
			ch <- err
			return
		}
		if price, ok := prices[fromAddress.Hex()]; ok {
			fromPriceUSD = price
		}
		if price, ok := prices[toAddress.Hex()]; ok {
			toPriceUSD = price
		}
	}(results)

	wg.Add(1)
	go func(ch chan error) {
		defer wg.Done()
		defer func() {
			if err := recover(); err != nil {
				ch <- fmt.Errorf("%s(panic)", err)
			}
		}()

		switch req.Type {
		case "INPUT":
			prices, _, err = global.QuoteExactInput(l.ctx, nil, req.ChainId, fromAddress.Hex(), toAddress.Hex(), amount, maxLength)
			if err != nil {
				ch <- err
				return
			}
		case "OUTPUT":
			prices, _, err = global.QuoteExactOutput(l.ctx, nil, req.ChainId, fromAddress.Hex(), toAddress.Hex(), amount, maxLength)
			if err != nil {
				ch <- err
				return
			}
		default:
			err = fmt.Errorf("%s", "Invalid type")
			return
		}
	}(results)

	wg.Wait()
	close(results)

	for err := range results {
		if err != nil {
			return nil, err
		}
	}

	bestRouterA := linq.From(prices).
		WhereT(func(item cache.Price) bool {
			return item.Value.Cmp(common.Big0) > 0
		}).
		OrderByT(func(item cache.Price) float64 {
			return item.PriceImpact
		}).First()
	var bestRouter cache.Price
	if bestRouterA != nil {
		bestRouter = bestRouterA.(cache.Price)
	}
	var router []types.TokenQuoteResponse_Router_Pool
	linq.From(bestRouter.Pools).
		SelectT(func(pool quoter.Pool) types.TokenQuoteResponse_Router_Pool {
			return types.TokenQuoteResponse_Router_Pool{
				Address: pool.Address.Hex(),
				Token0:  pool.Token0.Hex(),
				Token1:  pool.Token1.Hex(),
				Fee:     pool.Fee,
			}
		}).
		ToSlice(&router)
	var sqrtPriceX96AfterList []string
	linq.From(bestRouter.SqrtPriceX96AfterList).SelectT(func(sqrtPriceX96After *big.Int) string { return sqrtPriceX96After.String() }).ToSlice(&sqrtPriceX96AfterList)
	var initializedTicksCrossedList []int64
	linq.From(bestRouter.InitializedTicksCrossedList).SelectT(func(initializedTicksCrossed uint32) int64 { return int64(initializedTicksCrossed) }).ToSlice(&initializedTicksCrossedList)
	var amountIn *big.Int
	var amountOut *big.Int
	switch req.Type {
	case "INPUT":
		amountIn, amountOut = amount, bestRouter.Value
	case "OUTPUT":
		amountIn, amountOut = bestRouter.Value, amount
	}
	bestRouterResult := &types.TokenQuoteResponse_Router{
		TokenIn:                     fromAddress.Hex(),
		TokenOut:                    toAddress.Hex(),
		AmountIn:                    amountIn.String(),
		AmountOut:                   amountOut.String(),
		Stablecoin:                  stablecoinAddress.Hex(),
		PriceUSDIn:                  fromPriceUSD.String(),
		PriceUSDOut:                 toPriceUSD.String(),
		Router:                      router,
		Path:                        fmt.Sprintf("0x%s", common.Bytes2Hex(bestRouter.Path)),
		SqrtPriceX96AfterList:       sqrtPriceX96AfterList,
		InitializedTicksCrossedList: initializedTicksCrossedList,
		GasEstimate:                 bestRouter.GasEstimate.String(),
		PriceImpact:                 fmt.Sprintf("%.2f%%", bestRouter.PriceImpact*100),
		FeeRate:                     fmt.Sprintf("%.2f%%", bestRouter.FeeRate*100),
	}

	routerResult := []types.TokenQuoteResponse_Router{}
	linq.From(prices).
		SelectT(func(item cache.Price) types.TokenQuoteResponse_Router {
			var router []types.TokenQuoteResponse_Router_Pool
			linq.From(item.Pools).
				SelectT(func(pool quoter.Pool) types.TokenQuoteResponse_Router_Pool {
					return types.TokenQuoteResponse_Router_Pool{
						Address: pool.Address.Hex(),
						Token0:  pool.Token0.Hex(),
						Token1:  pool.Token1.Hex(),
						Fee:     pool.Fee,
					}
				}).
				ToSlice(&router)
			var sqrtPriceX96AfterList []string
			linq.From(item.SqrtPriceX96AfterList).SelectT(func(sqrtPriceX96After *big.Int) string { return sqrtPriceX96After.String() }).ToSlice(&sqrtPriceX96AfterList)
			var initializedTicksCrossedList []int64
			linq.From(item.InitializedTicksCrossedList).SelectT(func(initializedTicksCrossed uint32) int64 { return int64(initializedTicksCrossed) }).ToSlice(&initializedTicksCrossedList)
			var amountIn *big.Int
			var amountOut *big.Int
			switch req.Type {
			case "INPUT":
				amountIn, amountOut = amount, item.Value
			case "OUTPUT":
				amountIn, amountOut = item.Value, amount
			}
			return types.TokenQuoteResponse_Router{
				TokenIn:                     fromAddress.Hex(),
				TokenOut:                    toAddress.Hex(),
				AmountIn:                    amountIn.String(),
				AmountOut:                   amountOut.String(),
				Stablecoin:                  stablecoinAddress.Hex(),
				PriceUSDIn:                  fromPriceUSD.String(),
				PriceUSDOut:                 toPriceUSD.String(),
				Router:                      router,
				Path:                        fmt.Sprintf("0x%s", common.Bytes2Hex(item.Path)),
				SqrtPriceX96AfterList:       sqrtPriceX96AfterList,
				InitializedTicksCrossedList: initializedTicksCrossedList,
				GasEstimate:                 item.GasEstimate.String(),
				PriceImpact:                 fmt.Sprintf("%f%%", item.PriceImpact*100),
				FeeRate:                     fmt.Sprintf("%.2f%%", item.FeeRate*100),
			}
		}).
		ToSlice(&routerResult)

	return &types.TokenQuoteResponse{
		SwapRouter: l.svcCtx.Config.GetSonexConfig(req.ChainId).SwapRouter,
		BestRouter: bestRouterResult,
		Routers:    routerResult,
	}, nil
}
