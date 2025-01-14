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

	// Query
	swapRecords, total, err := eventSonexPoolSwapModel.QueryMergeWithExtra(l.ctx, &options)
	if err != nil {
		return nil, err
	}
	totalResult = total

	// Handle data
	tokenAddresses := []string{}
	for _, record := range swapRecords {
		tokenAddresses = append(tokenAddresses, record.Tokens()...)
	}
	linq.From(tokenAddresses).Distinct().ToSlice(&tokenAddresses)
	prices, err := tokenPriceModel.GetCurrentPrices(tokenAddresses)
	if err != nil {
		return nil, err
	}

	poolMap := map[string]types.Pool{}
	for _, record := range swapRecords {
		pools := []types.Pool{}
		handlePool := func(poolAddress string, token0Address string, token1Address string) (*types.Pool, error) {
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
					Fee: record.ExtraFee_0,
				}
			}
			return &pool, nil
		}
		pool, err := handlePool(record.ExtraPool_0, record.ExtraToken0_0, record.ExtraToken1_0)
		if err != nil {
			return nil, err
		}
		pools = append(pools, *pool)
		if record.IsCross() {
			pool, err := handlePool(record.ExtraPool_1.String, record.ExtraToken0_1.String, record.ExtraToken1_1.String)
			if err != nil {
				return nil, err
			}
			pools = append(pools, *pool)
		}
		fromToken, err := global.GetToken(req.ChainId, record.FromToken())
		if err != nil {
			return nil, err
		}
		usdToken, err := global.GetToken(chainConf.ChainId, common.HexToAddress(chainConf.USD).Hex())
		if err != nil {
			return nil, err
		}
		usdAccuracy := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(usdToken.Decimals), nil))
		fromAmount := new(big.Float).SetInt(record.FromAmount())
		fromAddress := common.HexToAddress(record.FromToken())
		fromPrice := big.NewFloat(0)
		if fromAddress == (common.Address{}) {
			fromAddress = common.HexToAddress(chainConf.WETH)
		}
		price, ok := prices[fromAddress.Hex()]
		if ok {
			fromPrice, _ = new(big.Float).SetString(price.Price)
		}
		fromAccuracy := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(fromToken.Decimals), nil))
		valueUSD, _ := new(big.Float).Quo(new(big.Float).Quo(new(big.Float).Mul(fromAmount, fromPrice), fromAccuracy), usdAccuracy).Float64()

		dataResult = append(dataResult, types.TransactionsResponse_Transaction{
			Wallet:     record.From,
			Pools:      pools,
			TxHash:     record.TxHash,
			TxAt:       record.TxAt.Unix(),
			From:       record.FromToken(),
			To:         record.ToToken(common.HexToAddress(sonexConf.SwapRouter).Hex()),
			FromAmount: record.FromAmount().String(),
			ToAmount:   record.ToAmount().String(),
			ValueUSD:   fmt.Sprintf("%f", valueUSD),
		})
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
