package sonex

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"sonex_api_backend_go/internal/handler/response"
	"sonex_api_backend_go/internal/lib/cache"
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
	st := time.Now()

	if req.ChainId <= 0 {
		return nil, fmt.Errorf("invalid parameter")
	}

	chainConf := l.svcCtx.Config.GetChainConfig(req.ChainId)
	if chainConf == nil {
		return nil, fmt.Errorf("%s", "Miss chain configure")
	}

	sonexConf := l.svcCtx.Config.GetSonexConfig(req.ChainId)
	if sonexConf == nil {
		return nil, fmt.Errorf("%s", "Miss sonex configure")
	}

	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	eventSonexPoolSwapModel := repository.NewEventSonexPoolSwapModel(conn)
	tokenPriceModel := repository.NewSonexTokenPriceRecordModel(conn)

	dataResult := []types.TransactionsResponse_Transaction{}
	var totalResult int64 = -1

	// Filter
	options := repository.QueryMergePoolSwapWithExtraOptions{}
	if req.Wallet != "" && common.IsHexAddress(req.Wallet) {
		walletAddress := common.HexToAddress(req.Wallet).Hex()
		options.Wallet = &walletAddress
	}
	if req.Pool != "" && common.IsHexAddress(req.Pool) {
		poolAddress := common.HexToAddress(req.Pool).Hex()
		options.Pool = &poolAddress
	}
	if req.Token0 != "" && common.IsHexAddress(req.Token0) {
		token0Address := common.HexToAddress(req.Token0)
		var tokenAddress string
		switch token0Address {
		case (common.Address{}):
			tokenAddress = common.HexToAddress(chainConf.WETH).Hex()
			options.IsETH = true
		default:
			tokenAddress = token0Address.Hex()
		}
		if options.Tokens == nil {
			options.Tokens = []string{}
		}
		options.Tokens = append(options.Tokens, tokenAddress)
	}
	if req.Token1 != "" && common.IsHexAddress(req.Token1) {
		token1Address := common.HexToAddress(req.Token1)
		var tokenAddress string
		switch token1Address {
		case (common.Address{}):
			tokenAddress = common.HexToAddress(chainConf.WETH).Hex()
			options.IsETH = true
		default:
			tokenAddress = token1Address.Hex()
		}
		if options.Tokens == nil {
			options.Tokens = []string{}
		}
		options.Tokens = append(options.Tokens, tokenAddress)
	}

	// Page
	if req.OrderBy != "" {
		options.OrderBy = &req.OrderBy
	}
	if req.PageSize == nil || *req.PageSize < 1 {
		var defaultSize int64 = 20
		req.PageSize = &defaultSize
	}
	options.PageSize = req.PageSize
	if req.PageIndex == nil || *req.PageIndex < 1 {
		var defaultIndex int64 = 1
		req.PageIndex = &defaultIndex
	}
	options.PageIndex = req.PageIndex

	fmt.Printf("1 ====> %s\n", time.Since(st))
	st = time.Now()

	// Query
	swapRecords, total, err := eventSonexPoolSwapModel.QueryTransactions(l.ctx, &options)
	if err != nil {
		return nil, err
	}
	totalResult = total

	// Handle data
	transactions := map[string][]repository.EventSonexPoolSwapWithPool{}
	linq.From(swapRecords).GroupByT(
		func(item repository.EventSonexPoolSwapWithPool) string { return item.TxHash },
		func(item repository.EventSonexPoolSwapWithPool) repository.EventSonexPoolSwapWithPool { return item },
	).ToMapByT(
		&transactions,
		func(item linq.Group) string { return item.Key.(string) },
		func(item linq.Group) []repository.EventSonexPoolSwapWithPool {
			var value []repository.EventSonexPoolSwapWithPool
			linq.From(item.Group).SelectT(func(subItem interface{}) repository.EventSonexPoolSwapWithPool {
				return subItem.(repository.EventSonexPoolSwapWithPool)
			}).OrderByT(func(item repository.EventSonexPoolSwapWithPool) int64 { return item.TxIndex }).ToSlice(&value)
			return value
		},
	)

	fmt.Printf("2 ====> %s\n", time.Since(st))
	st = time.Now()

	tokenAddressMap := map[string]struct{}{}
	for _, swaps := range transactions {
		for _, swap := range swaps {
			tokenAddressMap[swap.ExtraToken0] = struct{}{}
			tokenAddressMap[swap.ExtraToken1] = struct{}{}
		}
	}
	tokenAddresses := []string{}
	for key, _ := range tokenAddressMap {
		tokenAddresses = append(tokenAddresses, key)
	}
	prices, err := tokenPriceModel.GetCurrentPrices(l.ctx, tokenAddresses)
	if err != nil {
		return nil, err
	}

	fmt.Printf("3 ====> %s\n", time.Since(st))
	st = time.Now()

	poolMap := map[string]types.Pool{}
	handlePool := func(poolAddress string, token0Address string, token1Address string, fee int64) (*types.Pool, error) {
		pool, ok := poolMap[poolAddress]
		if !ok {
			token0, err := global.GetToken(req.ChainId, token0Address)
			if err != nil {
				return nil, err
			}
			token1, err := global.GetToken(req.ChainId, token1Address)
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
				Address: poolAddress,
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
				Fee: fee,
			}
		}
		return &pool, nil
	}

	var wg sync.WaitGroup
	results := make(chan struct {
		Tx  *types.TransactionsResponse_Transaction
		Err error
	}, len(transactions))

	usdToken, err := global.GetToken(chainConf.ChainId, common.HexToAddress(chainConf.USD).Hex())
	if err != nil {
		return nil, err
	}
	usdAccuracy := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(usdToken.Decimals), nil))
	for txHash, swaps := range transactions {
		wg.Add(1)

		go func(ch chan struct {
			Tx  *types.TransactionsResponse_Transaction
			Err error
		}) {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					ch <- struct {
						Tx  *types.TransactionsResponse_Transaction
						Err error
					}{Err: fmt.Errorf("%s(panic)", err)}
				}
			}()

			startTx := swaps[0]
			endTx := swaps[len(swaps)-1]

			var fromToken *cache.Token
			var fromAmount *big.Int
			var fromPrice *big.Int
			var toToken *cache.Token
			var toAmount *big.Int
			var toPrice *big.Int

			{
				amount0, _ := new(big.Int).SetString(startTx.EventAmount0, 10)
				amount1, _ := new(big.Int).SetString(startTx.EventAmount1, 10)
				if amount0.Cmp(amount1) > 0 {
					fromToken, err = global.GetToken(req.ChainId, startTx.ExtraToken0)
					if err != nil {
						ch <- struct {
							Tx  *types.TransactionsResponse_Transaction
							Err error
						}{Err: err}
						return
					}
					fromAmount = new(big.Int).Abs(amount0)
				} else {
					fromToken, err = global.GetToken(req.ChainId, startTx.ExtraToken1)
					if err != nil {
						ch <- struct {
							Tx  *types.TransactionsResponse_Transaction
							Err error
						}{Err: err}
						return
					}
					fromAmount = new(big.Int).Abs(amount1)
				}
				price, ok := prices[fromToken.Address]
				if ok {
					fromPrice, _ = new(big.Int).SetString(price.Price, 10)
				}
			}
			{
				amount0, _ := new(big.Int).SetString(endTx.EventAmount0, 10)
				amount1, _ := new(big.Int).SetString(endTx.EventAmount1, 10)
				if amount0.Cmp(amount1) < 0 {
					toToken, err = global.GetToken(req.ChainId, endTx.ExtraToken0)
					if err != nil {
						ch <- struct {
							Tx  *types.TransactionsResponse_Transaction
							Err error
						}{Err: err}
						return
					}
					toAmount = new(big.Int).Abs(amount0)
				} else {
					toToken, err = global.GetToken(req.ChainId, endTx.ExtraToken1)
					if err != nil {
						ch <- struct {
							Tx  *types.TransactionsResponse_Transaction
							Err error
						}{Err: err}
						return
					}
					toAmount = new(big.Int).Abs(amount1)
				}
				price, ok := prices[toToken.Address]
				if ok {
					toPrice, _ = new(big.Int).SetString(price.Price, 10)
				}
			}

			pools := []types.Pool{}
			for _, swap := range swaps {
				pool, err := handlePool(swap.ExtraPool, swap.ExtraToken0, swap.ExtraToken1, swap.ExtraFee)
				if err != nil {
					ch <- struct {
						Tx  *types.TransactionsResponse_Transaction
						Err error
					}{Err: err}
					return
				}
				pools = append(pools, *pool)
			}

			_ = toPrice

			fromAccuracy := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(fromToken.Decimals), nil))
			valueUSD, _ := new(big.Float).Quo(new(big.Float).Quo(new(big.Float).SetInt(new(big.Int).Mul(fromAmount, fromPrice)), fromAccuracy), usdAccuracy).Float64()

			ch <- struct {
				Tx  *types.TransactionsResponse_Transaction
				Err error
			}{Tx: &types.TransactionsResponse_Transaction{
				Wallet:     startTx.From,
				Pools:      pools,
				TxHash:     txHash,
				TxAt:       startTx.TxAt.Unix(),
				From:       fromToken.Address,
				To:         toToken.Address,
				FromAmount: fromAmount.String(),
				ToAmount:   toAmount.String(),
				ValueUSD:   fmt.Sprintf("%f", valueUSD),
			}}
		}(results)
	}

	wg.Wait()
	close(results)

	// Handle data
	for result := range results {
		if result.Err != nil {
			return nil, err
		}
		dataResult = append(dataResult, *result.Tx)
	}

	fmt.Printf("4 ====> %s\n", time.Since(st))

	return &response.ApiResponse{
		Data: dataResult,
		Meta: types.PageResponse{
			PageSize:  *req.PageSize,
			PageIndex: *req.PageIndex,
			Total:     totalResult,
		},
	}, nil
}
