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

	var swapRecords []repository.EventSonexMergePoolSwapWithExtra
	dataResult := []types.TransactionsResponse_Transaction{}
	var totalResult int64 = -1

	// Filter
	if req.Wallet != "" && common.IsHexAddress(req.Wallet) {
		walletAddress := common.HexToAddress(req.Wallet)
		if swapRecords == nil {
			swapRecords, err = eventSonexPoolSwapModel.QueryMergeByWalletWithExtra(l.ctx, walletAddress.Hex())
			if err != nil {
				return nil, err
			}
			if swapRecords == nil {
				swapRecords = []repository.EventSonexMergePoolSwapWithExtra{}
			}
		} else {
			linq.From(swapRecords).
				WhereT(func(item repository.EventSonexMergePoolSwapWithExtra) bool { return item.From == walletAddress.Hex() }).
				ToSlice(&swapRecords)
		}
	}
	if req.Pool != "" && common.IsHexAddress(req.Pool) {
		poolAddress := common.HexToAddress(req.Pool)
		if swapRecords == nil {
			swapRecords, err = eventSonexPoolSwapModel.QueryMergeByPoolWithExtra(l.ctx, poolAddress.Hex())
			if err != nil {
				return nil, err
			}
			if swapRecords == nil {
				swapRecords = []repository.EventSonexMergePoolSwapWithExtra{}
			}
		} else {
			linq.From(swapRecords).
				WhereT(func(item repository.EventSonexMergePoolSwapWithExtra) bool {
					return item.ExtraPool_0.String == poolAddress.Hex() || item.ExtraPool_1 == poolAddress.Hex()
				}).
				ToSlice(&swapRecords)
		}
	}
	if req.Token0 != "" && common.IsHexAddress(req.Token0) {
		tokens := []string{common.HexToAddress(req.Token0).Hex()}
		if swapRecords == nil {
			swapRecords, err = eventSonexPoolSwapModel.QueryMergeByTokensWithExtra(l.ctx, tokens)
			if err != nil {
				return nil, err
			}
			if swapRecords == nil {
				swapRecords = []repository.EventSonexMergePoolSwapWithExtra{}
			}
		} else {
			linq.From(swapRecords).
				WhereT(func(item repository.EventSonexMergePoolSwapWithExtra) bool {
					return linq.From(tokens).Contains(item.ExtraToken0) || linq.From(tokens).Contains(item.ExtraToken1)
				}).
				ToSlice(&swapRecords)
		}
	}
	if req.Token1 != "" && common.IsHexAddress(req.Token1) {
		tokens := []string{common.HexToAddress(req.Token1).Hex()}
		if swapRecords == nil {
			swapRecords, err = eventSonexPoolSwapModel.QueryMergeByTokensWithExtra(l.ctx, tokens)
			if err != nil {
				return nil, err
			}
			if swapRecords == nil {
				swapRecords = []repository.EventSonexMergePoolSwapWithExtra{}
			}
		} else {
			linq.From(swapRecords).
				WhereT(func(item repository.EventSonexMergePoolSwapWithExtra) bool {
					return linq.From(tokens).Contains(item.ExtraToken0) || linq.From(tokens).Contains(item.ExtraToken1)
				}).
				ToSlice(&swapRecords)
		}
	}

	// Sort
	if swapRecords != nil {
		switch req.OrderBy {
		case "ASC":
			linq.From(swapRecords).OrderByT(func(item repository.EventSonexMergePoolSwapWithExtra) int64 {
				return item.TxAt.Unix()
			}).ToSlice(&swapRecords)
		case "DESC":
			fallthrough
		default:
			linq.From(swapRecords).OrderByDescendingT(func(item repository.EventSonexMergePoolSwapWithExtra) int64 {
				return item.TxAt.Unix()
			}).ToSlice(&swapRecords)
		}
	}

	// Default
	if swapRecords == nil {
		swapRecords, totalResult, err = eventSonexPoolSwapModel.PageQueryMergeOrderByTxAtWithExtra(l.ctx, req.OrderBy, *req.PageSize, *req.PageIndex)
		if err != nil {
			return nil, err
		}
	}

	// Handle data
	poolMap := map[string]types.Pool{}
	for _, record := range swapRecords {
		handlePool := func(poolAddress common.Address) (*types.Pool, *types.Token, *types.Token, *big.Int, *big.Int, float64, error) {
			pool, ok := poolMap[poolAddress.Hex()]
			if !ok {
				token0, err := global.GetToken(req.ChainId, record.ExtraToken0())
				if err != nil {
					return nil, nil, nil, nil, nil, 0, err
				}
				token1, err := global.GetToken(req.ChainId, record.ExtraToken1())
				if err != nil {
					return nil, nil, nil, nil, nil, 0, err
				}

				prices, err := tokenPriceModel.GetCurrentPrices([]string{token0.Address, token1.Address})
				if err != nil {
					return nil, nil, nil, nil, nil, 0, err
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
						PriceUSD: fmt.Sprintf("%f", price0),
					},
					Token1: types.Token{
						Address:  token1.Address,
						Name:     token1.Name,
						Symbol:   token1.Symbol,
						Decimals: token1.Decimals,
						PriceUSD: fmt.Sprintf("%f", price1),
					},
					Fee: record.ExtraFee(),
				}
			}

			var from types.Token
			var to types.Token
			var fromAmount *big.Int
			var toAmount *big.Int
			amount0, _ := new(big.Int).SetString(record.EventAmount0(), 10)
			amount1, _ := new(big.Int).SetString(record.EventAmount1(), 10)
			price0, _ := new(big.Float).SetString(pool.Token0.PriceUSD)
			price1, _ := new(big.Float).SetString(pool.Token1.PriceUSD)
			accuracy0 := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(pool.Token0.Decimals), nil))
			accuracy1 := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(pool.Token1.Decimals), nil))
			valueUSD0, _ := new(big.Float).Quo(new(big.Float).Mul(new(big.Float).SetInt(amount0), price0), accuracy0).Float64()
			valueUSD1, _ := new(big.Float).Quo(new(big.Float).Mul(new(big.Float).SetInt(amount1), price1), accuracy1).Float64()

			var valueUSD float64
			if valueUSD0 > valueUSD1 {
				from = pool.Token0
				to = pool.Token1
				fromAmount = new(big.Int).Abs(amount0)
				toAmount = new(big.Int).Abs(amount1)
				valueUSD = valueUSD0
			} else {
				from = pool.Token1
				to = pool.Token0
				fromAmount = new(big.Int).Abs(amount1)
				toAmount = new(big.Int).Abs(amount0)
				valueUSD = valueUSD1
			}
			return &pool, &from, &to, fromAmount, toAmount, valueUSD, nil
		}
		pools := []types.Pool{}
		var from string
		var to string
		var fromAmount *big.Int
		var toAmount *big.Int
		var valueUSD float64
		if record.ExtraPool_0.Valid {
			pool, from_, _, fromAmount_, _, valueUSD_, err := handlePool(common.HexToAddress(record.ExtraPool_0.String))
			if err != nil {
				return nil, err
			}
			from = from_.Address
			fromAmount = fromAmount_
			valueUSD = valueUSD_
			poolMap[record.ExtraPool_0.String] = *pool
			pools = append(pools, *pool)
		}
		{
			pool, from_, to_, fromAmount_, toAmount_, valueUSD_, err := handlePool(common.HexToAddress(record.ExtraPool_1))
			if err != nil {
				return nil, err
			}
			if from == "" {
				from = from_.Address
				fromAmount = fromAmount_
				valueUSD = valueUSD_
			}
			to = to_.Address
			toAmount = toAmount_
			poolMap[record.ExtraPool_1] = *pool
			pools = append(pools, *pool)
		}

		// ETH
		valueI, _ := new(big.Int).SetString(record.Value, 10)
		if valueI.Cmp(common.Big0) > 0 {
			from = (common.Address{}).Hex()
		}

		dataResult = append(dataResult, types.TransactionsResponse_Transaction{
			Wallet:     record.From,
			Pools:      pools,
			TxHash:     record.TxHash,
			TxAt:       record.TxAt.Unix(),
			From:       from,
			To:         to,
			FromAmount: fromAmount.String(),
			ToAmount:   toAmount.String(),
			ValueUSD:   fmt.Sprintf("%f", valueUSD),
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
