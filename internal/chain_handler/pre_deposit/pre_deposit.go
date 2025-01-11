package pre_deposit

import (
	"context"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/contract/pre_deposit"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/lib/node_caller"
	"sonex_api_backend_go/internal/repository"

	"github.com/ahmetb/go-linq/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type PreDepositHandler struct {
	dataSource string
	chain      *config.ChainConfig
	config     *config.PreDepositConfig

	onStaked    []func(chainId int64, tx *node_caller.Transaction, event *pre_deposit.BridgeStakerStaked)
	onWithdrawn []func(chainId int64, tx *node_caller.Transaction, event *pre_deposit.BridgeStakerWithdrawn)
}

func NewPreDepositHandler(dataSource string, chain config.ChainConfig, config config.PreDepositConfig) (*PreDepositHandler, error) {
	return &PreDepositHandler{
		dataSource:  dataSource,
		chain:       &chain,
		config:      &config,
		onStaked:    []func(chainId int64, tx *node_caller.Transaction, event *pre_deposit.BridgeStakerStaked){},
		onWithdrawn: []func(chainId int64, tx *node_caller.Transaction, event *pre_deposit.BridgeStakerWithdrawn){},
	}, nil
}

func (handler *PreDepositHandler) OnStaked(callback func(chainId int64, tx *node_caller.Transaction, event *pre_deposit.BridgeStakerStaked)) {
	handler.onStaked = append(handler.onStaked, callback)
}

func (handler *PreDepositHandler) OnWithdrawn(callback func(chainId int64, tx *node_caller.Transaction, event *pre_deposit.BridgeStakerWithdrawn)) {
	handler.onWithdrawn = append(handler.onWithdrawn, callback)
}

func (handler *PreDepositHandler) Execute(client *ethclient.Client, opt *bind.FilterOpts, isReplay bool) (error, bool) {
	nodeCaller, err := node_caller.NewNodeCaller(*handler.chain)
	if err != nil {
		return err, false
	}
	bridgeStaker, err := pre_deposit.NewBridgeStaker(common.HexToAddress(handler.config.BridgeStaker), client)
	if err != nil {
		return err, false
	}

	numberSet := map[int64]struct{}{}

	stakedLogs, err := bridgeStaker.FilterStaked(opt, nil, nil)
	if err != nil {
		return err, false
	}
	stakedEvents := []*pre_deposit.BridgeStakerStaked{}
	for stakedLogs.Next() {
		if err = stakedLogs.Error(); err != nil {
			return err, false
		}
		numberSet[int64(stakedLogs.Event.Raw.BlockNumber)] = struct{}{}
		stakedEvents = append(stakedEvents, stakedLogs.Event)
	}

	withdrawnLogs, err := bridgeStaker.FilterWithdrawn(opt, nil)
	if err != nil {
		return err, false
	}
	withdrawnEvents := []*pre_deposit.BridgeStakerWithdrawn{}
	for withdrawnLogs.Next() {
		if err = withdrawnLogs.Error(); err != nil {
			return err, false
		}
		numberSet[int64(withdrawnLogs.Event.Raw.BlockNumber)] = struct{}{}
		withdrawnEvents = append(withdrawnEvents, withdrawnLogs.Event)
	}

	var numbers []int64
	linq.From(numberSet).SelectT(func(item linq.KeyValue) int64 { return item.Key.(int64) }).ToSlice(&numbers)
	transactions, err := nodeCaller.BatchGetTranscationsByNumber(opt.Context, numbers)
	if err != nil {
		return err, false
	}

	conn := repository.NewConn(handler.dataSource)
	err = conn.TransactCtx(opt.Context, func(ctx context.Context, s sqlx.Session) error {
		eventPreDepositStakedModel := repository.NewEventPreDepositStakedModel(sqlx.NewSqlConnFromSession(s))
		eventPreDepositWithdrawnModel := repository.NewEventPreDepositWithdrawnModel(sqlx.NewSqlConnFromSession(s))

		// Staked Event
		{
			for _, event := range stakedEvents {
				blockNumber := int64(event.Raw.BlockNumber)
				tx := transactions[event.Raw.TxHash.Hex()]
				txIndex := int64(event.Raw.Index)
				txHash := tx.Tx().Hash().Hex()
				txAt := tx.TxAt
				from := tx.From.Hex()
				to := tx.Tx().To().Hex()
				value := tx.Tx().Value().String()
				user := event.User.Hex()
				token := event.Token.Hex()
				amount := event.Amount.String()
				timestamp := event.Timestamp.Int64()

				// If handled, ignore
				_, err := eventPreDepositStakedModel.FindOneByChainIdTxHashTxIndex(ctx, handler.chain.ChainId, txHash, txIndex)
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
				err = eventPreDepositStakedModel.New(ctx, id, handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value, user, token, amount, timestamp)
				if err != nil {
					return err
				}
			}
		}

		// Withdrawn Event
		{
			for _, event := range withdrawnEvents {
				blockNumber := int64(event.Raw.BlockNumber)
				tx := transactions[event.Raw.TxHash.Hex()]
				txIndex := int64(event.Raw.Index)
				txHash := tx.Tx().Hash().Hex()
				txAt := tx.TxAt
				from := tx.From.Hex()
				to := tx.Tx().To().Hex()
				value := tx.Tx().Value().String()
				user := event.User.Hex()
				wastrAmount := event.WastrAmount.String()
				vastrAmount := event.VastrAmount.String()
				timestamp := event.Timestamp.Int64()

				// If handled, ignore
				_, err := eventPreDepositWithdrawnModel.FindOneByChainIdTxHashTxIndex(ctx, handler.chain.ChainId, txHash, txIndex)
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
				err = eventPreDepositWithdrawnModel.New(ctx, id, handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value, user, wastrAmount, vastrAmount, timestamp)
				if err != nil {
					return err
				}
			}
		}

		// Callback
		for _, event := range stakedEvents {
			for _, callback := range handler.onStaked {
				tx := transactions[event.Raw.TxHash.Hex()]
				go callback(handler.chain.ChainId, tx, event)
			}
		}
		for _, event := range withdrawnEvents {
			for _, callback := range handler.onWithdrawn {
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
