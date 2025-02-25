package sonex

import (
	"context"
	"sonex_api_backend_go/internal/config"
	contract_common "sonex_api_backend_go/internal/lib/contract/common"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/lib/node_caller"
	"sonex_api_backend_go/internal/repository"

	"github.com/ahmetb/go-linq/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type XTokenHandler struct {
	dataSource string
	chain      config.ChainConfig
	config     config.SonexConfig

	onApproval []func(chainId int64, tx *node_caller.Transaction, event *contract_common.XtokenApproval)
	onTransfer []func(chainId int64, tx *node_caller.Transaction, event *contract_common.XtokenTransfer)
}

func NewXTokenHandler(dataSource string, chain config.ChainConfig, config config.SonexConfig) (*XTokenHandler, error) {
	return &XTokenHandler{
		dataSource: dataSource,
		chain:      chain,
		config:     config,
	}, nil
}

func (handler *XTokenHandler) OnApproval(callback func(chainId int64, tx *node_caller.Transaction, event *contract_common.XtokenApproval)) {
	handler.onApproval = append(handler.onApproval, callback)
}

func (handler *XTokenHandler) OnTransfer(callback func(chainId int64, tx *node_caller.Transaction, event *contract_common.XtokenTransfer)) {
	handler.onTransfer = append(handler.onTransfer, callback)
}

func (handler *XTokenHandler) Execute(client *ethclient.Client, opt *bind.FilterOpts, isReplay bool) (error, bool) {
	nodeCaller, err := node_caller.NewNodeCaller(handler.chain, client)
	if err != nil {
		return err, false
	}
	xtoken, err := contract_common.NewXtoken(common.HexToAddress(handler.config.XToken), client)
	if err != nil {
		return err, false
	}

	numberSet := map[int64]struct{}{}

	// Approval Event
	approvalLogs, err := xtoken.FilterApproval(opt, nil, nil)
	if err != nil {
		return err, false
	}
	approvalEvents := []*contract_common.XtokenApproval{}
	for approvalLogs.Next() {
		if err = approvalLogs.Error(); err != nil {
			return err, false
		}
		numberSet[int64(approvalLogs.Event.Raw.BlockNumber)] = struct{}{}
		approvalEvents = append(approvalEvents, approvalLogs.Event)
	}

	// Transfer Event
	transferLogs, err := xtoken.FilterTransfer(opt, nil, nil)
	if err != nil {
		return err, false
	}
	transferEvents := []*contract_common.XtokenTransfer{}
	for transferLogs.Next() {
		if err = transferLogs.Error(); err != nil {
			return err, false
		}
		numberSet[int64(transferLogs.Event.Raw.BlockNumber)] = struct{}{}
		transferEvents = append(transferEvents, transferLogs.Event)
	}

	var numbers []int64

	linq.From(numberSet).SelectT(func(item linq.KeyValue) int64 { return item.Key.(int64) }).ToSlice(&numbers)
	transactions, err := nodeCaller.BatchGetTranscationsByNumber(opt.Context, numbers)
	if err != nil {
		return err, false
	}

	conn := repository.NewConn(handler.dataSource)
	err = conn.TransactCtx(opt.Context, func(ctx context.Context, s sqlx.Session) error {
		eventXtokenApprovalModel := repository.NewEventXtokenApprovalModel(sqlx.NewSqlConnFromSession(s))
		eventXtokenTransferModel := repository.NewEventXtokenTransferModel(sqlx.NewSqlConnFromSession(s))

		callableApprovalEvents := []*contract_common.XtokenApproval{}
		// Approval Event
		{
			for _, event := range approvalEvents {
				blockNumber := int64(event.Raw.BlockNumber)
				tx := transactions[event.Raw.TxHash.Hex()]
				txIndex := int64(event.Raw.Index)
				txHash := tx.Tx().Hash().Hex()
				txAt := tx.TxAt
				from := tx.From.Hex()
				to := tx.Tx().To().Hex()
				value := tx.Tx().Value().String()
				eventOwner := event.Owner.Hex()
				eventSpender := event.Spender.Hex()
				eventValue := event.Value

				// If handled, ignore
				_, err := eventXtokenApprovalModel.FindOneByChainIdTxHashTxIndex(ctx, handler.chain.ChainId, txHash, txIndex)
				switch err {
				case repository.ErrNotFound:
				case nil:
					continue
				default:
					return err
				}

				id, err := global.NewId()
				if err != nil {
					return err
				}
				err = eventXtokenApprovalModel.New(ctx, id, handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value, eventOwner, eventSpender, eventValue)
				if err != nil {
					return err
				}

				callableApprovalEvents = append(callableApprovalEvents, event)
			}
		}

		callableTransferEvents := []*contract_common.XtokenTransfer{}
		// Transfer Event
		{
			for _, event := range transferEvents {
				blockNumber := int64(event.Raw.BlockNumber)
				tx := transactions[event.Raw.TxHash.Hex()]
				txIndex := int64(event.Raw.Index)
				txHash := tx.Tx().Hash().Hex()
				txAt := tx.TxAt
				from := tx.From.Hex()
				to := tx.Tx().To().Hex()
				value := tx.Tx().Value().String()
				eventFrom := event.From.Hex()
				eventTo := event.To.Hex()
				eventValue := event.Value

				// If handled, ignore
				_, err := eventXtokenTransferModel.FindOneByChainIdTxHashTxIndex(ctx, handler.chain.ChainId, txHash, txIndex)
				switch err {
				case repository.ErrNotFound:
				case nil:
					continue
				default:
					return err
				}

				id, err := global.NewId()
				if err != nil {
					return err
				}
				err = eventXtokenTransferModel.New(ctx, id, handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value, eventFrom, eventTo, eventValue)
				if err != nil {
					return err
				}

				callableTransferEvents = append(callableTransferEvents, event)
			}
		}

		// Callback
		for _, event := range callableApprovalEvents {
			for _, callback := range handler.onApproval {
				tx := transactions[event.Raw.TxHash.Hex()]
				go callback(handler.chain.ChainId, tx, event)
			}
		}

		for _, event := range callableTransferEvents {
			for _, callback := range handler.onTransfer {
				tx := transactions[event.Raw.TxHash.Hex()]
				go callback(handler.chain.ChainId, tx, event)
			}
		}

		return nil
	})
	if err != nil {
		return err, false
	}

	return nil, false
}
