package sonex

import (
	"context"
	"fmt"
	"math/big"

	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

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

	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	eventSonexPoolSwapModel := repository.NewEventSonexPoolSwapModel(conn)

	var swapRecords []repository.EventSonexPoolSwapWithExtra
	dataResult := []types.TransactionsResponse_Transaction{}
	// var totalResult int64
	if req.Wallet != "" && common.IsHexAddress(req.Wallet) {
		walletAddress := common.HexToAddress(req.Wallet)
		swapRecords, _, err = eventSonexPoolSwapModel.PageQueryByWalletWithExtra(l.ctx, walletAddress.Hex(), req.PageSize, req.PageIndex)
		if err != nil {
			return nil, err
		}
	}

	pools := map[string]types.Pool{}
	for _, record := range swapRecords {
		poolAddress := common.HexToAddress(record.ExtraPool)
		pool, ok := pools[poolAddress.Hex()]
		if !ok {

		}

		var from string
		var to string
		amount0, _ := new(big.Int).SetString(record.EventAmount0, 10)
		amount1, _ := new(big.Int).SetString(record.EventAmount1, 10)
		if amount0.Cmp(amount1) < 0 {
			from = record.ExtraToken0
			to = record.ExtraToken1
		} else {
			from = record.ExtraToken1
			to = record.ExtraToken0
		}

		dataResult = append(dataResult, types.TransactionsResponse_Transaction{
			Wallet: record.From,
			Pool:   pool,
			TxHash: record.TxHash,
			TxAt:   record.TxAt.Unix(),
			From:   from,
			To:     to,
		})
	}

	return &types.TransactionsResponse{
		Data: dataResult,
	}, nil
}
