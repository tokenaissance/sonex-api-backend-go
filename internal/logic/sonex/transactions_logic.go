package sonex

import (
	"context"
	"fmt"
	"math/big"

	"sonex_api_backend_go/internal/handler/response"
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

func (l *TransactionsLogic) Transactions(req *types.TransactionsRequest) (resp *response.ApiResponse, err error) {
	if req.ChainId <= 0 {
		return nil, fmt.Errorf("invalid parameter")
	}
	if req.PageSize == nil || *req.PageSize < 1 {
		var defaultSize int64 = 20
		req.PageSize = &defaultSize
	}
	if req.PageIndex == nil || *req.PageIndex < 1 {
		var defaultIndex int64 = 1
		req.PageIndex = &defaultIndex
	}

	chainConf := l.svcCtx.Config.GetChainConfig(req.ChainId)
	if chainConf == nil {
		return nil, fmt.Errorf("%s", "Miss chain configure")
	}

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
		swapRecords, totalResult, err = eventSonexPoolSwapModel.PageQueryOrderByTxAtWithExtra(l.ctx, req.OrderBy, *req.PageSize, *req.PageIndex)
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
			valueI, _ := new(big.Int).SetString(record.Value, 10)
			if valueI.Cmp(common.Big0) > 0 {
				if token0.Address == common.HexToAddress(chainConf.WETH).Hex() {
					token0.Address = (common.Address{}).Hex()
					token0.Name = chainConf.ETH.Name
					token0.Symbol = chainConf.ETH.Symbol
					token0.Decimals = chainConf.ETH.Decimals
				} else if token1.Address == common.HexToAddress(chainConf.WETH).Hex() {
					token1.Address = (common.Address{}).Hex()
					token1.Name = chainConf.ETH.Name
					token1.Symbol = chainConf.ETH.Symbol
					token1.Decimals = chainConf.ETH.Decimals
				}
			}

			prices, err := tokenPriceModel.GetCurrentPrices([]string{token0.Address, token1.Address})
			if err != nil {
				return nil, err
			}
			price0I := big.NewInt(0)
			price1I := big.NewInt(0)
			price, ok := prices[token0.Address]
			if ok {
				price0I, _ = new(big.Int).SetString(price.Price, 10)
			}
			price0, _ := global.PriceToUSD(req.ChainId, price0I)
			price, ok = prices[token1.Address]
			if ok {
				price1I, _ = new(big.Int).SetString(price.Price, 10)
			}
			price1, _ := global.PriceToUSD(req.ChainId, price1I)

			pool = types.Pool{
				Address: poolAddress.Hex(),
				Token0: types.Token{
					Address:  token0.Address,
					Name:     token0.Name,
					Symbol:   token0.Symbol,
					Decimals: token0.Decimals,
					PriceUSD: price0,
				},
				Token1: types.Token{
					Address:  token1.Address,
					Name:     token1.Name,
					Symbol:   token1.Symbol,
					Decimals: token1.Decimals,
					PriceUSD: price1,
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
		price0 := big.NewFloat(pool.Token0.PriceUSD)
		price1 := big.NewFloat(pool.Token1.PriceUSD)
		accuracy0 := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(pool.Token0.Decimals), nil))
		accuracy1 := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(pool.Token1.Decimals), nil))
		valueUSD0, _ := new(big.Float).Quo(new(big.Float).Mul(new(big.Float).SetInt(amount0), price0), accuracy0).Float64()
		valueUSD1, _ := new(big.Float).Quo(new(big.Float).Mul(new(big.Float).SetInt(amount1), price1), accuracy1).Float64()

		var valueUSD float64
		if valueUSD0 > valueUSD1 {
			from = record.ExtraToken0
			to = record.ExtraToken1
			valueUSD = valueUSD0
		} else {
			from = record.ExtraToken1
			to = record.ExtraToken0
			valueUSD = valueUSD1
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
		linq.From(dataResult).Skip(int((*req.PageIndex - 1) * *req.PageSize)).Take(int(*req.PageSize)).ToSlice(&dataResult)
	}

	return &response.ApiResponse{
		Data: dataResult,
		Meta: types.PageResponse{
			PageSize:  *req.PageSize,
			PageIndex: *req.PageIndex,
			Total:     totalResult,
		},
	}, nil
}
