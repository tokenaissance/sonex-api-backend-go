package sonex

import (
	"context"
	"fmt"
	"math/big"

	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/ahmetb/go-linq/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type TransactionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTransactionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransactionsLogic {
	return &TransactionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TransactionsLogic) Transactions(req *types.TransactionsRequest) (resp *types.TransactionsResponse, err error) {
	if req.ChainId <= 0 {
		return nil, fmt.Errorf("invalid parameter")
	}
	if req.PageSize < 1 {
		req.PageSize = 20
	}
	if req.PageIndex < 1 {
		req.PageIndex = 1
	}

	chainConf := l.svcCtx.Config.GetChainConfig(req.ChainId)
	if chainConf == nil {
		return nil, fmt.Errorf("%s", "Miss chain configure")
	}
	usdToken, err := global.GetToken(req.ChainId, common.HexToAddress(chainConf.USD).Hex())
	if err != nil {
		return nil, err
	}
	usdAccuracy := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(usdToken.Decimals), nil))

	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	eventSonexPoolSwapModel := repository.NewEventSonexPoolSwapModel(conn)
	tokenPriceModel := repository.NewSonexTokenPriceRecordModel(conn)

	var swapRecords []repository.EventSonexPoolSwapWithExtra
	dataResult := []types.TransactionsResponse_Transaction{}
	var totalResult int64 = -1

	// Filter
	if req.Wallet != "" && common.IsHexAddress(req.Wallet) {
		walletAddress := common.HexToAddress(req.Wallet)
		if swapRecords == nil {
			swapRecords, err = eventSonexPoolSwapModel.QueryByWalletWithExtra(l.ctx, walletAddress.Hex())
			if err != nil {
				return nil, err
			}
			if swapRecords == nil {
				swapRecords = []repository.EventSonexPoolSwapWithExtra{}
			}
		} else {
			linq.From(swapRecords).
				WhereT(func(item repository.EventSonexPoolSwapWithExtra) bool { return item.From == walletAddress.Hex() }).
				ToSlice(&swapRecords)
		}
	}
	if req.Pool != "" && common.IsHexAddress(req.Pool) {
		poolAddress := common.HexToAddress(req.Pool)
		if swapRecords == nil {
			swapRecords, err = eventSonexPoolSwapModel.QueryByPoolWithExtra(l.ctx, poolAddress.Hex())
			if err != nil {
				return nil, err
			}
			if swapRecords == nil {
				swapRecords = []repository.EventSonexPoolSwapWithExtra{}
			}
		} else {
			linq.From(swapRecords).
				WhereT(func(item repository.EventSonexPoolSwapWithExtra) bool { return item.ExtraPool == poolAddress.Hex() }).
				ToSlice(&swapRecords)
		}
	}
	if req.Token0 != "" && common.IsHexAddress(req.Token0) {
		tokens := []string{common.HexToAddress(req.Token0).Hex()}
		if swapRecords == nil {
			swapRecords, err = eventSonexPoolSwapModel.QueryByTokensWithExtra(l.ctx, tokens)
			if err != nil {
				return nil, err
			}
			if swapRecords == nil {
				swapRecords = []repository.EventSonexPoolSwapWithExtra{}
			}
		} else {
			linq.From(swapRecords).
				WhereT(func(item repository.EventSonexPoolSwapWithExtra) bool {
					return linq.From(tokens).Contains(item.ExtraToken0) || linq.From(tokens).Contains(item.ExtraToken1)
				}).
				ToSlice(&swapRecords)
		}
	}
	if req.Token1 != "" && common.IsHexAddress(req.Token1) {
		tokens := []string{common.HexToAddress(req.Token1).Hex()}
		if swapRecords == nil {
			swapRecords, err = eventSonexPoolSwapModel.QueryByTokensWithExtra(l.ctx, tokens)
			if err != nil {
				return nil, err
			}
			if swapRecords == nil {
				swapRecords = []repository.EventSonexPoolSwapWithExtra{}
			}
		} else {
			linq.From(swapRecords).
				WhereT(func(item repository.EventSonexPoolSwapWithExtra) bool {
					return linq.From(tokens).Contains(item.ExtraToken0) || linq.From(tokens).Contains(item.ExtraToken1)
				}).
				ToSlice(&swapRecords)
		}
	}

	// Sort
	if swapRecords != nil && req.OrderBy != "" {
		switch req.OrderBy {
		case "ASC":
			linq.From(swapRecords).OrderByT(func(item repository.EventSonexPoolSwapWithExtra) int64 {
				return item.TxAt.Unix()
			}).ToSlice(&swapRecords)
		case "DESC":
			linq.From(swapRecords).OrderByDescendingT(func(item repository.EventSonexPoolSwapWithExtra) int64 {
				return item.TxAt.Unix()
			}).ToSlice(&swapRecords)
		}
	}

	// Default
	if swapRecords == nil {
		swapRecords, totalResult, err = eventSonexPoolSwapModel.PageQueryOrderByTxAtWithExtra(l.ctx, req.OrderBy, req.PageSize, req.PageIndex)
		if err != nil {
			return nil, err
		}
	}

	// Handle data
	pools := map[string]types.Pool{}
	for _, record := range swapRecords {
		poolAddress := common.HexToAddress(record.ExtraPool)
		pool, ok := pools[poolAddress.Hex()]
		if !ok {
			token0, err := global.GetToken(req.ChainId, record.ExtraToken0)
			if err != nil {
				return nil, err
			}
			token1, err := global.GetToken(req.ChainId, record.ExtraToken1)
			if err != nil {
				return nil, err
			}
			prices, err := tokenPriceModel.GetCurrentPrices([]string{token0.Address, token1.Address})
			if err != nil {
				return nil, err
			}
			price0 := big.NewInt(0)
			price1 := big.NewInt(0)
			price, ok := prices[token0.Address]
			if ok {
				price0, _ = new(big.Int).SetString(price.Price, 10)
			}
			price, ok = prices[token1.Address]
			if ok {
				price1, _ = new(big.Int).SetString(price.Price, 10)
			}

			pool = types.Pool{
				Address: poolAddress.Hex(),
				Token0: types.Token{
					Address:  token0.Address,
					Name:     token0.Name,
					Symbol:   token0.Symbol,
					Decimals: token0.Decimals,
					Price:    price0.String(),
				},
				Token1: types.Token{
					Address:  token1.Address,
					Name:     token1.Name,
					Symbol:   token1.Symbol,
					Decimals: token1.Decimals,
					Price:    price1.String(),
				},
				Fee: record.ExtraFee,
			}
			pools[poolAddress.Hex()] = pool
		}

		var from string
		var to string
		var fromAmount *big.Int
		var toAmount *big.Int
		amount0, _ := new(big.Int).SetString(record.EventAmount0, 10)
		amount1, _ := new(big.Int).SetString(record.EventAmount1, 10)
		price0, _ := new(big.Float).SetString(pool.Token0.Price)
		price1, _ := new(big.Float).SetString(pool.Token1.Price)
		var valueUSD float64

		if amount0.Cmp(amount1) < 0 {
			from = record.ExtraToken0
			to = record.ExtraToken1
			fromAmount = new(big.Int).Abs(amount0)
			toAmount = new(big.Int).Abs(amount1)
			fromAccuracy := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(pool.Token0.Decimals), nil))
			toAccuracy := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(pool.Token1.Decimals), nil))
			fromValueUSD := new(big.Float).Quo(new(big.Float).Quo(new(big.Float).Mul(new(big.Float).SetInt(fromAmount), price0), fromAccuracy), usdAccuracy)
			toValueUSD := new(big.Float).Quo(new(big.Float).Quo(new(big.Float).Mul(new(big.Float).SetInt(toAmount), price1), toAccuracy), usdAccuracy)
			valueUSD, _ = new(big.Float).Add(fromValueUSD, toValueUSD).Float64()
		} else {
			from = record.ExtraToken1
			to = record.ExtraToken0
			fromAmount = new(big.Int).Abs(amount1)
			toAmount = new(big.Int).Abs(amount0)
			fromAccuracy := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(pool.Token1.Decimals), nil))
			toAccuracy := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(pool.Token0.Decimals), nil))
			fromValueUSD := new(big.Float).Quo(new(big.Float).Quo(new(big.Float).Mul(new(big.Float).SetInt(fromAmount), price1), fromAccuracy), usdAccuracy)
			toValueUSD := new(big.Float).Quo(new(big.Float).Quo(new(big.Float).Mul(new(big.Float).SetInt(toAmount), price0), toAccuracy), usdAccuracy)
			valueUSD, _ = new(big.Float).Add(fromValueUSD, toValueUSD).Float64()
		}

		dataResult = append(dataResult, types.TransactionsResponse_Transaction{
			Wallet:     record.From,
			Pool:       pool,
			TxHash:     record.TxHash,
			TxAt:       record.TxAt.Unix(),
			From:       from,
			To:         to,
			FromAmount: fromAmount.String(),
			ToAmount:   toAmount.String(),
			ValueUSD:   valueUSD,
		})
	}

	// Page
	if totalResult < 0 {
		totalResult = int64(len(dataResult))
		linq.From(dataResult).Skip(int((req.PageIndex - 1) * req.PageSize)).Take(int(req.PageSize)).ToSlice(&dataResult)
	}

	return &types.TransactionsResponse{
		Data:  dataResult,
		Total: totalResult,
	}, nil
}
